package hub

import (
	"context"
	"fmt"
	"strings"

	version "github.com/bittorrent/go-btfs"

	"github.com/bittorrent/go-btfs/core"

	hubpb "github.com/bittorrent/go-btfs-common/protos/hub"
	"github.com/bittorrent/go-btfs-common/utils/grpc"
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

// CheckValidMode takes in a raw mode string and returns the host request mode type
// if valid, and if local is true and mode is empty, return prefix for storing such
// information into local datastore.
func CheckValidMode(mode string, local bool) (hubpb.HostsReq_Mode, string, error) {
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
			Id:         node.Identity.Pretty(),
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
			Id:         node.Identity.Pretty(),
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
