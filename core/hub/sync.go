package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	version "github.com/bittorrent/go-btfs"
	shell "github.com/bittorrent/go-btfs-api"
	config "github.com/bittorrent/go-btfs-config"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/chain/abi"
	chainconfig "github.com/bittorrent/go-btfs/chain/config"
	"github.com/bittorrent/go-btfs/core"

	hubpb "github.com/bittorrent/go-btfs-common/protos/hub"
	"github.com/bittorrent/go-btfs-common/utils/grpc"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	HubModeAll = "all" // special all case (for local reading)

	AllModeHelpText = `
- "score":   top overall score
- "geo":     closest location
- "rep":     highest reputation
- "price":   lowest price
- "speed":   highest transfer speed
- "testnet": testnet-specific
- "all":     all existing hosts`
)

const SP_MODE = "sp"

var SP_MODE_MAP = map[string]int32{
	SP_MODE: 6, // SP_MODE is a special mode for storage providers
}

// CheckValidMode takes in a raw mode string and returns the host request mode type
// if valid, and if local is true and mode is empty, return prefix for storing such
// information into local datastore.
func CheckValidMode(mode string, local bool) (hubpb.HostsReq_Mode, string, error) {
	if mode == SP_MODE {
		return hubpb.HostsReq_Mode(SP_MODE_MAP[SP_MODE]), strings.ToUpper(SP_MODE), nil
	}
	if mode == HubModeAll && local {
		return -1, "", nil
	}
	// Consistent with grpc consts
	modeKey := strings.ToUpper(mode)
	if m, ok := hubpb.HostsReq_Mode_value[modeKey]; ok {
		return hubpb.HostsReq_Mode(m), modeKey, nil
	}
	return -1, "", fmt.Errorf("Invalid Hub query mode: %s", mode)
}

// QueryHosts queries the BTFS-Hub to retrieve the latest list of hosts info
// according to a certain mode.
func QueryHosts(ctx context.Context, node *core.IpfsNode, mode string) ([]*hubpb.Host, error) {
	hrm, _, err := CheckValidMode(mode, false)
	if err != nil {
		return nil, err
	}
	config, err := node.Repo.Config()
	if err != nil {
		return nil, err
	}
	var resp *hubpb.HostsResp
	err = grpc.HubQueryClient(config.Services.HubDomain).WithContext(ctx, func(ctx context.Context,
		client hubpb.HubQueryServiceClient) error {
		resp, err = client.GetHosts(ctx, &hubpb.HostsReq{
			Id:         node.Identity.String(),
			Mode:       hrm,
			Version:    version.CurrentVersionNumber,
			NewVersion: hubpb.HubRouter_V2,
		})
		if err != nil {
			return err
		}
		if resp.Code != hubpb.ResponseCode_SUCCESS {
			return fmt.Errorf(resp.Message)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to query hosts from Hub service: %v", err)
	}

	return resp.Hosts.Hosts, nil
}

// GetSP get sp node from proposal contract
func GetSP(ctx context.Context, cfg *config.Config) ([]*hubpb.Host, error) {
	cli := chain.ChainObject.Backend
	if cli == nil {
		client, err := ethclient.Dial(cfg.ChainInfo.Endpoint)
		if err != nil {
			fmt.Printf("failed to dial eth client: %v", err)
			return nil, err
		}
		cli = client
	}
	cc, ok := chainconfig.GetChainConfig(cfg.ChainInfo.ChainId)
	if !ok {
		return nil, fmt.Errorf("chain %d is not supported yet", cfg.ChainInfo.ChainId)
	}
	proposal, err := abi.NewProposalContract(cc.ProposalAddress, cli)
	if err != nil {
		return nil, err
	}
	proposals, err := proposal.GetProposalsByType(nil, 0)
	if err != nil {
		return nil, err
	}
	uris := make([]string, 0)
	for _, p := range proposals {
		if p.Status != 1 {
			continue
		}
		uris = append(uris, p.Uri)
	}
	// request uri to get json file
	hosts := make([]*hubpb.Host, 0)
	for _, uri := range uris {
		// parse json file
		us := strings.Split(uri, "/")
		if len(us) < 5 {
			continue
		}
		cid := us[len(us)-1]
		cat, err := shell.NewLocalShell().Cat(cid)
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(cat)
		if err != nil {
			fmt.Println("read body failed:", err)
			return nil, err
		}

		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println("parse json failed:", err)
			return nil, err
		}

		hosts = append(hosts, &hubpb.Host{
			NodeId: data["storageProviderNodeId"].(string),
		})
	}
	hosts = append(hosts, &hubpb.Host{
		NodeId: "16Uiu2HAmKsMZH6RNi6rWCpfcnzek2iPjcNcvrcgjL7pERFCpr1AV",
	})

	return hosts, nil
}

// QueryStats queries the BTFS-Hub to retrieve the latest storage stats on this host.
func QueryStats(ctx context.Context, node *core.IpfsNode, v2 bool) (*hubpb.StatsResp, error) {
	config, err := node.Repo.Config()
	if err != nil {
		return nil, err
	}

	newVersion := hubpb.HubRouter_V2
	if !v2 {
		newVersion = hubpb.HubRouter_V1
	}

	var resp *hubpb.StatsResp
	err = grpc.HubQueryClient(config.Services.HubDomain).WithContext(ctx, func(ctx context.Context,
		client hubpb.HubQueryServiceClient) error {
		resp, err = client.GetStats(ctx, &hubpb.StatsReq{
			Id:         node.Identity.String(),
			NewVersion: newVersion,
		})
		if err != nil {
			return err
		}
		if resp.Code != hubpb.ResponseCode_SUCCESS {
			return fmt.Errorf(resp.Message)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to query stats from Hub service: %v", err)
	}

	return resp, nil
}
