package cheque

import (
	"fmt"
	"io"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	oldcmds "github.com/bittorrent/go-btfs/commands"
	"github.com/tron-us/go-btfs-common/crypto"
)

type ChainInfoRet struct {
	ChainId            int64  `json:"chain_id"`
	NodeAddr           string `json:"node_addr"`
	VaultAddr          string `json:"vault_addr"`
	WalletImportPrvKey string `json:"wallet_import_prv_key"`
}

var ChequeChainInfoCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Show current chain info.",
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		cctx := env.(*oldcmds.Context)
		cfg, err := cctx.GetConfig()
		if err != nil {
			return err
		}
		privKey, err := crypto.ToPrivKey(cfg.Identity.PrivKey)
		if err != nil {
			return err
		}
		keys, err := crypto.FromIcPrivateKey(privKey)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &ChainInfoRet{
			ChainId:            chain.ChainObject.ChainID,
			NodeAddr:           chain.ChainObject.OverlayAddress.String(),
			VaultAddr:          chain.SettleObject.VaultService.Address().String(),
			WalletImportPrvKey: keys.HexPrivateKey,
		})
	},
	Type: &ChainInfoRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *ChainInfoRet) error {
			_, err := fmt.Fprintf(w, "chain id:\t%d\nnode addr:\t%s\nvault addr:\t%s\nwallet import private key:\t%s\n", out.ChainId,
				out.NodeAddr, out.VaultAddr, out.WalletImportPrvKey)
			return err
		}),
	},
}
