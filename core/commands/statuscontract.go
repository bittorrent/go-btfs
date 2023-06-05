package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	onlinePb "github.com/bittorrent/go-btfs-common/protos/online"
	"github.com/bittorrent/go-btfs/utils"
	"io"
	"math/big"
	"strconv"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/reportstatus"
	"github.com/bittorrent/go-btfs/spin"
)

var StatusContractCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "report status-contract cmd.",
		ShortDescription: `
report status-contract cmd, total cmd and list cmd.`,
	},
	Subcommands: map[string]*cmds.Command{
		"total":      TotalCmd,
		"reportlist": ReportListCmd,
		"config":     StatusConfigCmd,
		//"report_status_contract": ReportStatusContractCmd,

		"lastinfo":             LastInfoCmd,
		"report_online_server": ReportOnlineServerCmd,

		"daily_report_online_server": ReportOnlineDailyCmd,
		"daily_report_list":          ReportListDailyCmd,
		"daily_total":                TotalDailyCmd,
		"daily_last_report_time":     ReportLastTimeDailyCmd,
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
		Tagline: "(old)report status-contract total info, (total count, total gas spend, and contract address)",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

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
		Tagline: "(old)report status-contract list, and input from and limit to get its.",
	},
	RunTimeout: 5 * time.Minute,
	Arguments: []cmds.Argument{
		cmds.StringArg("from", true, false, "page offset"),
		cmds.StringArg("limit", true, false, "page limit."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

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
		total := len(list)
		// order by time desc
		for i, j := 0, total-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
		if from < total {
			if (from + limit) <= len(list) {
				list = list[from : from+limit]
			} else {
				list = list[from:]
			}
		}

		return cmds.EmitOnce(res, &ReportListCmdRet{
			Records: list,
			Total:   total,
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

var LastInfoCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "get reporting status-contract last info",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		last, err := chain.GetLastOnline()
		if err != nil {
			return err
		}
		if last == nil {
			return errors.New("not found. ")
		}

		r := chain.LastOnlineInfoRet{
			LastTime:      last.LastTime,
			LastSignature: last.LastSignature,
			LastSignedInfo: onlinePb.SignedInfo{
				Peer:        last.LastSignedInfo.Peer,
				CreatedTime: last.LastSignedInfo.CreatedTime,
				Version:     last.LastSignedInfo.Version,
				Nonce:       last.LastSignedInfo.Nonce,
				BttcAddress: last.LastSignedInfo.BttcAddress,
				SignedTime:  last.LastSignedInfo.SignedTime,
			},
		}

		return cmds.EmitOnce(res, &r)
	},
	Type: chain.LastOnlineInfoRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *chain.LastOnlineInfoRet) error {
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

var StatusConfigCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "(old)get reporting status-contract config. ",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		rs, err := chain.GetReportStatus()
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, rs)
	},
	Type: chain.ReportStatusInfo{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *chain.ReportStatusInfo) error {
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

var ReportOnlineServerCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "report online server. ",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}

		spin.DC.SendDataOnline(node, cfg)

		return cmds.EmitOnce(res, "report online server ok!")
	},
}

var ReportStatusContractCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "(old,drop it)report status-contract. ",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := reportstatus.CmdReportStatus()
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, "report status contract ok!")
	},
}
