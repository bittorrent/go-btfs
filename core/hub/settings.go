package hub

import (
	"context"
	"errors"

	hubpb "github.com/bittorrent/go-btfs-common/protos/hub"
	nodepb "github.com/bittorrent/go-btfs-common/protos/node"
	"github.com/bittorrent/go-btfs-common/utils/grpc"
)

func GetHostSettings(ctx context.Context, addr, peerId string) (*nodepb.Node_Settings, error) {
	// get from remote
	ns := new(nodepb.Node_Settings)
	err := grpc.HubQueryClient(addr).WithContext(ctx, func(ctx context.Context, client hubpb.HubQueryServiceClient) error {
		req := new(hubpb.SettingsReq)
		req.Id = peerId
		req.NewVersion = hubpb.HubRouter_V2
		resp, err := client.GetSettings(ctx, req)
		if err != nil {
			return err
		}
		if resp.Code != hubpb.ResponseCode_SUCCESS {
			return errors.New(resp.Message)
		}
		ns.StoragePriceAsk = uint64(resp.SettingsData.StoragePriceAsk)
		ns.StoragePriceDefault = ns.StoragePriceAsk
		ns.CustomizedPricing = false
		// XXX: These configs need to be controlled by a customized flag as well
		ns.StorageTimeMin = uint64(resp.SettingsData.StorageTimeMin)
		ns.BandwidthLimit = resp.SettingsData.BandwidthLimit
		ns.BandwidthPriceAsk = uint64(resp.SettingsData.BandwidthPriceAsk)
		ns.CollateralStake = uint64(resp.SettingsData.CollateralStake)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ns, nil
}
