package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/bittorrent/go-btfs/gc"
	ft "github.com/bittorrent/go-unixfs"
	importer "github.com/bittorrent/go-unixfs/importer"
	uio "github.com/bittorrent/go-unixfs/io"
	"github.com/bittorrent/go-unixfs/mod"
	ftutil "github.com/bittorrent/go-unixfs/util"
	bserv "github.com/ipfs/go-blockservice"
	merkledag "github.com/ipfs/go-merkledag"

	chunker "github.com/bittorrent/go-btfs-chunker"
	config "github.com/bittorrent/go-btfs-config"
	files "github.com/bittorrent/go-btfs-files"
	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/coreapi"
	"github.com/bittorrent/go-btfs/core/coreunix"
	"github.com/bittorrent/go-btfs/repo"
	coreiface "github.com/bittorrent/interface-go-btfs-core"
	ipath "github.com/bittorrent/interface-go-btfs-core/path"
	cid "github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	syncds "github.com/ipfs/go-datastore/sync"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	u "github.com/ipfs/go-ipfs-util"
	ipld "github.com/ipfs/go-ipld-format"
)

type MetaStruct struct {
	Price  float64
	Number int
	NodeId string
}

func getDagserv(t *testing.T) ipld.DAGService {
	db := dssync.MutexWrap(ds.NewMapDatastore())
	bs := bstore.NewBlockstore(db)
	blockserv := bserv.New(bs, offline.Exchange(bs))
	return merkledag.NewDAGService(blockserv)
}

func TestMetadata(t *testing.T) {
	ctx := context.Background()
	// Make some random node
	ds := getDagserv(t)
	data := make([]byte, 1000)
	_, err := io.ReadFull(u.NewTimeSeededRand(), data)
	if err != nil {
		t.Fatal(err)
	}
	r := bytes.NewReader(data)
	nd, err := importer.BuildDagFromReader(ds, chunker.DefaultSplitter(r))
	if err != nil {
		t.Fatal(err)
	}

	c := nd.Cid()

	m := new(ft.Metadata)
	m.MimeType = "THIS IS A TEST"

	// Such effort, many compromise
	ipfsnode := &core.IpfsNode{DAG: ds}

	mdk, err := coreunix.AddMetadataTo(ipfsnode, c.String(), m)
	if err != nil {
		t.Fatal(err)
	}

	rec, err := coreunix.Metadata(ipfsnode, mdk)
	if err != nil {
		t.Fatal(err)
	}
	if rec.MimeType != m.MimeType {
		t.Fatalf("something went wrong in conversion: '%s' != '%s'", rec.MimeType, m.MimeType)
	}

	cdk, err := cid.Decode(mdk)
	if err != nil {
		t.Fatal(err)
	}

	retnode, err := ds.Get(ctx, cdk)
	if err != nil {
		t.Fatal(err)
	}

	rtnpb, ok := retnode.(*merkledag.ProtoNode)
	if !ok {
		t.Fatal("expected protobuf node")
	}

	ndr, err := uio.NewDagReader(ctx, rtnpb, ds)
	if err != nil {
		t.Fatal(err)
	}

	out, err := ioutil.ReadAll(ndr)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(out, data) {
		t.Fatal("read incorrect data")
	}
}

func addFileWithMetadata(node *core.IpfsNode, data []byte, metadata string) (ipath.Path, error) {
	return addFileToBtfs(node, data, metadata)
}

func addFile(node *core.IpfsNode, data []byte) (ipath.Path, error) {
	return addFileToBtfs(node, data, "")
}

func addFileWithoutMetadata(node *core.IpfsNode, data []byte) (ipath.Path, error) {
	return addFileToBtfs(node, data, "")
}

func addFileToBtfs(node *core.IpfsNode, data []byte, metadata string) (ipath.Path, error) {
	// Create out chan and adder
	output := make(chan interface{})
	adder, err := coreunix.NewAdder(context.Background(), node.Pinning, node.Blockstore, node.DAG)
	if err != nil {
		return nil, err
	}
	adder.Out = output
	if metadata != "" {
		adder.TokenMetadata = metadata
	}

	// Create a file with goroutine
	file := files.NewBytesFile(data)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	addDone := make(chan struct{})
	go func() {
		defer close(output)
		defer close(addDone)
		_, err := adder.AddAllAndPin(ctx, file)
		if err != nil {
			output <- err
		}
	}()

	var addedFileHash cid.Cid
	select {
	case o := <-output:
		err, ok := o.(error)
		if ok {
			return nil, err
		}
		addedFileHash = o.(*coreiface.AddEvent).Path.Cid()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	<-addDone

	return ipath.IpfsPath(addedFileHash), nil
}

func addDirectoryToBtfs(node *core.IpfsNode, file files.Node, metadata string, rs bool) (ipath.Path, error) {
	// Make sure given `file` is a directory
	dir, ok := file.(files.Directory)
	if !ok {
		return nil, errors.New("expected directory node")
	}

	// Create out chan and adder
	output := make(chan interface{})
	adder, err := coreunix.NewAdder(context.Background(), node.Pinning, node.Blockstore, node.DAG)
	if err != nil {
		return nil, err
	}
	adder.Out = output
	if metadata != "" {
		adder.TokenMetadata = metadata
		//
		if _, ok := file.(files.Directory); ok {
			adder.MetaForDirectory = true
		}
	}
	if rs {
		dsize, psize, csize := TestRsDataSize, TestRsParitySize, 262144
		adder.Chunker = fmt.Sprintf("reed-solomon-%d-%d-%d", dsize, psize, csize)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var rsfileAdder *coreunix.ReedSolomonAdder
	if rs {
		rsfileAdder, err = coreunix.NewReedSolomonAdder(adder)
		if err != nil {
			return nil, err
		}
		rsfileAdder.IsDir = true
		if files.IsMultiPartDirectory(dir) {
			rsfileAdder.FileType = coreunix.MultipartFile
		} else if files.IsMapDirectory(dir) {
			rsfileAdder.FileType = coreunix.SliceFile
		} else if files.IsSerialFileDirectory(dir) {
			rsfileAdder.FileType = coreunix.SerialFile
		} else {
			return nil, fmt.Errorf("unexpected files.Directory type [%T]", dir)
		}
	}

	addDone := make(chan struct{})
	go func() {
		defer close(output)
		defer close(addDone)
		if !rs {
			_, err = adder.AddAllAndPin(ctx, file)

		} else {
			_, err = rsfileAdder.AddAllAndPin(ctx, file)
		}
		if err != nil {
			output <- err
		}
	}()

	var addedFileHash cid.Cid
	var size int
	if rs {
		size = 1
	} else {
		size, err = directorySize(dir)
		if err != nil {
			return nil, err
		}
	}

	for i := 0; i < int(size); i++ {
		select {
		case o := <-output:
			err, ok := o.(error)
			if ok {
				return nil, err
			}
			addedFileHash = o.(*coreiface.AddEvent).Path.Cid()
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	<-addDone

	return ipath.IpfsPath(addedFileHash), nil
}

func verifyUnpinViaGC(node *core.IpfsNode, removed ipld.Node) error {
	var gcOutput <-chan gc.Result
	gcStarted := make(chan struct{})
	go func() {
		defer close(gcStarted)
		gcOutput = gc.GC(context.Background(), node.Blockstore, node.Repo.Datastore(), node.Pinning, nil)
	}()

	<-gcStarted

	removedHashes := make(map[string]struct{})
	for item := range gcOutput {
		if item.Error != nil {
			return item.Error
		}
		removedHashes[item.KeyRemoved.String()] = struct{}{}
	}
	if _, found := removedHashes[removed.Cid().String()]; !found {
		return errors.New("Previous node not gc'ed")
	}
	return nil
}

func addMetadata(node *core.IpfsNode, path ipath.Path, meta string, rootWillBeUnpinned bool) (ipath.Path, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	nd, err := api.ResolveNode(ctx, path)
	if err != nil {
		return nil, err
	}

	output := make(chan interface{})
	modifier, err := coreunix.NewMetaModifier(context.Background(), node.Pinning, node.Blockstore, node.DAG)
	if err != nil {
		return nil, err
	}
	modifier.Out = output
	modifier.Overwrite = true
	metaBytes := []byte(meta)
	modifier.TokenMetadata = string(metaBytes)

	var n ipld.Node
	go func() {
		defer close(output)
		n, _, err = modifier.AddMetaAndPin(nd)
		if err != nil {
			output <- err
		} else {
			modifier.Out <- n.Cid()
		}
	}()

	var modifiedFileHash cid.Cid
	select {
	case o := <-output:
		err, ok := o.(error)
		if ok {
			return nil, err
		}
		modifiedFileHash = o.(cid.Cid)
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if rootWillBeUnpinned {
		err = verifyUnpinViaGC(node, nd)
		if err != nil {
			return nil, err
		}
	}

	return ipath.IpfsPath(modifiedFileHash), nil
}

// TestAddFileMetadata tests the functionality
// to add a file with meta data items to BTFS network.
func TestAddFileWithMetadata(t *testing.T) {
	// Create repo.
	r := &repo.Mock{
		C: config.Config{
			Identity: config.Identity{
				PeerID: testPeerID, // required by offline node
			},
		},
		D: syncds.MutexWrap(datastore.NewMapDatastore()),
	}
	// Build Ipfs node.
	node, err := core.NewNode(context.Background(), &core.BuildCfg{Repo: r})
	if err != nil {
		t.Fatal(err)
	}

	// Add a file to BTFS network.
	expected := "{\"price\":11.2,\"number\":2368}"
	path, err := addFileWithMetadata(node, []byte("existing file contents for adding file with metadata"), expected)
	if err != nil {
		t.Fatal(err)
	}

	// Verify modified file.
	verifyMetadataItems(t, node, path, &MetaStruct{Price: 11.2, Number: 2368})
	if path.String() != "/btfs/QmRM4vKujpTaafJs6DfF2jpxadvgktYAdbK7JZYhswr9V6" {
		// note: the exact number will depend on the size and the sharding algo. used
		t.Fatalf("expected %s, got %s", "/btfs/QmRM4vKujpTaafJs6DfF2jpxadvgktYAdbK7JZYhswr9V6", path.String())
	}
}

func TestAppendMetadata(t *testing.T) {
	// Create repo.
	r := &repo.Mock{
		C: config.Config{
			Identity: config.Identity{
				PeerID: testPeerID, // required by offline node
			},
		},
		D: syncds.MutexWrap(datastore.NewMapDatastore()),
	}
	// Build Ipfs node.
	node, err := core.NewNode(context.Background(), &core.BuildCfg{Repo: r})
	if err != nil {
		t.Fatal(err)
	}

	// Add a file to BTFS network.
	data := []byte("existing file contents for appending metadata")
	path, err := addFileWithMetadata(node, data, `{"price":11.21}`)
	if err != nil {
		t.Fatal(err)
	}

	// Append token metadata to the BTFS file.
	p, err := addMetadata(node, path, fmt.Sprintf(`{"number":%d}`, 1234), false)
	if err != nil {
		t.Fatal(err)
	}

	// Verify modified file.
	verifyMetadataItems(t, node, p, &MetaStruct{Price: 11.21, Number: 1234})

	if p.String() != "/btfs/QmYu6CoRhJATcu5QM9XhsT85jE5sBrV61SByALHG6kz4kH" {
		// note: the exact number will depend on the size and the sharding algo. used
		t.Fatalf("expected %s, got %s", "/btfs/QmYu6CoRhJATcu5QM9XhsT85jE5sBrV61SByALHG6kz4kH", p.String())
	}
}

// TestAddMetadataToFileWithoutMeta tests the functionality
// to add metadata items to an existing BTFS file without metadata.
func TestAddMetadataToFileWithoutMeta(t *testing.T) {
	// Create repo.
	r := &repo.Mock{
		C: config.Config{
			Identity: config.Identity{
				PeerID: testPeerID, // required by offline node
			},
		},
		D: syncds.MutexWrap(datastore.NewMapDatastore()),
	}
	// Build Ipfs node.
	node, err := core.NewNode(context.Background(), &core.BuildCfg{Repo: r})
	if err != nil {
		t.Fatal(err)
	}

	// Add a file to BTFS network.
	path, err := addFile(node, []byte("existing file contents for adding without meta"))
	if err != nil {
		t.Fatal(err)
	}

	// Add token metadata to the BTFS file without existing metadata.
	expected := `{"number":2368}#{}`
	p, err := addMetadata(node, path, `{"number":2368}`, false)
	if err != nil {
		t.Fatal(err)
	}

	// Verify modified file.
	metaSize, err := metaDataSize(node, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(expected) != metaSize {
		t.Fatalf("expected %d, got %d", len(expected), int(metaSize))
	}
	verifyMetadataItems(t, node, p, &MetaStruct{Number: 2368})
	if p.String() != "/btfs/QmPob5g1dZfqFeL56SRqpVwXr7YVwkW86iMgsD7xc5iKwb" {
		// note: the exact number will depend on the size and the sharding algo. used
		t.Fatalf("expected %s, got %s", "/btfs/QmPob5g1dZfqFeL56SRqpVwXr7YVwkW86iMgsD7xc5iKwb", p.String())
	}
}

func TestUpdateMetadata(t *testing.T) {
	// Create repo.
	r := &repo.Mock{
		C: config.Config{
			Identity: config.Identity{
				PeerID: testPeerID, // required by offline node
			},
		},
		D: syncds.MutexWrap(datastore.NewMapDatastore()),
	}
	// Build Ipfs node.
	node, err := core.NewNode(context.Background(), &core.BuildCfg{Repo: r})
	if err != nil {
		t.Fatal(err)
	}

	// Add a file to BTFS network.
	data := []byte("existing file contents for updating metadata")
	path, err := addFileWithMetadata(node, data, `{"price":11.23}`)
	if err != nil {
		t.Fatal(err)
	}

	// Append token metadata to the BTFS file.
	p, err := addMetadata(node, path, `{"price":23.56,"number":4356,"NewItem":"justSave"}`, false)
	if err != nil {
		t.Fatal(err)
	}

	// Verify modified file.
	verifyMetadataItems(t, node, p, &MetaStruct{Price: 23.56, Number: 4356})

	if p.String() != "/btfs/QmUZVy85njUTjZygd1KeHbWfYumKjsE5ewFD8W7RYNVTwQ" {
		// note: the exact number will depend on the size and the sharding algo. used
		t.Fatalf("expected %s, got %s", "/btfs/QmUZVy85njUTjZygd1KeHbWfYumKjsE5ewFD8W7RYNVTwQ", p.String())
	}
}

// removeMetadata removes the metadata items for the key values from the given `meta` string.
// `meta` will have a format of key list. E.g., "price,nodeId".
func removeMetadata(node *core.IpfsNode, path ipath.Path, meta string, rootWillBeUnpinned bool) (ipath.Path, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	nd, err := api.ResolveNode(ctx, path)
	if err != nil {
		return nil, err
	}

	output := make(chan interface{})
	modifier, err := coreunix.NewMetaModifier(context.Background(), node.Pinning, node.Blockstore, node.DAG)
	if err != nil {
		return nil, err
	}
	modifier.Out = output
	modifier.TokenMetadata = meta

	var n ipld.Node
	go func() {
		defer close(output)
		n, _, err = modifier.RemoveMetaAndPin(nd)
		if err != nil {
			output <- err
		} else {
			modifier.Out <- n.Cid()
		}
	}()

	var modifiedFileHash cid.Cid
	select {
	case o := <-output:
		err, ok := o.(error)
		if ok {
			return nil, err
		}
		modifiedFileHash = o.(cid.Cid)
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if rootWillBeUnpinned {
		err = verifyUnpinViaGC(node, nd)
		if err != nil {
			return nil, err
		}
	}
	return ipath.IpfsPath(modifiedFileHash), nil
}

func TestRemoveMetadata(t *testing.T) {
	// Create repo.
	r := &repo.Mock{
		C: config.Config{
			Identity: config.Identity{
				PeerID: testPeerID, // required by offline node
			},
		},
		D: syncds.MutexWrap(datastore.NewMapDatastore()),
	}
	// Build Ipfs node.
	node, err := core.NewNode(context.Background(), &core.BuildCfg{Repo: r})
	if err != nil {
		t.Fatal(err)
	}

	// Add a file to BTFS network.
	data := []byte("existing file contentsfor removing")
	path, err := addFileWithMetadata(node, data, `{"price":11.25,"number":1234}`)
	if err != nil {
		t.Fatal(err)
	}

	// Append token metadata to the BTFS file.
	p, err := removeMetadata(node, path, `price,id`, false)
	if err != nil {
		t.Fatal(err)
	}

	// Verify modified file.
	verifyMetadataItems(t, node, p, &MetaStruct{Number: 1234})

	if p.String() != "/btfs/QmcSVqDuQaNDPrgMUGhmcWWR6qyavjKYtsnqn6xFVicufz" {
		// note: the exact number will depend on the size and the sharding algo. used
		t.Fatalf("expected %s, got %s", "/btfs/QmcSVqDuQaNDPrgMUGhmcWWR6qyavjKYtsnqn6xFVicufz", p.String())
	}
}

// metaDataSize returns the size of the metadata sub-DAG of
// the DAG of the given `p`.
func metaDataSize(node *core.IpfsNode, p ipath.Path) (int, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return 0, err
	}
	nd, err := api.ResolveNode(ctx, p)
	if err != nil {
		return 0, err
	}
	mnode, err := ft.GetMetaSubdagRoot(ctx, nd, node.DAG)
	if err != nil {
		return 0, err
	}
	fileSize, err := mod.FileSize(mnode)
	if err != nil {
		return 0, err
	}
	return int(fileSize), nil
}

func verifyMetadataItems(t *testing.T, node *core.IpfsNode, p ipath.Path, exp *MetaStruct) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		t.Fatal(err)
	}

	b, err := coreunix.GetMetaData(ctx, api, p)
	if err != nil {
		t.Fatal(err)
	}
	var meta MetaStruct
	b1 := ftutil.GetMetadataElement(b)
	if err := json.Unmarshal(b1, &meta); err != nil {
		t.Fatal(err)
	}

	if meta.Price != exp.Price {
		t.Fatalf("expected %.2f, got %.2f", exp.Price, meta.Price)
	}
	if meta.Number != exp.Number {
		t.Fatalf("expected %d, got %d", exp.Number, meta.Number)
	}
	if meta.NodeId != exp.NodeId {
		t.Fatalf("expected %s, got %s", exp.NodeId, meta.NodeId)
	}
}

// TestAddAddOneLevelDirecotryWithMetadata tests the functionality
// to add a one level directory with meta data items to BTFS network.
func TestAddOneLevelDirectoryWithMetadata(t *testing.T) {
	testAddDirectoryWithMetadata(t, oneLevelDirectory(), "/btfs/QmTsvHDpjEGEEUKZsUGWgzDmszmkwgteMSZvHRjCCnRJps", false)
}

func TestAddOneLevelDirectoryWithMetadataReedSolomon(t *testing.T) {
	testAddDirectoryWithMetadata(t, oneLevelDirectory(), "/btfs/QmeFqYjDyioqi7SRDZmJJivPd6s9q6Cu7Mgt1u4RTdKRpJ", true)
}

// TestAddAddOneLevelDirecotryWithMetadata tests the functionality
// to add a one level directory with meta data items to BTFS network.
func TestAddTwoLevelDirectoryWithMetadata(t *testing.T) {
	testAddDirectoryWithMetadata(t, twoLevelDirectory(), "/btfs/QmVmGgo864dMRAFkW8iUMN42eaP4yPbuF9PjQdb25aQwMQ", false)
}

func TestAddTwoLevelDirectoryWithMetadataReedSolomon(t *testing.T) {
	testAddDirectoryWithMetadata(t, twoLevelDirectory(), "/btfs/QmeMVoYW8z6cpCByMzn5E6YMK2Ah9fwkhy4huanep5c5bz", true)
}

func testAddDirectoryWithMetadata(t *testing.T, file files.Node, originalMetadata string, reedSolomon bool) {
	testAddDirectory(t, file, originalMetadata, reedSolomon)
}

func testAddDirectory(t *testing.T, file files.Node, originalMetadataString string, reedSolomon bool) {
	// Create repo.
	r := &repo.Mock{
		C: config.Config{
			Identity: config.Identity{
				PeerID: testPeerID, // required by offline node
			},
		},
		D: syncds.MutexWrap(datastore.NewMapDatastore()),
	}
	// Build Ipfs node.
	node, err := core.NewNode(context.Background(), &core.BuildCfg{Repo: r})
	if err != nil {
		t.Fatal(err)
	}

	// Add a file to BTFS network.
	expected := "{\"price\":11.2,\"number\":2368}"
	path, err := addDirectoryToBtfs(node, file, expected, reedSolomon)
	if err != nil {
		t.Fatal(err)
	}

	// Verify modified file.
	verifyMetadataItems(t, node, path, &MetaStruct{Price: 11.2, Number: 2368})
	if path.String() != originalMetadataString {
		// note: the exact number will depend on the size and the sharding algo. used
		t.Fatalf("expected %s, got %s", originalMetadataString, path.String())
	}
}

func oneLevelDirectory() files.Node {
	return files.NewMapDirectory(map[string]files.Node{
		"file1": files.NewBytesFile([]byte("contents for file one")),
		"file2": files.NewBytesFile([]byte("hello, world")),
	})
}

func twoLevelDirectory() files.Node {
	return files.NewMapDirectory(map[string]files.Node{
		"file1": files.NewBytesFile([]byte("contents for file one")),
		"file2": files.NewBytesFile([]byte("hello, world")),
		"dir1": files.NewMapDirectory(map[string]files.Node{
			"file3": files.NewBytesFile([]byte("test data")),
		}),
	})
}

func directorySize(f files.Node) (int, error) {
	// Error cases
	if f == nil {
		return 0, errors.New("nil argument encountered")
	}
	curr, ok := f.(files.Directory)
	if !ok {
		return 0, errors.New("unexpected node type")
	}

	// Normal/recursive case
	size := 1
	it := curr.Entries()

	for it.Next() {
		_, ok := it.Node().(files.Directory)
		if ok {
			subSize, err := directorySize(it.Node())
			if err != nil {
				return 0, err
			}
			size += subSize
		} else {
			size++
		}

	}
	return size, nil
}
