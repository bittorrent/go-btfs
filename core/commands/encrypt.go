package commands

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"os"

	shell "github.com/bittorrent/go-btfs-api"
	cmds "github.com/bittorrent/go-btfs-cmds"
	cp "github.com/bittorrent/go-btfs-common/crypto"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/corehttp/remote"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	peer "github.com/libp2p/go-libp2p/core/peer"
)

const toOption = "to"
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
	},
	Run: func(r *cmds.Request, re cmds.ResponseEmitter, e cmds.Environment) error {
		n, err := cmdenv.GetNode(e)
		if err != nil {
			return err
		}
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
		it := r.Files.Entries()
		file, err := cmdenv.GetFileArg(it)
		if err != nil {
			return err
		}
		originalBytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		encryptedBytes, err := ECCEncrypt(originalBytes, *eciesPk)
		if err != nil {
			return err
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
	Arguments: []cmds.Argument{
		cmds.StringArg("cid", true, false, "the CID of the encrypted file"),
	},
	Options: []cmds.Option{
		cmds.StringOption(fromOption, "specify the source peerID of CID"),
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
		if ok {
			peerID, err := peer.Decode(from)
			if err != nil {
				return err
			}
			b, err := remote.P2PCallStrings(r.Context, n, api, peerID, "/decryption", cid)
			if err != nil {
				return err
			}
			readClose = io.NopCloser(bytes.NewReader(b))
		} else {
			readClose, err = shell.NewLocalShell().Cat(cid)
			if err != nil {
				return err
			}
		}
		ecdsaPrivateKey, err := ethCrypto.HexToECDSA(conf.Identity.HexPrivKey)
		if err != nil {
			return err
		}
		eciesPrivateKey := ecies.ImportECDSA(ecdsaPrivateKey)
		endata, err := io.ReadAll(readClose)
		if err != nil {
			return err
		}
		defer readClose.Close()
		dedata, err := ECCDecrypt(endata, *eciesPrivateKey)
		if err != nil {
			panic(err)
		}
		fileName := "./decrypt-file-of-" + cid
		f, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(dedata)
		if err != nil {
			return err
		}
		return re.Emit("decrypted file name is: " + fileName)
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
