package cheque

import (
	"fmt"
	"github.com/bittorrent/go-btfs/utils"
	"io"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
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
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		walletImportPrvKey, err := chain.GetWalletImportPrvKey(env)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &ChainInfoRet{
			ChainId:            chain.ChainObject.ChainID,
			NodeAddr:           chain.ChainObject.OverlayAddress.String(),
			VaultAddr:          chain.SettleObject.VaultService.Address().String(),
			WalletImportPrvKey: walletImportPrvKey,
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
