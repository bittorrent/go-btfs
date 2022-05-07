package commands

import (
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/bittorrent/go-btfs/core/commands/storage/path"

	"github.com/bittorrent/go-btfs-cmds"
	"github.com/cenkalti/backoff/v4"
)

var daemonStartup = func() *backoff.ExponentialBackOff {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 2 * time.Second
	bo.MaxElapsedTime = 300 * time.Second
	bo.Multiplier = 1
	bo.MaxInterval = 2 * time.Second
	return bo
}()

var (
	shutdownReady = make(chan struct{}) // used to make sure the daemon shutdown be ready
	restartLock   sync.Mutex            // used to stop the daemon exit imediatly
)

// NotifyAndWaitIfOnRestaring will called by the daemon goroutine
func NotifyAndWaitIfOnRestarting() {
	close(shutdownReady)
	restartLock.Lock()
}

const (
	postPathModificationName = "post-path-modification"
)

// When restart the daemon:
// 1. Restarting goroutine set a lock, this lock will stop the daemon goroutine exit
// 2. Restarting goroutine send shutdown command to the daemon goroutine, and then wait the daemon's  shutdown state ready
// 3. When the daemon goroutine accepted the shutdown command, it will close the node and shutdown all the other services
// 4. When the daemon shutdown is ready or timeout, it will notify the restarting goroutine
// 5. The daemon goroutine then wait the restating gorouting by waitting to set the restarting lock
// (If the daemon is not on restarting, the daemon will set the lock immediatly, and exit directly)
// 6. Restarting goroutine will complete it's work and exec start daemon command
// 7. Restarting goroutine exit the process
var restartCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Restart the daemon.",
		ShortDescription: `
Shutdown the runnning daemon and start a new daemon process.
And if specified a new btfs path, it will be applied.
`,
	},
	Options: []cmds.Option{
		cmds.BoolOption(postPathModificationName, "p", "post path modification").WithDefault(false),
	}, Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) (err error) {
		go func() {
			restartLock.Lock()
			defer restartLock.Unlock()
			shutdownCmd := exec.Command(path.Executable, "shutdown")
			if err = shutdownCmd.Run(); err != nil {
				return
			}
			<-shutdownReady
			defer func() {
				daemonCmd := exec.Command(path.Executable, "daemon")
				err = daemonCmd.Start()
				os.Exit(0)
			}()
			if req.Options[postPathModificationName].(bool) && path.StorePath != "" && path.OriginPath != "" {
				if err = path.MoveFolder(); err != nil {
					return
				}
				if err = path.WriteProperties(); err != nil {
					return
				}
			}
		}()
		return nil
	},
}
