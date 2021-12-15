package spin

import (
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/upload"

	cmds "github.com/TRON-US/go-btfs-cmds"

	"go4.org/syncutil"
)

func RenterSessions(req *cmds.Request, env cmds.Environment) {
	go func() {
		params, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return
		}

		cursor, err := sessions.GetRenterSessionsCursor(params)
		if err != nil {
			return
		}
		// Limit to 10 at a time to lower resource consumption
		sem := syncutil.NewSem(10)
		for {
			session, err := cursor.NextSession(sessions.RssWaitUploadReqSignedStatus)
			if err != nil {
				break
			}
			if session == nil {
				break
			}
			sem.Acquire(1)
			go func(session *sessions.RenterSession) {
				defer sem.Release(1)
				if e := upload.ResumeWaitUploadOnSigning(session); e != nil {
					return
				}
			}(session)
		}
	}()
}
