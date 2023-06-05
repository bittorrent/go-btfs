package repo

import (
	"testing"

	"github.com/bittorrent/go-btfs-common/protos/hub"

	"github.com/ipfs/go-datastore"
	syncds "github.com/ipfs/go-datastore/sync"
	"github.com/stretchr/testify/assert"
)

func TestGetPut(t *testing.T) {
	d := syncds.MutexWrap(datastore.NewMapDatastore())
	settingData := &hub.SettingsData{
		StoragePriceAsk:   1.1,
		BandwidthPriceAsk: 1.2,
		StorageTimeMin:    1.3,
		BandwidthLimit:    1.4,
		CollateralStake:   1.5,
	}
	k := "setting"
	err := Put(d, k, settingData)
	if err != nil {
		t.Fatal(err)
	}
	m := new(hub.SettingsData)
	sd, err := Get(d, k, m)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, sd, settingData)
}
