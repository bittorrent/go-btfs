package mock_test

import (
	"testing"

	"github.com/bittorrent/go-btfs/statestore/mock"
	"github.com/bittorrent/go-btfs/statestore/test"
	"github.com/bittorrent/go-btfs/transaction/storage"
)

func TestMockStateStore(t *testing.T) {
	test.Run(t, func(t *testing.T) storage.StateStorer {
		return mock.NewStateStore()
	})
}
