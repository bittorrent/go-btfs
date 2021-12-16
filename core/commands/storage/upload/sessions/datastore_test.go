package sessions

import (
	"testing"
	"time"

	coremock "github.com/bittorrent/go-btfs/core/mock"
	sessionpb "github.com/bittorrent/go-btfs/protos/session"

	"github.com/stretchr/testify/assert"
)

func TestSaveGetRemove(t *testing.T) {
	node, err := coremock.NewMockNode()
	if err != nil {
		t.Fatal(err)
	}

	shardHashes := make([]string, 0)
	shardHashes = append(shardHashes, "Qm1")
	shardHashes = append(shardHashes, "Qm2")
	shardHashes = append(shardHashes, "Qm3")
	current := time.Now().UTC()
	md := &sessionpb.Metadata{
		TimeCreate:  current,
		RenterId:    node.Identity.String(),
		FileHash:    "Qm123",
		ShardHashes: shardHashes,
	}
	err = Save(node.Repo.Datastore(), "ds.ds.test", md)
	if err != nil {
		t.Fatal(err)
	}
	newMd := &sessionpb.Metadata{}
	err = Get(node.Repo.Datastore(), "ds.ds.test", newMd)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, md, newMd)
	err = Remove(node.Repo.Datastore(), "ds.ds.test")
	if err != nil {
		t.Fatal(err)
	}
	err = Get(node.Repo.Datastore(), "ds.ds.test", newMd)
	if err == nil {
		t.Fatal("ds.ds.test should have been removed")
	}
}
