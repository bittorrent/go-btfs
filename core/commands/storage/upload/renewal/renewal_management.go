package renewal

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/utils"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

const (
	filterOptionName = "filter"
)

const (
	RenewTypeAll    = "all"
	RenewTypeAuto   = "auto"
	RenewTypeManual = "manual"
)

var (
	renewKeyPrefix = "/btfs/%s/renew/"
	autoRenewKey   = renewKeyPrefix + "auto"
	manualRenewKey = renewKeyPrefix + "manual"
)

// RenewStatusResponse represents the status of a renewal operation
type RenewStatusResponse struct {
	FileHash  string    `json:"file_hash"`
	Duration  int       `json:"duration"`
	TotalCost int64     `json:"total_cost"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RenewListResponse represents a list of renewals
type RenewListResponse struct {
	Renewals []RenewStatusResponse `json:"renewals"`
	Total    int                   `json:"total"`
}

// StorageRenewInfoCmd checks the status of a specific renewal
var StorageRenewInfoCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Check the status of a storage renewal for a specific CID.",
		ShortDescription: `
This command checks the status of a storage renewal for a specific CID.

Example:
    $ btfs storage upload renew info <cid>
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("cid", true, false, "CID of the renewal file to check."),
	},
	Options: []cmds.Option{
		cmds.StringOption(filterOptionName, "-f", "Filter renewals by type [auto|manual]").WithDefault(RenewTypeAuto),
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

		renewalCID := req.Arguments[0]

		// Get context parameters
		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}

		// Get renewal information
		renewalInfo, err := getRenewalInfo(ctxParams, renewalCID, req.Options[filterOptionName].(string))
		if err != nil {
			return fmt.Errorf("failed to get renewal info: %v", err)
		}

		if renewalInfo == nil {
			return fmt.Errorf("renewal cid not found: %s", renewalCID)
		}

		return res.Emit(renewalInfo)
	},
	Type: RenewalInfo{},
}

// StorageRenewListCmd lists all renewals for the current node
var StorageRenewListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List all storage renewals for the current node.",
		ShortDescription: `
This command lists all storage renewal operations performed by the current node.

Example:
    $ btfs storage upload renew list
`,
	},
	Options: []cmds.Option{
		cmds.StringOption(filterOptionName, "-f", "Filter renewals by type [all|auto|manual]").WithDefault("all"),
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

		// Get all renewals
		renewals, err := getRenewalsFiles(ctxParams, req.Options[filterOptionName].(string))
		if err != nil {
			return fmt.Errorf("failed to get renewals: %v", err)
		}

		response := &RenewListResponse{
			Renewals: renewals,
			Total:    len(renewals),
		}

		return res.Emit(response)
	},
	Type: RenewListResponse{},
}

var StorageRenewEnableCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Enable storage renewals for a specific CID.",
		ShortDescription: `
This command enables storage renewals for a specific CID.

Example:
    $ btfs storage upload renew enable <cid>
`,
	},
	Type: Res{},
	Arguments: []cmds.Argument{
		cmds.StringArg("cid", true, false, "CID of the file to enable renewals for."),
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
		return enableAutoRenewal(ctxParams, req.Arguments[0])
	},
}

var StorageRenewDisableCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Disable storage renewals for a specific CID.",
		ShortDescription: `
This command disables storage renewals for a specific CID.

Example:
    $ btfs storage upload renew disable <cid>
`,
	},
	Type: Res{},
	Arguments: []cmds.Argument{
		cmds.StringArg("cid", true, false, "CID of the file to disable renewals for."),
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
		return disableAutoRenewal(ctxParams, req.Arguments[0])
	},
}

// StoreRenewalInfo stores auto-renewal configuration for a file
func StoreRenewalInfo(ctxParams *uh.ContextParams, info *RenewalInfo, renewType string) error {

	if renewType != RenewTypeAuto && renewType != RenewTypeManual {
		return fmt.Errorf("invalid filter type: %s", renewType)
	}
	configKey := fmt.Sprintf(autoRenewKey+"/%s", ctxParams.N.Identity.String(), info.CID)
	if renewType == RenewTypeManual {
		configKey = fmt.Sprintf(manualRenewKey+"/%s", ctxParams.N.Identity.String(), info.CID)
	}

	configData, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return ctxParams.N.Repo.Datastore().Put(ctxParams.Ctx, datastore.NewKey(configKey), configData)
}

// getRenewalInfo retrieves renewal information from datastore
func getRenewalInfo(ctxParams *uh.ContextParams, cid string, renewType string) (*RenewalInfo, error) {
	if renewType != RenewTypeAuto && renewType != RenewTypeManual {
		return nil, fmt.Errorf("invalid filter type: %s", renewType)
	}
	if cid == "" {
		return nil, fmt.Errorf("cid cannot be empty")
	}
	renewalKey := fmt.Sprintf(autoRenewKey+"/%s", ctxParams.N.Identity.String(), cid)
	if renewType == RenewTypeManual {
		renewalKey = fmt.Sprintf(manualRenewKey+"/%s", ctxParams.N.Identity.String(), cid)
	}

	data, err := ctxParams.N.Repo.Datastore().Get(ctxParams.Ctx, datastore.NewKey(renewalKey))
	if err != nil {
		if errors.Is(err, datastore.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var renewalInfo RenewalInfo
	err = json.Unmarshal(data, &renewalInfo)
	if err != nil {
		return nil, err
	}

	return &renewalInfo, nil
}

// getRenewalsFiles retrieves all renewal information for the current node
func getRenewalsFiles(ctxParams *uh.ContextParams, filterType string) ([]RenewStatusResponse, error) {
	if filterType != RenewTypeAll && filterType != RenewTypeAuto && filterType != RenewTypeManual {
		return nil, fmt.Errorf("invalid filter type: %s", filterType)
	}

	prefix := fmt.Sprintf(renewKeyPrefix, ctxParams.N.Identity.String())
	if filterType == RenewTypeAuto {
		prefix = fmt.Sprintf(autoRenewKey, ctxParams.N.Identity.String())
	}
	if filterType == RenewTypeManual {
		prefix = fmt.Sprintf(manualRenewKey, ctxParams.N.Identity.String())
	}

	q := query.Query{
		Prefix: prefix,
	}

	results, err := ctxParams.N.Repo.Datastore().Query(ctxParams.Ctx, q)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var renewals []RenewStatusResponse

	for result := range results.Next() {
		if result.Error != nil {
			continue
		}

		var renewalInfo RenewalInfo
		if err := json.Unmarshal(result.Value, &renewalInfo); err != nil {
			continue
		}

		status := RenewStatusResponse{
			FileHash:  renewalInfo.CID,
			Duration:  renewalInfo.RenewalDuration,
			TotalCost: renewalInfo.TotalPay,
			CreatedAt: renewalInfo.CreatedAt,
			ExpiresAt: renewalInfo.CreatedAt.Add(time.Duration(renewalInfo.RenewalDuration) * 24 * time.Hour),
		}

		renewals = append(renewals, status)
	}

	return renewals, nil
}

func enableAutoRenewal(ctxParams *uh.ContextParams, fileHash string) error {
	return EnableAutoRenewalForFile(ctxParams, fileHash)
}

func disableAutoRenewal(ctxParams *uh.ContextParams, fileHash string) error {
	return DisableAutoRenewalForFile(ctxParams, fileHash)
}

type Res struct {
	ID string
}
