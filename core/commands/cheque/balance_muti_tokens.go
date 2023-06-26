package cheque

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
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/net/context"
)

var ChequeAllTokenBalanceCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get all token balance by addr.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("addr", true, false, "bttc account address"),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		addr := req.Arguments[0]

		mp := make(map[string]*big.Int, 0)
		for k := range tokencfg.MpTokenAddr {
			balance, err := chain.SettleObject.VaultService.TokenBalanceOf(context.Background(), common.HexToAddress(addr), k)
			if err != nil {
				return err
			}

			mp[k] = balance
		}

		return cmds.EmitOnce(res, &mp)
	},
	Type: &ChequeBttBalanceCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *ChequeBttBalanceCmdRet) error {
			_, err := fmt.Fprintf(w, "the balance is: %v\n", out.Balance)
			return err
		}),
	},
}

var ChequeTokenBalanceCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get one token balance by addr.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("addr", true, false, "bttc account address"),
	},
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		addr := req.Arguments[0]

		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		//fmt.Printf("... token:%+v\n", tokenStr)
		_, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		balance, err := chain.SettleObject.VaultService.TokenBalanceOf(context.Background(), common.HexToAddress(addr), tokenStr)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &ChequeBttBalanceCmdRet{
			Balance: balance,
		})
	},
	Type: &ChequeBttBalanceCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *ChequeBttBalanceCmdRet) error {
			_, err := fmt.Fprintf(w, "the balance is: %v\n", out.Balance)
			return err
		}),
	},
}
