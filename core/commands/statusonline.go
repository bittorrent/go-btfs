package commands

import (
	"encoding/json"
	"fmt"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/spin"
	"io"
	"strconv"
	"time"
)

// ReportOnlineDailyCmd (report online daily)
var ReportOnlineDailyCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "report online server. ",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}

		spin.DC.SendOnlineDaily(node, cfg)

		return cmds.EmitOnce(res, "report online server ok!")
	},
}

type RetReportOnlineListDaily struct {
	Records []*chain.LastOnlineInfo `json:"records"`
	Total   int                     `json:"total"`
	PeerId  string                  `json:"peer_id"`
}

// ReportListDailyCmd (report list daily)
var ReportListDailyCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "report list daily, and input from and limit to get its.",
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

		list, err := chain.GetReportOnlineListDailyOK()
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

		return cmds.EmitOnce(res, &RetReportOnlineListDaily{
			Records: list,
			Total:   total,
			PeerId:  peerId,
		})
	},
	Type: RetReportOnlineListDaily{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *RetReportOnlineListDaily) error {
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

type RetTotalDaily struct {
	PeerId         string `json:"peer_id"`
	StatusContract string `json:"status_contract"`
	TotalCount     int    `json:"total_count"`
}

// TotalDailyCmd (total daily)
var TotalDailyCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "report status-contract total info, (total count, total gas spend, and contract address)",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		peerId := n.Identity.Pretty()

		list, err := chain.GetReportOnlineListDailyOK()
		if err != nil {
			return err
		}
		if list == nil {
			return nil
		}

		return cmds.EmitOnce(res, &RetTotalDaily{
			PeerId:         peerId,
			StatusContract: "0x111",
			TotalCount:     len(list),
		})
	},
	Type: RetTotalDaily{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *RetTotalDaily) error {
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

type RetReportLastTime struct {
	LastTime time.Time `json:"last_time"`
}

// ReportLastTimeDailyCmd (last time daily)
var ReportLastTimeDailyCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "report status-contract total info, (total count, total gas spend, and contract address)",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		last, err := chain.GetReportOnlineLastTimeDaily()
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &RetReportLastTime{
			LastTime: last.LastReportTime,
		})
	},
	Type: RetReportLastTime{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *RetReportLastTime) error {
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
