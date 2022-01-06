package stake

import (
	"fmt"
	"io"
	"math/big"
	"time"

	cmds "github.com/TRON-US/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

type StakeInfoCmdRet struct {
	Amount   *big.Int `json:"amount"`
	LockTime *big.Int `json:"locktime"`
}

var StakeInfoCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get stake status.",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		amount := chain.SettleObject.StakeService.CurrentStakeAmount()
		locktime := chain.SettleObject.StakeService.CurrentStakeLockTime()

		return cmds.EmitOnce(res, &StakeInfoCmdRet{
			Amount:   amount,
			LockTime: locktime,
		})
	},
	Type: &StakeInfoCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *StakeInfoCmdRet) error {
			_, err := fmt.Fprintf(w, "the amount of staking is : %s, the lock time of unstake is to : %v \n", out.Amount, out.LockTime)
			return err
		}),
	},
}
