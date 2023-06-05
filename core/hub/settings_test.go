package hub

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	nodepb "github.com/bittorrent/go-btfs-common/protos/node"
)

func TestGetSettings(t *testing.T) {
	ns, err := GetHostSettings(context.Background(), "https://score.btfs.io",
		"QmWJWGxKKaqZUW4xga2BCzT5FBtYDL8Cc5Q5jywd6xPt1g")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("GetHostSettings, ns = %+v \n", ns)

	defNs := &nodepb.Node_Settings{StoragePriceAsk: 125000, StorageTimeMin: 30, StoragePriceDefault: 125000}
	if !reflect.DeepEqual(ns, defNs) {
		t.Fatal("default settings not equal")
	}
}
