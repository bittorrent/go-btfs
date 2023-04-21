package offline

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/utils"

	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"

	cmds "github.com/bittorrent/go-btfs-cmds"
	renterpb "github.com/bittorrent/go-btfs/protos/renter"

	cmap "github.com/orcaman/concurrent-map"
)

var StorageUploadSignCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Return the signed data to the upload session.",
		ShortDescription: `
This command returns the signed data (From BTFS SDK application's perspective)
to the upload session.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("session-id", true, false, "ID for the entire storage upload session."),
		cmds.StringArg("peer-id", true, false, "Offline signs needed for this particular client."),
		cmds.StringArg("nonce-timestamp", true, false, "Nonce timestamp string for this upload signing."),
		cmds.StringArg("upload-session-signature", true, false, "Private key-signed string of peer-id:nonce-timestamp"),
		cmds.StringArg("session-status", true, false, "current upload session status."),
		cmds.StringArg("signed", true, false, "signed json data."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		ssId := req.Arguments[0]
		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}
		rss, err := sessions.GetRenterSession(ctxParams, ssId, "", make([]string, 0))
		if err != nil {
			return err
		}
		err = verifyReceivedMessage(req, rss)
		if err != nil {
			return err
		}
		status, err := rss.Status()
		if err != nil {
			return err
		}
		if status.Status != req.Arguments[4] {
			return fmt.Errorf("error status, want: %s, actual: %s", status.Status, req.Arguments[4])
		}
		bytes, err := helper.StringToBytes(req.Arguments[5], helper.Base64)
		if err != nil {
			return err
		}
		var cm cmap.ConcurrentMap
		switch req.Arguments[4] {
		case sessions.RssSubmitStatus:
			cm = uh.BalanceChanMaps
		//case sessions.RssSubmitBalanceReqSignedStatus:
		//	cm = uh.SignedChannelCommitChanMaps
		case sessions.RssGuardStatus:
			cm = uh.FileMetaChanMaps
		case sessions.RssGuardFileMetaSignedStatus:
			cm = uh.QuestionsChanMaps
		case sessions.RssWaitUploadStatus:
			cm = uh.WaitUploadChanMap
		//case sessions.RssPayStatus:
		//cm = uh.PayinReqChanMaps
		default:
			return errors.New("wrong status:" + req.Arguments[4])
		}
		if bc, ok := cm.Get(ssId); ok {
			bc.(chan []byte) <- bytes
		}
		return rss.SaveOfflineSigning(&renterpb.OfflineSigning{
			Sig: bytes,
		})
	},
}
