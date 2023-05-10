package commands

import (
	"context"
	"fmt"
	"github.com/bittorrent/go-btfs/utils"
	"io"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

const (
	CheckBackoffDuration = 20 * time.Second
	CheckMaxRetries      = 3
)

type NetworkRet struct {
	CodeBttc   int    `json:"code_bttc"`
	ErrBttc    string `json:"err_bttc"`
	CodeStatus int    `json:"code_status"`
	ErrStatus  string `json:"err_status"`
}

var NetworkCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get btfs network information",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		timeoutCtx, _ := context.WithTimeout(context.Background(), CheckBackoffDuration*time.Duration(CheckMaxRetries))
		_, err = chain.ChainObject.Backend.BlockNumber(timeoutCtx)
		if err != nil {
			chain.CodeBttc = chain.ConstCodeError
			chain.ErrBttc = err
		} else {
			chain.CodeBttc = chain.ConstCodeSuccess
			chain.ErrBttc = nil
		}

		//chain.ErrStatus = errors.New("network111")
		ret := NetworkRet{
			CodeBttc:   chain.CodeBttc,
			ErrBttc:    switchErrToString(chain.ErrBttc),
			CodeStatus: chain.CodeStatus,
			ErrStatus:  switchErrToString(chain.ErrStatus),
		}
		return cmds.EmitOnce(res, &ret)
	},
	Type: NetworkRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *NetworkRet) error {
			_, err := fmt.Fprintf(w, "code bttc:\t%d\nerr bttc:\t%s\ncode status:\t%d\nerr status:\t%s\n",
				out.CodeBttc, out.ErrBttc, out.CodeStatus, out.ErrStatus)
			return err
		}),
	},
}

func switchErrToString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
