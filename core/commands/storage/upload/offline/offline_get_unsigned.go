package offline

import (
	"errors"
	"github.com/bittorrent/go-btfs/utils"

	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"

	cmds "github.com/bittorrent/go-btfs-cmds"
)

var StorageUploadGetUnsignedCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get the input data for upload signing.",
		ShortDescription: `
This command obtains the upload signing input data for from the upload session 
(From BTFS SDK application's perspective) and returns to the caller.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("session-id", true, false, "ID for the entire storage upload session."),
		cmds.StringArg("peer-id", true, false, "Offline signs needed for this particular client."),
		cmds.StringArg("nonce-timestamp", true, false, "Nonce timestamp string for this upload signing."),
		cmds.StringArg("upload-session-signature", true, false, "Private key-signed string of peer-id:nonce-timestamp"),
		cmds.StringArg("session-status", true, false, "current upload session status."),
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
		rssStatus, err := rss.Status()
		if err != nil {
			return err
		}
		if req.Arguments[4] != rssStatus.Status {
			return errors.New("unexpected session status from SDK during communication in upload signing")
		}
		signing, err := rss.OfflineSigning()
		if err != nil {
			return err
		}

		rawString, err := helper.BytesToString(signing.Raw, helper.Base64)
		if err != nil {
			return err
		}
		return res.Emit(&GetUnsignedRes{Unsigned: rawString, Price: signing.Price})
	},
	Type: GetUnsignedRes{},
}

type GetUnsignedRes struct {
	Opcode   string
	Unsigned string
	Price    int64
}
