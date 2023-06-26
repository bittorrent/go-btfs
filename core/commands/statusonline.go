package commands

import (
	"encoding/json"
	"fmt"
	cmds "github.com/bittorrent/go-btfs-cmds"
	onlinePb "github.com/bittorrent/go-btfs-common/protos/online"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/spin"
	"github.com/bittorrent/go-btfs/utils"
	"io"
	"strconv"
	"time"
)

// ReportOnlineDailyCmd (report online daily)
var ReportOnlineDailyCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "daily report online server. ",
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

		msg, err := spin.DC.SendOnlineDaily(node, cfg)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, "daily report online server ok! "+msg)
	},
}

type RetReportOnlineListDaily struct {
	Records  []*chain.LastOnlineInfoRet `json:"records"`
	Total    int                        `json:"total"`
	PeerId   string                     `json:"peer_id"`
	BttcAddr string                     `json:"bttc_addr"`
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
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		peerId := n.Identity.Pretty()

		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}
		bttcAddr := cfg.Identity.BttcAddr

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

		rs := make([]*chain.LastOnlineInfoRet, 0)
		for _, v := range list {
			r := chain.LastOnlineInfoRet{
				LastTime:      v.LastTime,
				LastSignature: v.LastSignature,
				LastSignedInfo: onlinePb.SignedInfo{
					Peer:        v.LastSignedInfo.Peer,
					CreatedTime: v.LastSignedInfo.CreatedTime,
					Version:     v.LastSignedInfo.Version,
					Nonce:       v.LastSignedInfo.Nonce,
					BttcAddress: v.LastSignedInfo.BttcAddress,
					SignedTime:  v.LastSignedInfo.SignedTime,
				},
			}
			rs = append(rs, &r)
		}

		return cmds.EmitOnce(res, &RetReportOnlineListDaily{
			Records:  rs,
			Total:    total,
			PeerId:   peerId,
			BttcAddr: bttcAddr,
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
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

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
			StatusContract: chain.ChainObject.Chainconfig.StatusAddress.String(),
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
	EveryDaySeconds int64     `json:"every_day_seconds"`
	LastTime        time.Time `json:"last_time"`
}

// ReportLastTimeDailyCmd (last time daily)
var ReportLastTimeDailyCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "report status-contract total info, (total count, total gas spend, and contract address)",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		last, err := chain.GetReportOnlineLastTimeDaily()
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &RetReportLastTime{
			EveryDaySeconds: last.EveryDaySeconds,
			LastTime:        last.LastReportTime,
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
