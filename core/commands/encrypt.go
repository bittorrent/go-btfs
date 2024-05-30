package commands

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	shell "github.com/bittorrent/go-btfs-api"
	cmds "github.com/bittorrent/go-btfs-cmds"
	cp "github.com/bittorrent/go-btfs-common/crypto"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/corehttp/remote"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/mitchellh/go-homedir"
)

const toOption = "to"
const passOption = "pass"
const fromOption = "from"

var encryptCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "encrypt file with the public key of the peer",
	},
	Arguments: []cmds.Argument{
		cmds.FileArg("path", true, true, "The path to a file to be added to btfs.").EnableRecursive().EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.StringOption(toOption, "the peerID of the node which you want to share with"),
		cmds.StringOption(passOption, "p", "the password that you want to encrypt the file by AES"),
	},
	Run: func(r *cmds.Request, re cmds.ResponseEmitter, e cmds.Environment) error {
		n, err := cmdenv.GetNode(e)
		if err != nil {
			return err
		}
		it := r.Files.Entries()
		file, err := cmdenv.GetFileArg(it)
		if err != nil {
			return err
		}
		originalBytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		var encryptedBytes []byte
		pass, ok := r.Options[passOption].(string)
		if ok {
			// That means it's symmetrical encryption
			hasher := md5.New()
			hasher.Write([]byte(pass))
			secretBytes := hasher.Sum(nil)
			encryptedBytes, err = AesEncrypt(originalBytes, secretBytes)
			if err != nil {
				return err
			}
		} else {
			// That means it's asymmetric encryption
			to, ok := r.Options[toOption].(string)
			if !ok {
				to = n.Identity.String()
			}
			id, err := peer.Decode(to)
			if err != nil {
				return errors.New("the to option must be a valid peerID")
			}
			p2pPk, err := id.ExtractPublicKey()
			if err != nil {
				return errors.New("can't extract public key from peerID")
			}
			pkBytes, err := cp.Secp256k1PublicKeyRaw(p2pPk)
			if err != nil {
				return errors.New("can't change from p2p public key to secp256k1 public key from peerID")
			}

			ethPk, err := ethCrypto.UnmarshalPubkey(pkBytes)
			if err != nil {
				return errors.New("can't unmarshall public key from peerID")
			}

			eciesPk := ecies.ImportECDSAPublic(ethPk)
			encryptedBytes, err = ECCEncrypt(originalBytes, *eciesPk)
			if err != nil {
				return err
			}
		}

		btfsClient := shell.NewLocalShell()
		cid, err := btfsClient.Add(bytes.NewReader(encryptedBytes), shell.Pin(true))
		if err != nil {
			return err
		}
		return re.Emit(cid)
	},
}

var decryptCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "decrypt the content of a CID with the private key of this peer",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("cid", true, false, "the CID of the encrypted file"),
	},
	Options: []cmds.Option{
		cmds.StringOption(fromOption, "specify the source peerID of CID"),
		cmds.StringOption(passOption, "p", "the password that you want to decrypt the file by AES"),
	},
	Run: func(r *cmds.Request, re cmds.ResponseEmitter, e cmds.Environment) error {
		conf, err := cmdenv.GetConfig(e)
		if err != nil {
			return err
		}
		n, err := cmdenv.GetNode(e)
		if err != nil {
			return err
		}
		api, err := cmdenv.GetApi(e, r)
		if err != nil {
			return err
		}

		var readClose io.ReadCloser
		cid := r.Arguments[0]
		from, ok := r.Options[fromOption].(string)
		if ok && strings.TrimSpace(from) != "" && strings.TrimSpace(from) != n.Identity.String() {
			peerID, err := peer.Decode(from)
			if err != nil {
				return err
			}
			b, err := remote.P2PCallStrings(r.Context, n, api, peerID, "/decryption", cid)
			if err != nil && strings.Contains(err.Error(), "unsupported path namespace") {
				return errors.New("cid not found")
			}
			if err != nil {
				return err
			}
			readClose = io.NopCloser(bytes.NewReader(b))
		} else {
			c := &http.Client{
				Transport: &http.Transport{
					Proxy:             http.ProxyFromEnvironment,
					DisableKeepAlives: true,
				},
				Timeout: 1 * time.Minute,
			}
			baseDir := os.Getenv("BTFS_PATH")
			if baseDir == "" {
				baseDir = "~/.btfs"
			}

			baseDir, err := homedir.Expand(baseDir)
			if err != nil {
				return nil
			}

			apiFile := path.Join(baseDir, "api")
			if _, err := os.Stat(apiFile); err != nil {
				return nil
			}
			api, err := os.ReadFile(apiFile)
			if err != nil {
				return nil
			}
			readClose, err = shell.NewShellWithClient(strings.TrimSpace(string(api)), c).Cat(cid)
			if err != nil && strings.Contains(err.Error(), "unsupported path namespace") {
				return errors.New("cid not found")
			}
			if err != nil && strings.Contains(err.Error(), "Timeout") {
				return fmt.Errorf("timeout when try to get cid: %s", cid)
			}
		}
		encryptedData, err := io.ReadAll(readClose)
		if err != nil {
			return err
		}
		defer readClose.Close()

		var decryptedData []byte
		pass, ok := r.Options[passOption].(string)
		if ok && strings.TrimSpace(pass) != "" {
			// That means it's symmetrical encryption
			hasher := md5.New()
			hasher.Write([]byte(pass))
			secretBytes := hasher.Sum(nil)
			decryptedData, err = AesDecrypt(encryptedData, secretBytes)
			if err != nil {
				log.Error(err)
				return errors.New("error happens during decryption by AES, may be you have a wrong password or it is not encrypted by AES")
			}
		} else {
			// That means it's asymmetric encryption
			pkbytesOri, err := base64.StdEncoding.DecodeString(conf.Identity.PrivKey)
			if err != nil {
				return err
			}
			ecdsaPrivateKey, err := ethCrypto.ToECDSA(pkbytesOri[4:])
			if err != nil {
				return err
			}
			eciesPrivateKey := ecies.ImportECDSA(ecdsaPrivateKey)

			decryptedData, err = ECCDecrypt(encryptedData, *eciesPrivateKey)
			if err != nil {
				log.Error(err)
				return errors.New("decryption is failed, maybe the content of encryption is not encrypted by your public key")
			}
		}
		return re.Emit(bytes.NewReader(decryptedData))
	},
}

func ECCEncrypt(pt []byte, puk ecies.PublicKey) ([]byte, error) {
	ct, err := ecies.Encrypt(rand.Reader, &puk, pt, nil, nil)
	return ct, err
}

func ECCDecrypt(ct []byte, prk ecies.PrivateKey) ([]byte, error) {
	pt, err := prk.Decrypt(ct, nil, nil)
	return pt, err
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AES decryption,CBC
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AES decryption
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
