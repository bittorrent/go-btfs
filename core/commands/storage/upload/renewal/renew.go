package renewal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	nodepb "github.com/bittorrent/go-btfs-common/protos/node"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	"github.com/bittorrent/go-btfs/utils"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p/core/peer"
)

var log = logging.Logger("renew")

const (
	renewDurationOptionName = "duration"
	renewTokenOptionName    = "renew-token"
	renewPriceOptionName    = "renew-price"
)

// RenewRequest represents a file renewal request
type RenewRequest struct {
	CID         string         `json:"cid"`
	Token       common.Address `json:"token"`
	Price       uint64         `json:"price"`
	Duration    int            `json:"duration"`
	SpId        string         `json:"sp_id"`
	RenterID    peer.ID        `json:"renter_id"`
	ShardId     string         `json:"shard_id"`
	ShardSize   int64          `json:"shard_size"`
	ContractId  string         `json:"contract_id"`
	OriginalEnd time.Time      `json:"original_end"`
	NewEnd      time.Time      `json:"new_end"`
	TotalCost   int64          `json:"total_cost"`
}

// RenewResponse represents the response of a renewal operation
type RenewResponse struct {
	Success       bool      `json:"success"`
	CID           string    `json:"cid"`
	NewExpiration time.Time `json:"new_expiration"`
	TotalCost     int64     `json:"total_cost"`
}

// StorageRenewCmd implements the storage renew command
var StorageRenewCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Renew storage duration for uploaded files but not set auto-renewed.",
		ShortDescription: `
This command allows users to extend the storage duration of previously uploaded files
without re-uploading the content. The renewal extends the storage contract with
existing storage providers.

Examples:
    # Renew a file for 30 days
    $ btfs storage upload renew <file-hash> --duration 30

    # Renew with specific token and price
    $ btfs storage upload renew <file-hash> --duration 60 --token WBTT --price 1000
`,
	},
	Subcommands: map[string]*cmds.Command{
		"enable":  StorageRenewEnableCmd,
		"disable": StorageRenewDisableCmd,
		"info":    StorageRenewInfoCmd,
		"list":    StorageRenewListCmd,
		"service": StorageRenewServiceCmd,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("cid", true, false, "cid of the file to renew."),
	},
	Options: []cmds.Option{
		cmds.IntOption(renewDurationOptionName, "d", "Renewal duration in days.").WithDefault(30),
		cmds.StringOption(renewTokenOptionName, "rt", "Token type for payment (WBTT/TRX/USDD/USDT).").WithDefault("WBTT"),
		cmds.Int64Option(renewPriceOptionName, "rp", "Max price per GiB per day in µBTT."),
	},
	RunTimeout: 10 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		if !nd.IsOnline {
			return errors.New("node must be online to renew storage")
		}

		err = utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		// Extract parameters
		cid := req.Arguments[0]
		duration, _ := req.Options[renewDurationOptionName].(int)
		tokenStr, _ := req.Options[renewTokenOptionName].(string)
		priceOpt, hasPriceOpt := req.Options[renewPriceOptionName].(int64)

		// Validate parameters
		if duration <= 0 {
			return errors.New("duration must be positive")
		}

		// Get context parameters
		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}

		// check if the cid enabled autorenew
		info, err := getRenewalInfo(ctxParams, cid, RenewTypeAuto)
		if err != nil {
			return err
		}
		if info != nil && info.Enabled {
			return fmt.Errorf("file %s is already auto-renewed and cannot be renewed manually", cid)
		}
		// Get token address
		token, exists := tokencfg.MpTokenAddr[tokenStr]
		if !exists {
			return fmt.Errorf("invalid token type: %s", tokenStr)
		}

		// Get current price if not specified
		var price int64
		if hasPriceOpt {
			price = priceOpt
		} else {
			priceObj, err := chain.SettleObject.OracleService.CurrentPrice(token)
			if err != nil {
				return err
			}
			price = priceObj.Int64()
		}

		contracts, err := sessions.ListShardsContracts(ctxParams.N.Repo.Datastore(), ctxParams.N.Identity.String(), nodepb.ContractStat_RENTER.String())
		if err != nil {
			return fmt.Errorf("failed to get shard contract, you can sync first, then try it again")
		}

		var totalCost int64
		for _, c := range contracts {
			fileHash, err := ctxParams.N.Repo.Datastore().Get(ctxParams.Ctx, datastore.NewKey(
				fmt.Sprintf(userFileShard, ctxParams.N.Identity, c.Meta.ContractId)))
			if err != nil {
				return fmt.Errorf("failed to get file hash for contract %s: %v", c.Meta.ContractId, err)
			}
			if cid != string(fileHash) {
				continue
			}

			renewReq := &RenewRequest{
				CID:         cid,
				Token:       token,
				Price:       uint64(price),
				Duration:    duration,
				SpId:        c.Meta.SpId,
				RenterID:    ctxParams.N.Identity,
				ShardId:     c.Meta.ShardHash,
				ShardSize:   int64(c.Meta.ShardSize),
				ContractId:  c.Meta.ContractId,
				OriginalEnd: time.Unix(int64(c.Meta.StorageEnd), 0),
				NewEnd:      time.Unix(int64(c.Meta.StorageEnd), 0).Add(time.Duration(duration) * 24 * time.Hour), // This should be calculated from existing contract
			}

			resp, err := executeRenewal(ctxParams, renewReq)
			if err != nil {
				return fmt.Errorf("renewal failed: %v", err)
			}

			totalCost += resp.TotalCost

		}

		// TODO fill the field
		info = &RenewalInfo{
			CID: cid,
		}
		StoreRenewalInfo(ctxParams, info, RenewTypeManual)

		return res.Emit(RenewResponse{
			Success:       true,
			CID:           cid,
			NewExpiration: time.Now(),
			TotalCost:     totalCost,
		})
	},
	Type: RenewResponse{},
}

// executeRenewal performs the actual renewal operation
func executeRenewal(ctxParams *uh.ContextParams, renewReq *RenewRequest) (*RenewResponse, error) {
	// Calculate total cost
	rate, err := chain.SettleObject.OracleService.CurrentRate(renewReq.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to get token rate: %v", err)
	}

	totalCost, err := uh.TotalPay(renewReq.ShardSize, int64(renewReq.Price), renewReq.Duration, rate)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate total cost: %v", err)
	}

	renewReq.TotalCost = totalCost * rate.Int64()
	// Execute renewal with storage providers
	err = payRenewalCheque(ctxParams, renewReq, totalCost)
	if err != nil {
		return nil, fmt.Errorf("failed to execute renewal with providers: %v", err)
	}

	err = extendShardEndTime(ctxParams, renewReq.ContractId, renewReq.Duration)
	if err != nil {
		return nil, fmt.Errorf("failed to extend shard end time: %v", err)
	}

	return &RenewResponse{
		Success:       true,
		CID:           renewReq.CID,
		NewExpiration: renewReq.NewEnd,
		TotalCost:     totalCost,
	}, nil
}

// payRenewalCheque pays renewal fee directly via cheque to storage provider
func payRenewalCheque(ctxParams *uh.ContextParams, renewReq *RenewRequest, paymentAmount int64) error {
	// Get the original shard contract to find the storage provider
	// Get storage provider ID
	spId := renewReq.SpId
	if spId == "" {
		return fmt.Errorf("no storage provider ID found in contract")
	}

	fmt.Printf("Paying renewal cheque for shard %s to sp %s, amount: %d\n", renewReq.ShardId, spId, paymentAmount)

	// Check available balance before issuing cheque
	err := checkAvailableBalance(ctxParams.Ctx, paymentAmount, renewReq.Token)
	if err != nil {
		return fmt.Errorf("insufficient balance for renewal payment: %v", err)
	}

	// Issue cheque directly to storage provider for renewal
	err = issueRenewalCheque(ctxParams, spId, paymentAmount, renewReq.Token, renewReq.ShardId, renewReq.Duration, renewReq.ContractId)
	if err != nil {
		return fmt.Errorf("failed to issue renewal cheque to provider %s: %v", spId, err)
	}

	return nil
}

// issueRenewalCheque issues a cheque directly to storage provider for renewal payment
func issueRenewalCheque(ctxParams *uh.ContextParams, providerID string, amount int64, token common.Address, shardHash string, duration int, contractId string) error {
	fmt.Printf("Issuing renewal cheque to provider %s for shard %s, amount: %d, duration: %d days\n", providerID, shardHash, amount, duration)

	// Get settlement service
	if chain.SettleObject.SwapService == nil {
		return fmt.Errorf("settlement service not available")
	}

	// Convert amount to big.Int
	paymentAmount := big.NewInt(amount)

	// Generate a renewal contract ID for tracking
	renewalContractID := fmt.Sprintf("renewal_%d_%s", duration, contractId)

	// Issue cheque through settlement service
	// This directly pays the provider without creating a new contract
	// Pay method signature: Pay(ctx context.Context, peer string, amount *big.Int, contractId string, token common.Address)
	err := chain.SettleObject.SwapService.Settle(providerID, paymentAmount, renewalContractID, token)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully issued renewal cheque to provider %s for contract %s\n", providerID, renewalContractID)

	// Store cheque information for tracking
	err = storeRenewalChequeInfo(ctxParams, providerID, shardHash, renewalContractID, amount, duration)
	if err != nil {
		log.Errorf("Failed to store renewal cheque info: %v", err)
		// Don't fail the payment for storage issues
	}

	return nil
}

// storeRenewalChequeInfo stores cheque information for renewal tracking
func storeRenewalChequeInfo(ctxParams *uh.ContextParams, providerID, shardHash string, contractID string, amount int64, duration int) error {
	chequeInfo := map[string]interface{}{
		"provider_id": providerID,
		"shard_hash":  shardHash,
		"amount":      amount,
		"duration":    duration,
		"issued_at":   time.Now().Unix(),
		"contract_id": contractID, // Store the renewal contract ID
	}

	chequeKey := fmt.Sprintf("/btfs/%s/renewal_cheques/%s_%s_%d",
		ctxParams.N.Identity.String(), providerID, shardHash, time.Now().Unix())

	chequeData, err := json.Marshal(chequeInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal cheque info: %v", err)
	}

	return ctxParams.N.Repo.Datastore().Put(ctxParams.Ctx, datastore.NewKey(chequeKey), chequeData)
}

// Note: We no longer use contract-based renewal. Instead, we directly issue cheques to providers.
// This simplifies the renewal process and avoids the complexity of contract renegotiation.

// RenewalShardInfo represents auto-renewal configuration for a file
type RenewalShardInfo struct {
	ShardId   string `json:"shard_id"`
	ShardSize int    `jons:"shard_size"`
	SPId      string `json:"sp_id"`
}
type RenewalInfo struct {
	CID             string              `json:"cid"`
	ShardsInfo      []*RenewalShardInfo `json:"shards_info"`
	RenewalDuration int                 `json:"renewal_duration"`
	Token           common.Address      `json:"token"`
	Price           int64               `json:"price"`
	Enabled         bool                `json:"enabled"` // if auto renew is enabled
	TotalPay        int64               `json:"total_pay"`
	CreatedAt       time.Time           `json:"created_at"`
	LastRenewalAt   *time.Time          `json:"last_renewal_at,omitempty"`
	NextRenewalAt   time.Time           `json:"next_renewal_at"`
}

func checkAvailableBalance(ctx context.Context, amount int64, token common.Address) error {
	realAmount, err := getRealAmount(amount, token)
	if err != nil {
		return err
	}

	// token: get available balance of token.
	// AvailableBalance, err := chain.SettleObject.VaultService.AvailableBalance(ctx, token)
	AvailableBalance, err := chain.SettleObject.VaultService.AvailableBalance(ctx, token)
	if err != nil {
		return err
	}

	fmt.Printf("check,  balance=%v, realAmount=%v \n", AvailableBalance, realAmount)
	if AvailableBalance.Cmp(realAmount) < 0 {
		fmt.Println("check, err: ", vault.ErrInsufficientFunds)
		return vault.ErrInsufficientFunds
	}
	return nil
}

func getRealAmount(amount int64, token common.Address) (*big.Int, error) {
	// this is price's rate [Compatible with older versions]
	rateObj, err := chain.SettleObject.OracleService.CurrentRate(token)
	if err != nil {
		return nil, err
	}

	realAmount := big.NewInt(0).Mul(big.NewInt(amount), rateObj)
	return realAmount, nil
}

func extendShardEndTime(ctxParams *uh.ContextParams, contractId string, duration int) error {
	k, c, err := sessions.GetUserShardContract(ctxParams.N.Repo.Datastore(), ctxParams.N.Identity.String(), nodepb.ContractStat_RENTER.String(), contractId)
	if err != nil {
		return err
	}

	if c.Meta.ContractId == contractId {
		c.Meta.StorageEnd = uint64(time.Unix(int64(c.Meta.StorageEnd), 0).Add(time.Duration(duration) * time.Hour * 24).Unix())
		err = sessions.UpdateShardContract(ctxParams.N.Repo.Datastore(), c, k)
		if err != nil {
			return err
		}
	}

	return nil
}
