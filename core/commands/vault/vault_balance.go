package vault

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"io"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"golang.org/x/net/context"
)

type VaultBalanceCmdRet struct {
	Balance *big.Int `json:"balance"`
}

var VaultBalanceCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get vault balance.",
	},
	RunTimeout: 5 * time.Minute,
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		//fmt.Printf("... token:%+v\n", tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		balance, err := chain.SettleObject.VaultService.AvailableBalance(context.Background(), token)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &VaultBalanceCmdRet{
			Balance: balance,
		})
	},
	Type: &VaultBalanceCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *VaultBalanceCmdRet) error {
			_, err := fmt.Fprintf(w, "the vault available balance: %v\n", out.Balance)
			return err
		}),
	},
}
