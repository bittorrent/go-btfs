package cheque

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"io"
	"math/big"
	"sort"
	"strconv"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

type ChequeCashListRet struct {
	TxHash   string   `json:"tx_hash"`
	PeerID   string   `json:"peer_id"`
	Vault    string   `json:"vault"`
	Amount   *big.Int `json:"amount"`
	CashTime int64    `json:"cash_time"`
	Status   string   `json:"status"`
}
type ChequeCashListRsp struct {
	Records []ChequeCashListRet `json:"records"`
	Total   int                 `json:"total"`
}

var ChequeCashListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get cash status by peerID.",
	},
	RunTimeout: 5 * time.Minute,
	Arguments: []cmds.Argument{
		cmds.StringArg("from", true, false, "page offset"),
		cmds.StringArg("limit", true, false, "page limit."),
	},
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		from, err := strconv.Atoi(req.Arguments[0])
		if err != nil {
			return fmt.Errorf("parse from:%v failed", req.Arguments[0])
		}
		limit, err := strconv.Atoi(req.Arguments[1])
		if err != nil {
			return fmt.Errorf("parse limit:%v failed", req.Arguments[1])
		}
		if from < 0 {
			return fmt.Errorf("invalid from: %d", from)
		}
		if limit < 0 {
			return fmt.Errorf("invalid limit: %d", limit)
		}

		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		fmt.Printf("... token:%+v\n", tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		results, err := chain.SettleObject.CashoutService.CashoutResults(token)
		if err != nil {
			return err
		}
		sort.Slice(results, func(i, j int) bool {
			return results[i].CashTime > results[j].CashTime
		})

		total := len(results)
		ret := make([]ChequeCashListRet, 0, limit)
		if from < len(results) {
			if (from + limit) <= len(results) {
				results = results[from : from+limit]
			} else {
				results = results[from:]
			}
			for _, result := range results {
				peer, known, err := chain.SettleObject.SwapService.VaultPeer(result.Vault)
				if err == nil {
					if !known {
						peer = "unkonwn"
					}
					r := ChequeCashListRet{
						TxHash:   result.TxHash.String(),
						PeerID:   peer,
						Vault:    result.Vault.String(),
						Amount:   result.Amount,
						CashTime: result.CashTime,
						Status:   result.Status,
					}
					ret = append(ret, r)
				}
			}
		}

		return cmds.EmitOnce(res, &ChequeCashListRsp{
			Records: ret,
			Total:   total,
		})
	},
	Type: &ChequeCashListRsp{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *ChequeCashListRsp) error {
			marshaled, err := json.MarshalIndent(out, "", "\t")
			if err != nil {
				return err
			}
			marshaled = append(marshaled, byte('\n'))
			fmt.Fprintln(w, string(marshaled))
			return nil
		}),
	},
}
