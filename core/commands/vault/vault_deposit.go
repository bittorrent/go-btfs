package vault

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"io"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/utils"
	"golang.org/x/net/context"
)

type VaultDepositCmdRet struct {
	Hash string `json:"hash"`
}

var VaultDepositCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Deposit from beneficiary to vault contract account.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("amount", true, false, "deposit amount."),
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

		argAmount := utils.RemoveSpaceAndComma(req.Arguments[0])
		amount, ok := new(big.Int).SetString(argAmount, 10)
		if !ok {
			return fmt.Errorf("amount:%s cannot be parsed", req.Arguments[0])
		}

		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		//fmt.Printf("... token:%+v\n", tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		hash, err := chain.SettleObject.VaultService.Deposit(context.Background(), amount, token)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &VaultDepositCmdRet{
			Hash: hash.String(),
		})
	},
	Type: &VaultDepositCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *VaultDepositCmdRet) error {
			_, err := fmt.Fprintf(w, "the hash of transaction: %s\n", out.Hash)
			return err
		}),
	},
}
