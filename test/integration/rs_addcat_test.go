package integrationtest

import (
	"github.com/bittorrent/go-btfs/thirdparty/unit"

	files "github.com/TRON-US/go-btfs-files"

	"testing"

	"github.com/TRON-US/interface-go-btfs-core/options"
	"github.com/stretchr/testify/assert"

	testutil "github.com/libp2p/go-libp2p-testing/net"
)

var conf = testutil.LatencyConfig{
	NetworkLatency:    0,
	RoutingLatency:    0,
	BlockstoreLatency: 0,
}

func TestAddCatFile(t *testing.T) {
	data := RandomBytes(1 * unit.KB)
	buf, err := addCat("", files.NewBytesFile(data), conf, options.Unixfs.Chunker("reed-solomon"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, data, buf.Bytes(), "catted data does not match added data")
}

func TestAddCatNestedFile(t *testing.T) {
	data := RandomBytes(1 * unit.KB)
	f1 := files.NewMapDirectory(map[string]files.Node{
		"dir1": files.NewMapDirectory(map[string]files.Node{
			"dir11": files.NewMapDirectory(map[string]files.Node{
				"dir111": files.NewMapDirectory(map[string]files.Node{
					"file1.txt": files.NewBytesFile(RandomBytes(2 * unit.KB)),
					"file2.txt": files.NewBytesFile(data),
					"file3.txt": files.NewBytesFile(RandomBytes(2 * unit.KB)),
				}),
			}),
		}),
	})
	buf, err := addCat("dir1/dir11/dir111/file2.txt", f1, conf, options.Unixfs.Chunker("reed-solomon"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, data, buf.Bytes(), "catted data does not match added data")
}
