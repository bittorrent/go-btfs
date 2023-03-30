package coreunix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gopath "path"

	core "github.com/bittorrent/go-btfs/core"

	chunker "github.com/TRON-US/go-btfs-chunker"
	"github.com/TRON-US/go-mfs"
	ft "github.com/TRON-US/go-unixfs"
	ihelper "github.com/TRON-US/go-unixfs/importer/helpers"
	uio "github.com/TRON-US/go-unixfs/io"
	"github.com/TRON-US/go-unixfs/mod"
	ftutil "github.com/TRON-US/go-unixfs/util"
	coreiface "github.com/TRON-US/interface-go-btfs-core"
	ipath "github.com/TRON-US/interface-go-btfs-core/path"
	cid "github.com/ipfs/go-cid"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	pin "github.com/ipfs/go-ipfs-pinner"
	posinfo "github.com/ipfs/go-ipfs-posinfo"
	ipld "github.com/ipfs/go-ipld-format"
	dag "github.com/ipfs/go-merkledag"
)

// MetaModifier contains the options to the `metadata` command.
type MetaModifier struct {
	ctx context.Context
	*Adder
	Overwrite bool
}

func NewMetaModifier(ctx context.Context, p pin.Pinner, bs bstore.GCLocker, ds ipld.DAGService) (*MetaModifier, error) {
	adder, err := NewAdder(ctx, p, bs, ds)
	if err != nil {
		return nil, err
	}

	return &MetaModifier{ctx, adder, false}, nil
}

// Add metadata to the given DAG root node for a BTFS file.
func (modifier *MetaModifier) AddMetaAndPin(node ipld.Node) (ipld.Node, ipath.Resolved, error) {
	if modifier.Pin {
		modifier.unlocker = modifier.gcLocker.PinLock(modifier.ctx)
	}
	defer func() {
		if modifier.unlocker != nil {
			modifier.unlocker.Unlock(modifier.ctx)
		}
	}()

	p, err := modifier.addMetaToFileNode(node)
	if err != nil {
		return nil, nil, err
	}

	root, err := modifier.updateMfs()
	if err != nil {
		return nil, nil, err
	}

	nd, err := root.GetNode()
	if err != nil {
		return nil, nil, err
	}

	if !modifier.Pin {
		return nd, p, nil
	}

	if node.Cid() != nd.Cid() {
		err = modifier.pinning.Unpin(modifier.ctx, node.Cid(), true)
		if err != nil {
			return nil, nil, err
		}
	}

	return nd, p, modifier.PinRoot(modifier.ctx, nd)
}

func (modifier *MetaModifier) updateMfs() (mfs.FSNode, error) {
	// get root
	mr, err := modifier.mfsRoot()
	if err != nil {
		return nil, err
	}
	var root mfs.FSNode
	rootdir := mr.GetDirectory()
	root = rootdir

	err = root.Flush()
	if err != nil {
		return nil, err
	}

	// We are adding a file without wrapping, so swap the root to it
	var name string
	children, err := rootdir.ListNames(modifier.ctx)
	if err != nil {
		return nil, err
	}

	if len(children) == 0 {
		return nil, fmt.Errorf("expected at least one child dir, got none")
	}

	// Replace root with the first child
	name = children[0]
	root, err = rootdir.Child(name)
	if err != nil {
		return nil, err
	}

	// Now close mr
	err = mr.Close()
	if err != nil {
		return nil, err
	}
	return root, nil
}

func (modifier *MetaModifier) addMetaToFileNode(node ipld.Node) (ipath.Resolved, error) {
	err := modifier.maybePauseForGC(modifier.ctx)
	if err != nil {
		return nil, err
	}

	dagnode, err := modifier.addMeta(node)
	if err != nil {
		return nil, err
	}

	// patch it into the root
	return modifier.addNodeToMfs(dagnode, "")
}

func (adder *Adder) addNodeToMfs(node ipld.Node, path string) (ipath.Resolved, error) {
	if path == "" {
		path = node.Cid().String()
	}

	if pi, ok := node.(*posinfo.FilestoreNode); ok {
		node = pi.Node
	}

	mr, err := adder.mfsRoot()
	if err != nil {
		return nil, err
	}
	// From core/coreunix/add.go
	// TODO: check with the same execution path in add.go for `btfs add` command.
	dir := gopath.Dir(path)
	if dir != "." {
		opts := mfs.MkdirOpts{
			Mkparents:  true,
			Flush:      false,
			CidBuilder: adder.CidBuilder,
		}
		if err := mfs.Mkdir(mr, dir, opts); err != nil {
			return nil, err
		}
	}

	if err := mfs.PutNode(mr, path, node); err != nil {
		return nil, err
	}

	c := node.Cid()

	return ipath.IpfsPath(c), nil
}

// Add the given `modifier.TokenMetadata` to the given `node`.
func (modifier *MetaModifier) addMeta(nd ipld.Node) (ipld.Node, error) {
	// Error cases
	if nd == nil {
		return nil, errors.New("Invalid argument: nil node value")
	}

	// Build Metadata DAG modifier for balanced format
	// TODO: trickle format
	var metadataNotExist bool
	var mnode ipld.Node
	dr := []byte(modifier.TokenMetadata)
	// Here please do not call CreateMetadataList() on modifier.TokenMetadata.
	// Metadata element should be manipulated first before creating the metadata list.

	children, err := ft.GetChildrenForDagWithMeta(modifier.ctx, nd, modifier.dagService)
	if err != nil {
		return nil, err
	} else {
		if children == nil || children.MetaNode == nil {
			metadataNotExist = true
		} else {
			mnode = children.MetaNode
		}
	}
	mdagmod, err := modifier.buildMetaDagModifier(nd, mnode, metadataNotExist, dr)
	if err != nil {
		return nil, err
	}

	// Add the given TokenMetadata to the metadata sub-DAG and return
	// a new root of the BTFS file DAG.
	newRoot, err := mdagmod.AddMetadata(nd, dr)
	if err != nil {
		return nil, err
	}

	return newRoot, nil
}

// RemoveMetaAndPin removes the metadata entries from the DAG
// of the given `node` root node of a BTFS file.
func (modifier *MetaModifier) RemoveMetaAndPin(node ipld.Node) (ipld.Node, ipath.Resolved, error) {
	if modifier.Pin {
		modifier.unlocker = modifier.gcLocker.PinLock(modifier.ctx)
	}
	defer func() {
		if modifier.unlocker != nil {
			modifier.unlocker.Unlock(modifier.ctx)
		}
	}()
	p, err := modifier.removeMetaItemsFromFileNode(node)
	if err != nil {
		return nil, nil, err
	}

	root, err := modifier.updateMfs()
	if err != nil {
		return nil, nil, err
	}

	nd, err := root.GetNode()
	if err != nil {
		return nil, nil, err
	}

	if node.Cid() != nd.Cid() {
		err = modifier.pinning.Unpin(modifier.ctx, node.Cid(), true)
		if err != nil {
			return nil, nil, err
		}
	}

	if !modifier.Pin {
		return nd, p, nil
	}
	return nd, p, modifier.PinRoot(modifier.ctx, nd)
}

func (modifier *MetaModifier) removeMetaItemsFromFileNode(node ipld.Node) (ipath.Resolved, error) {
	err := modifier.maybePauseForGC(modifier.ctx)
	if err != nil {
		return nil, err
	}

	dagnode, err := modifier.removeMeta(node)
	if err != nil {
		return nil, err
	}

	// patch it into the root
	return modifier.addNodeToMfs(dagnode, "")
}

// Remove the metadata entries keyed by the given `modifier.TokenMetadata`
// key values from the given `node`.
func (modifier *MetaModifier) removeMeta(nd ipld.Node) (ipld.Node, error) {
	// Error cases
	if nd == nil {
		return nil, errors.New("Invalid argument: nil node value")
	}

	// Build Metadata DAG modifier for balanced format
	// TODO: trickle format
	children, err := ft.GetChildrenForDagWithMeta(modifier.ctx, nd, modifier.dagService)
	if err != nil {
		return nil, err
	} else {
		if children == nil || children.MetaNode == nil {
			return nil, errors.New("no metadata exists for the given node")
		}
	}
	dr := []byte(modifier.TokenMetadata)
	mdagmod, err := modifier.buildMetaDagModifier(nd, children.MetaNode, false, dr)
	if err != nil {
		return nil, err
	}

	// Add the given TokenMetadata to the metadata sub-DAG and return
	// a new root of the BTFS file DAG.
	newRoot, err := mdagmod.RemoveMetadata(nd, dr)
	if err != nil {
		return nil, err
	}

	return newRoot, nil
}

func (modifier *MetaModifier) buildMetaDagModifier(nd ipld.Node, mnode ipld.Node, noMeta bool, dr []byte) (*mod.MetaDagModifier, error) {
	// Retrieve SuperMeta.
	var b []byte
	var err error
	if !noMeta {
		b, err = ftutil.ReadMetadataElementFromDag(modifier.ctx, nd, modifier.dagService, false)
		if err != nil {
			return nil, err
		}
	}

	sm, err := ihelper.GetOrDefaultSuperMeta(b)
	if err != nil {
		return nil, err
	}

	// Build a new Dag modifier.
	dagmod, err := mod.NewDagModifierBalanced(modifier.ctx, mnode, modifier.dagService, chunker.MetaSplitterGen(int64(sm.ChunkSize)), int(sm.MaxLinks), noMeta)
	if err != nil {
		return nil, err
	}

	// Create a Dag builder helper
	params := ihelper.DagBuilderParams{
		Dagserv:    modifier.dagService,
		CidBuilder: modifier.CidBuilder,
		ChunkSize:  sm.ChunkSize,
		Maxlinks:   int(sm.MaxLinks),
	}

	db, err := params.New(chunker.NewMetaSplitter(bytes.NewReader(dr), sm.ChunkSize))
	if err != nil {
		return nil, err
	}

	// Finally create a Meta Dag Modifier
	mdagmod := mod.NewMetaDagModifierBalanced(dagmod, db, modifier.Overwrite)

	return mdagmod, nil
}

func AddMetadataTo(n *core.IpfsNode, skey string, m *ft.Metadata) (string, error) {
	c, err := cid.Decode(skey)
	if err != nil {
		return "", err
	}

	nd, err := n.DAG.Get(n.Context(), c)
	if err != nil {
		return "", err
	}

	mdnode := new(dag.ProtoNode)
	mdata, err := ft.BytesForMetadata(m)
	if err != nil {
		return "", err
	}

	mdnode.SetData(mdata)
	if err := mdnode.AddNodeLink("file", nd); err != nil {
		return "", err
	}

	err = n.DAG.Add(n.Context(), mdnode)
	if err != nil {
		return "", err
	}

	return mdnode.Cid().String(), nil
}

func Metadata(n *core.IpfsNode, skey string) (*ft.Metadata, error) {
	c, err := cid.Decode(skey)
	if err != nil {
		return nil, err
	}

	nd, err := n.DAG.Get(n.Context(), c)
	if err != nil {
		return nil, err
	}

	pbnd, ok := nd.(*dag.ProtoNode)
	if !ok {
		return nil, dag.ErrNotProtobuf
	}

	return ft.MetadataFromBytes(pbnd.Data())
}

// MetaDataMap takes a root node path and returns an unmarshalled metadata map in bytes
func MetaDataMap(ctx context.Context, api coreiface.CoreAPI, path ipath.Path) (map[string]interface{}, error) {
	metaData, err := GetMetaData(ctx, api, path)
	if err != nil {
		return nil, err
	}

	jmap := make(map[string]interface{})
	err = json.Unmarshal(metaData, &jmap)
	if err != nil {
		return nil, err
	}

	return jmap, nil
}

// MetaData takes a root node path and returns the metadata in bytes
func MetaNodeFsData(ctx context.Context, api coreiface.CoreAPI, path ipath.Path) ([]byte, error) {
	nd, err := api.ResolveNode(ctx, path)
	if err != nil {
		return nil, err
	}

	_, ok := nd.(*dag.ProtoNode)
	if !ok {
		return nil, dag.ErrNotProtobuf
	}

	ds := api.Dag()

	mn, err := getMetaData(ctx, nd, ds)
	if err != nil {
		return nil, err
	}

	fsNode, err := ft.FSNodeFromBytes(mn)
	if err != nil {
		return nil, err
	}

	if fsNode.Type() != ft.TTokenMeta {
		return nil, fmt.Errorf("invalid metadata fs node type: %d", fsNode.Type())
	}

	return fsNode.Data(), nil
}

// GetMetaData returns metadata under the given 'nd' node
// if 'nd' has metadata root node aa its first child.
func getMetaData(ctx context.Context, nd ipld.Node, ds ipld.DAGService) ([]byte, error) {
	_, metadata, err := GetDataForUserAndMeta(ctx, nd, ds)
	return metadata, err
}

// GetDataForUserAndMeta returns raw dag node user data, metadata in byte array, and error.
// Note that this function assumes, if token metadata exists inside
// the DAG topped by the given 'node', the 'node' has metadata root
// as first child, user data root as second child.
func GetDataForUserAndMeta(ctx context.Context, nd ipld.Node, ds ipld.DAGService) ([]byte, []byte, error) {
	n := nd.(*dag.ProtoNode)

	fsType, err := ft.GetFSType(n)
	if err != nil {
		return nil, nil, err
	}

	if ft.TTokenMeta == fsType {
		return nil, n.Data(), nil
	}

	// Return user data and metadata if first child is of type TTokenMeta.
	if nd.Links() != nil && len(nd.Links()) >= 2 {
		childen, err := ft.GetChildrenForDagWithMeta(ctx, nd, ds)
		if err != nil {
			return nil, nil, err
		}
		if childen == nil {
			return nil, n.Data(), nil
		}
		return childen.DataNode.(*dag.ProtoNode).Data(), childen.MetaNode.(*dag.ProtoNode).Data(), nil
	}

	return n.Data(), nil, nil
}

// Returns the full metadata DAG data bytes for the DAG of the given file path.
func GetMetaData(ctx context.Context, api coreiface.CoreAPI, path ipath.Path) ([]byte, error) {
	nd, err := api.ResolveNode(ctx, path)
	if err != nil {
		return nil, err
	}

	ds := api.Dag()

	return uio.GetMetaDataFromDagRoot(ctx, nd, ds)
}
