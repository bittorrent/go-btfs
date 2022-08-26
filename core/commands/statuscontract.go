package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"time"

	cmdss "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	cmdenvv "github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/reportstatus"
	"github.com/bittorrent/go-btfs/spin"
)

var StatusContractCmd = &cmdss.Command{
	Helptext: cmdss.HelpText{
		Tagline: "report status-contract cmd.",
		ShortDescription: `
report status-contract cmd, total cmd and list cmd.`,
	},
	Subcommands: map[string]*cmdss.Command{
		"total":                  TotalCmd,
		"reportlist":             ReportListCmd,
		"lastinfo":               LastInfoCmd,
		"config":                 StatusConfigCmd,
		"report_online_server":   ReportOnlineServerCmd,
		"report_status_contract": ReportStatusContractCmd,
	},
}

type TotalCmdRet struct {
	PeerId         string `json:"peer_id"`
	StatusContract string `json:"status_contract"`
	TotalCount     int    `json:"total_count"`
	TotalGasSpend  string `json:"total_gas_spend"`
}

var TotalCmd = &cmdss.Command{
	Helptext: cmdss.HelpText{
		Tagline: "report status-contract total info, (total count, total gas spend, and contract address)",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmdss.Request, res cmdss.ResponseEmitter, env cmdss.Environment) error {
		n, err := cmdenvv.GetNode(env)
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
				//fmt.Println("r.GasSpend is zero. ")
				continue
			}
			//fmt.Println("r.GasSpend = ", r.GasSpend)

			n, ok := n.SetString(r.GasSpend, 10)
			if !ok {
				return errors.New("parse gas_spend is error. ")
			}
			totalGasSpend = totalGasSpend.Add(totalGasSpend, n)
		}

		return cmdss.EmitOnce(res, &TotalCmdRet{
			PeerId:         peerId,
			StatusContract: list[len(list)-1].StatusContract,
			TotalCount:     len(list),
			TotalGasSpend:  totalGasSpend.String(),
		})
	},
	Type: TotalCmdRet{},
	Encoders: cmdss.EncoderMap{
		cmdss.Text: cmdss.MakeTypedEncoder(func(req *cmdss.Request, w io.Writer, out *TotalCmdRet) error {
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

var ReportListCmd = &cmdss.Command{
	Helptext: cmdss.HelpText{
		Tagline: "report status-contract list, and input from and limit to get its.",
	},
	RunTimeout: 5 * time.Minute,
	Arguments: []cmdss.Argument{
		cmdss.StringArg("from", true, false, "page offset"),
		cmdss.StringArg("limit", true, false, "page limit."),
	},
	Run: func(req *cmdss.Request, res cmdss.ResponseEmitter, env cmdss.Environment) error {
		n, err := cmdenvv.GetNode(env)
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

		s := list[From:To]
		l := len(s)
		for i, j := 0, l-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}

		return cmdss.EmitOnce(res, &ReportListCmdRet{
			Records: s,
			Total:   len(list),
			PeerId:  peerId,
		})
	},
	Type: ReportListCmdRet{},
	Encoders: cmdss.EncoderMap{
		cmdss.Text: cmdss.MakeTypedEncoder(func(req *cmdss.Request, w io.Writer, out *ReportListCmdRet) error {
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

var LastInfoCmd = &cmdss.Command{
	Helptext: cmdss.HelpText{
		Tagline: "get reporting status-contract last info",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmdss.Request, res cmdss.ResponseEmitter, env cmdss.Environment) error {
		last, err := chain.GetLastOnline()
		if err != nil {
			return err
		}
		if last == nil {
			return errors.New("not found. ")
		}

		return cmdss.EmitOnce(res, last)
	},
	Type: chain.LastOnlineInfo{},
	Encoders: cmdss.EncoderMap{
		cmdss.Text: cmdss.MakeTypedEncoder(func(req *cmdss.Request, w io.Writer, out *chain.LastOnlineInfo) error {
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

var StatusConfigCmd = &cmdss.Command{
	Helptext: cmdss.HelpText{
		Tagline: "get reporting status-contract config. ",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmdss.Request, res cmdss.ResponseEmitter, env cmdss.Environment) error {
		rs, err := chain.GetReportStatus()
		if err != nil {
			return err
		}

		return cmdss.EmitOnce(res, rs)
	},
	Type: chain.ReportStatusInfo{},
	Encoders: cmdss.EncoderMap{
		cmdss.Text: cmdss.MakeTypedEncoder(func(req *cmdss.Request, w io.Writer, out *chain.ReportStatusInfo) error {
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

var ReportOnlineServerCmd = &cmdss.Command{
	Helptext: cmdss.HelpText{
		Tagline: "report online server. ",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmdss.Request, res cmdss.ResponseEmitter, env cmdss.Environment) error {
		node, err := cmdenvv.GetNode(env)
		if err != nil {
			return err
		}

		cfg, err := cmdenvv.GetConfig(env)
		if err != nil {
			return err
		}

		spin.DC.SendDataOnline(node, cfg)

		return cmdss.EmitOnce(res, "report online server ok!")
	},
}

var ReportStatusContractCmd = &cmdss.Command{
	Helptext: cmdss.HelpText{
		Tagline: "report status-contract. ",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmdss.Request, res cmdss.ResponseEmitter, env cmdss.Environment) error {
		err := reportstatus.CmdReportStatus()
		if err != nil {
			return err
		}

		return cmdss.EmitOnce(res, "report status contract ok!")
	},
}
