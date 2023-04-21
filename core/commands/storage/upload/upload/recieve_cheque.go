package upload

import (
	"context"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"github.com/ethereum/go-ethereum/common"
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
		cmds.StringArg("token", false, false, "token."),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

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
		tokenHex := req.Arguments[3]
		fmt.Printf("receive cheque, requestPid:%s contractId:%+v,encodedCheque:%+v price:%v token:%+v\n",
			requestPid.String(), contractId, encodedCheque, price, tokenHex)

		token := common.HexToAddress(tokenHex)
		_, bl := tokencfg.MpTokenStr[token]
		if !bl {
			return errors.New("your input token is none. ")
		}

		// check price
		priceStore, amountStore, rateStore, err := getInputPriceAmountRate(ctxParams, contractId)
		if err != nil {
			return err
		}
		if price.Int64() < priceStore {
			return errors.New(
				fmt.Sprintf("receive cheque, your input-price[%v] is less than store-price[%v]. ",
					price, priceStore),
			)
		}
		realAmount := new(big.Int).Mul(big.NewInt(amountStore), rateStore)
		fmt.Printf("receive cheque, price:%v amountStore:%v rateStore:%+v,realAmount:%+v \n",
			priceStore, amountStore, rateStore.String(), realAmount.String())

		// decode and deal the cheque
		err = swapprotocol.SwapProtocol.Handler(context.Background(), requestPid.String(), encodedCheque, realAmount, token)
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
