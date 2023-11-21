package commands

import (
	"crypto/rand"
	"errors"
	"io"

	cmds "github.com/bittorrent/go-btfs-cmds"
	cp "github.com/bittorrent/go-btfs-common/crypto"
	files "github.com/bittorrent/go-btfs-files"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	iface "github.com/bittorrent/interface-go-btfs-core"
	path "github.com/bittorrent/interface-go-btfs-core/path"
	"github.com/ethereum/go-ethereum/crypto"
	eth "github.com/ethereum/go-ethereum/crypto"
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
		api, err := cmdenv.GetApi(e, r)
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

		ethPk, err := eth.UnmarshalPubkey(pkBytes)
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
		result, err := api.Unixfs().Add(n.Context(), files.NewBytesFile(encryptedBytes))
		if err != nil {
			return err
		}
		re.Emit(result.Cid().String())
		return nil
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
		api, err := cmdenv.GetApi(e, r)
		if err != nil {
			return err
		}

		conf, err := cmdenv.GetConfig(e)
		if err != nil {
			return err
		}
		_, ok := r.Options[fromOption].(string)
		if ok {
			// TODO: get cid from fromOption(remoteCall)
		}
		cid := r.Arguments[0]
		p := path.New(cid)
		f, err := api.Unixfs().Get(r.Context, p)
		if err != nil {
			return err
		}
		var file files.File
		switch f := f.(type) {
		case files.File:
			file = f
		case files.Directory:
			return iface.ErrIsDir
		default:
			return iface.ErrNotSupported
		}
		ecdsaPrivateKey, err := crypto.HexToECDSA(conf.Identity.HexPrivKey)
		if err != nil {
			return err
		}
		eciesPrivateKey := ecies.ImportECDSA(ecdsaPrivateKey)
		endata, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		dedata, err := ECCDecrypt(endata, *eciesPrivateKey)
		if err != nil {
			panic(err)
		}
		re.Emit(string(dedata))
		return nil
	},
}

var getEncryptedCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "get encrypted file by cid",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("cid", true, false, "cid of encrypted file"),
	},

	Run: func(r *cmds.Request, re cmds.ResponseEmitter, e cmds.Environment) error {
		re.Emit("hello world!")
		return nil
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
