package upload

import (
	"context"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/corehttp/remote"
	"github.com/bittorrent/go-btfs/settlement/swap/swapprotocol"
)

var StorageUploadChequeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "receive upload cheque, do with cheque, and return it.",
		ShortDescription: `receive upload cheque, deal it and return it.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("encoded-cheque", true, false, "encoded-cheque from peer-id."),
		cmds.StringArg("amount", true, false, "amount"),
		cmds.StringArg("contract-id", false, false, "contract-id."),
		cmds.StringArg("token", true, false, "token"),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		fmt.Printf("receive cheque ...\n")

		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}
		if !ctxParams.Cfg.Experimental.StorageHostEnabled {
			return fmt.Errorf("storage host api not enabled")
		}

		requestPid, ok := remote.GetStreamRequestRemotePeerID(req, ctxParams.N)
		if !ok {
			return fmt.Errorf("fail to get peer ID from request")
		}

		price, ok := new(big.Int).SetString(req.Arguments[1], 10)
		if !ok {
			return fmt.Errorf("exchangeRate:%s cannot be parsed, err:%s", req.Arguments[2], err)
		}

		encodedCheque := req.Arguments[0]
		contractId := req.Arguments[2]
		tokenStr := req.Arguments[3]
		fmt.Printf("receive cheque, requestPid:%s contractId:%+v,encodedCheque:%+v token:%+v\n",
			requestPid.String(), contractId, encodedCheque, tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		// decode and deal the cheque
		err = swapprotocol.SwapProtocol.Handler(context.Background(), requestPid.String(), encodedCheque, price, token)
		if err != nil {
			fmt.Println("receive cheque, swapprotocol.SwapProtocol.Handler, error:", err)
			return err
		}

		// if receive cheque of contractId, set shard paid status.
		if len(contractId) > 0 {
			err := setPaidStatus(ctxParams, contractId)
			if err != nil {
				fmt.Println("receive cheque, setPaidStatus: contractId error:", contractId, err)
				return err
			}
		}

		return nil
	},
}
