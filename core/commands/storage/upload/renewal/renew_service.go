package renewal

import (
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/utils"

	cmds "github.com/bittorrent/go-btfs-cmds"
)

// StorageRenewServiceCmd manages the auto-renewal service
var StorageRenewServiceCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Manage the auto-renewal service for storage files.",
		ShortDescription: `
This command allows users to control the auto-renewal service that automatically
renews storage contracts before they expire.

Examples:
    # Check service status
    $ btfs storage renew service status

    # Start the service
    $ btfs storage renew service start

    # Stop the service
    $ btfs storage renew service stop

    # Restart the service
    $ btfs storage renew service restart
`,
	},
	Subcommands: map[string]*cmds.Command{
		"status":  StorageRenewServiceStatusCmd,
		"start":   StorageRenewServiceStartCmd,
		"stop":    StorageRenewServiceStopCmd,
		"restart": StorageRenewServiceRestartCmd,
	},
}

// StorageRenewServiceStatusCmd checks the status of the auto-renewal service
var StorageRenewServiceStatusCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Check the status of the auto-renewal service.",
		ShortDescription: `
This command shows the current status of the auto-renewal service,
including whether it's running and its configuration.

Example:
    $ btfs storage renew service status
`,
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		_, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		err = utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		status := GetRenewalServiceStatus()
		return res.Emit(status)
	},
	Type: RenewalServiceStatus{},
}

// StorageRenewServiceStartCmd starts the auto-renewal service
var StorageRenewServiceStartCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Start the auto-renewal service.",
		ShortDescription: `
This command starts the auto-renewal service if it's not already running.
The service will automatically check for files that need renewal and
process them according to their auto-renewal configuration.

Example:
    $ btfs storage renew service start
`,
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		_, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		err = utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		// Get context parameters
		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}

		// Start the service
		err = InitializeRenewalService(ctxParams)
		if err != nil {
			return err
		}

		response := map[string]interface{}{
			"success": true,
			"message": "Auto-renewal service started successfully",
		}

		return res.Emit(response)
	},
}

// StorageRenewServiceStopCmd stops the auto-renewal service
var StorageRenewServiceStopCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Stop the auto-renewal service.",
		ShortDescription: `
This command stops the auto-renewal service if it's currently running.
No automatic renewals will be processed while the service is stopped.

Example:
    $ btfs storage renew service stop
`,
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		_, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		err = utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		// Stop the service
		err = ShutdownRenewalService()
		if err != nil {
			return err
		}

		response := map[string]interface{}{
			"success": true,
			"message": "Auto-renewal service stopped successfully",
		}

		return res.Emit(response)
	},
}

// StorageRenewServiceRestartCmd restarts the auto-renewal service
var StorageRenewServiceRestartCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Restart the auto-renewal service.",
		ShortDescription: `
This command restarts the auto-renewal service, which is useful for
applying configuration changes or recovering from errors.

Example:
    $ btfs storage renew service restart
`,
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		_, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		err = utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		// Get context parameters
		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}

		// Restart the service
		manager := GetGlobalRenewalServiceManager()
		err = manager.RestartService(ctxParams)
		if err != nil {
			return err
		}

		response := map[string]interface{}{
			"success": true,
			"message": "Auto-renewal service restarted successfully",
		}

		return res.Emit(response)
	},
}
