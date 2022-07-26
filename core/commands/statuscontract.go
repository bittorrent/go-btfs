package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"io"
	"math/big"
	"strconv"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

var StatusContractCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Interact with vault services on BTFS.",
		ShortDescription: `
Vault services include issue cheque to peer, receive cheque and store operations.`,
	},
	Subcommands: map[string]*cmds.Command{
		"total":      TotalCmd,
		"reportlist": ReportListCmd,
	},
}

type TotalCmdRet struct {
	PeerId         string `json:"peer_id"`
	StatusContract string `json:"status_contract"`
	TotalCount     int    `json:"total_count"`
	TotalGasSpend  string `json:"total_gas_spend"`
}

var TotalCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "report status contract total info.",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		peerId := n.Identity.Pretty()

		list, err := chain.GetReportStatusListOK()
		if err != nil {
			return err
		}
		if list == nil {
			return nil
		}

		// get list to cal spend, from string to big.int...
		totalGasSpend := new(big.Int)
		for _, r := range list {
			n := new(big.Int)
			if len(r.GasSpend) <= 0 {
				fmt.Println("r.GasSpend is zero. ")
				continue
			}
			fmt.Println("r.GasSpend = ", r.GasSpend)

			n, ok := n.SetString(r.GasSpend, 10)
			if !ok {
				return errors.New("parse gas_spend is error. ")
			}
			totalGasSpend = totalGasSpend.Add(totalGasSpend, n)
		}

		return cmds.EmitOnce(res, &TotalCmdRet{
			PeerId:         peerId,
			StatusContract: list[len(list)-1].StatusContract,
			TotalCount:     len(list),
			TotalGasSpend:  totalGasSpend.String(),
		})
	},
	Type: TotalCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *TotalCmdRet) error {
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

type ReportListCmdRet struct {
	Records []*chain.LevelDbReportStatusInfo `json:"records"`
	Total   int                              `json:"total"`
	PeerId  string                           `json:"peer_id"`
}

var ReportListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "report status contract list.",
	},
	RunTimeout: 5 * time.Minute,
	Arguments: []cmds.Argument{
		cmds.StringArg("from", true, false, "page offset"),
		cmds.StringArg("limit", true, false, "page limit."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		peerId := n.Identity.Pretty()

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

		list, err := chain.GetReportStatusListOK()
		if err != nil {
			return err
		}
		if list == nil {
			return nil
		}
		//
		//from := 0
		//limit := 10

		From := len(list) - 1 - from - limit
		if From <= 0 {
			From = 0
		}
		To := len(list) - 1 - from
		if To > len(list)-1 {
			To = len(list) - 1
		}
		fmt.Println("From, To = ", From, To)

		return cmds.EmitOnce(res, &ReportListCmdRet{
			Records: list[From:To], //这里可能对。
			Total:   len(list),
			PeerId:  peerId,
		})
	},
	Type: ReportListCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *ReportListCmdRet) error {
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
