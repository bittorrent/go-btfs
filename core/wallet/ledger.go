package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	walletpb "github.com/bittorrent/go-btfs/protos/wallet"

	ledgerpb "github.com/tron-us/go-btfs-common/protos/ledger"

	"github.com/golang/protobuf/proto"
	"github.com/ipfs/go-datastore"
)

const (
	channelKeyPrefix   = "/ledger-channels"
	channelKeyTemplate = channelKeyPrefix + "/%d"
)

func save(ds datastore.Datastore, state *ledgerpb.ChannelState) error {
	if state == nil || state.Id == nil {
		return errors.New("state or state.Id is nil")
	}
	return sessions.Save(ds, k(state.Id.Id), &walletpb.ChannelState{State: state, TimeCreate: time.Now()})
}

func list(ds datastore.Datastore) ([]*walletpb.ChannelState, error) {
	list, err := sessions.List(ds, channelKeyPrefix)
	if err != nil {
		return nil, err
	}
	var states []*walletpb.ChannelState
	for _, e := range list {
		state := &walletpb.ChannelState{}
		if err := proto.Unmarshal(e, state); err != nil {
			log.Debug(err)
			continue
		}
		states = append(states, state)
	}
	return states, nil
}

func rm(ctx context.Context, ds datastore.Datastore, channelId int64) error {
	return ds.Delete(ctx, datastore.NewKey(k(channelId)))
}

func k(channelId int64) string {
	return fmt.Sprintf(channelKeyTemplate, channelId)
}
