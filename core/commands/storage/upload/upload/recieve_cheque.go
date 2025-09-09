package upload

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"github.com/ethereum/go-ethereum/common"

	nodepb "github.com/bittorrent/go-btfs-common/protos/node"

	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"

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

		// it means the cheque is used to renew
		if strings.HasPrefix(contractId, "renewal_") {
			contractId, duration, err := parseContractIdAndRenewalDuration(contractId)
			if err != nil {
				fmt.Println("receive cheque, parseContractIdAndRenewalDuration, error:", err)
				return err
			}
			err = extendShardEndTime(ctxParams, contractId, duration)
			if err != nil {
				fmt.Println("receive renewal cheque, updateShardEndTime error:", err)
				return err
			}
			return nil
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

func parseContractIdAndRenewalDuration(renewContractId string) (contractId string, duration int, err error) {
	// renewal_%d_%s", duration, contractId
	strings := strings.Split(renewContractId, "_")
	if len(strings) != 3 {
		return "", 0, fmt.Errorf("bad renew contract id: fewer than 3 segments")
	}
	duration, err = strconv.Atoi(strings[1])
	if err != nil {
		return "", 0, fmt.Errorf("bad renew contract id: duration is not a number")
	}
	contractId = strings[2]
	return
}

func extendShardEndTime(ctxParams *uh.ContextParams, contractId string, duration int) error {
	key, cs, err := sessions.GetUserShardContract(ctxParams.N.Repo.Datastore(), ctxParams.N.Identity.String(), nodepb.ContractStat_HOST.String(), contractId)
	if err != nil {
		return err
	}

	if cs.Meta.ContractId == contractId {
		cs.Meta.StorageEnd = uint64(time.Unix(int64(cs.Meta.StorageEnd), 0).Add(time.Duration(duration) * time.Hour * 24).Unix())
		err := sessions.UpdateShardContract(ctxParams.N.Repo.Datastore(), cs, key)
		if err != nil {
			return err
		}
	}
	return nil
}
