package upload

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/utils"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

// RenewStatusResponse represents the status of a renewal operation
type RenewStatusResponse struct {
	SessionID string    `json:"session_id"`
	FileHash  string    `json:"file_hash"`
	Status    string    `json:"status"`
	Duration  int       `json:"duration"`
	TotalCost int64     `json:"total_cost"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Message   string    `json:"message"`
}

// RenewListResponse represents a list of renewals
type RenewListResponse struct {
	Renewals []RenewStatusResponse `json:"renewals"`
	Total    int                   `json:"total"`
}

// StorageRenewStatusCmd checks the status of a specific renewal
var StorageRenewStatusCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Check the status of a storage renewal operation.",
		ShortDescription: `
This command checks the status of a specific storage renewal operation
using the renewal session ID.

Example:
    $ btfs storage renew status <renewal-session-id>
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("renewal-session-id", true, false, "ID of the renewal session to check."),
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

		renewalSessionID := req.Arguments[0]

		// Get context parameters
		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}

		// Get renewal information
		renewalInfo, err := getRenewalInfo(ctxParams, renewalSessionID)
		if err != nil {
			return fmt.Errorf("failed to get renewal info: %v", err)
		}

		if renewalInfo == nil {
			return fmt.Errorf("renewal session not found: %s", renewalSessionID)
		}

		// Create status response
		status := &RenewStatusResponse{
			SessionID: renewalSessionID,
			FileHash:  renewalInfo.FileHash,
			Status:    "completed", // TODO: Implement actual status tracking
			Duration:  renewalInfo.Duration,
			TotalCost: renewalInfo.TotalCost,
			CreatedAt: time.Now(), // TODO: Store actual creation time
			ExpiresAt: renewalInfo.NewEnd,
			Message:   fmt.Sprintf("Renewal for file %s is active", renewalInfo.FileHash),
		}

		return res.Emit(status)
	},
	Type: RenewStatusResponse{},
}

// StorageRenewListCmd lists all renewals for the current node
var StorageRenewListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List all storage renewals for the current node.",
		ShortDescription: `
This command lists all storage renewal operations performed by the current node.

Example:
    $ btfs storage renew list
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

		// Get all renewals
		renewals, err := getAllRenewals(ctxParams)
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

// getRenewalInfo retrieves renewal information from datastore
func getRenewalInfo(ctxParams *uh.ContextParams, sessionID string) (*RenewRequest, error) {
	renewalKey := fmt.Sprintf("/btfs/%s/renewals/%s", ctxParams.N.Identity.String(), sessionID)

	data, err := ctxParams.N.Repo.Datastore().Get(ctxParams.Ctx, datastore.NewKey(renewalKey))
	if err != nil {
		if err == datastore.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	var renewalInfo RenewRequest
	err = json.Unmarshal(data, &renewalInfo)
	if err != nil {
		return nil, err
	}

	return &renewalInfo, nil
}

// getAllRenewals retrieves all renewal information for the current node
func getAllRenewals(ctxParams *uh.ContextParams) ([]RenewStatusResponse, error) {
	prefix := fmt.Sprintf("/btfs/%s/renewals/", ctxParams.N.Identity.String())
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

		var renewalInfo RenewRequest
		if err := json.Unmarshal(result.Value, &renewalInfo); err != nil {
			continue
		}

		// Extract session ID from key
		sessionID := result.Key[len(prefix):]

		status := RenewStatusResponse{
			SessionID: sessionID,
			FileHash:  renewalInfo.FileHash,
			Status:    "completed", // TODO: Implement actual status tracking
			Duration:  renewalInfo.Duration,
			TotalCost: renewalInfo.TotalCost,
			CreatedAt: time.Now(), // TODO: Store actual creation time
			ExpiresAt: renewalInfo.NewEnd,
			Message:   fmt.Sprintf("Renewal for file %s is active", renewalInfo.FileHash),
		}

		renewals = append(renewals, status)
	}

	return renewals, nil
}
