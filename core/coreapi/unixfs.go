package coreapi

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"sync"

	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/coreunix"
	ci "github.com/libp2p/go-libp2p/core/crypto"

	chunker "github.com/TRON-US/go-btfs-chunker"
	files "github.com/TRON-US/go-btfs-files"
	ecies "github.com/TRON-US/go-eccrypto"
	mfs "github.com/TRON-US/go-mfs"
	ft "github.com/TRON-US/go-unixfs"
	unixfile "github.com/TRON-US/go-unixfs/file"
	"github.com/TRON-US/go-unixfs/importer/helpers"
	uio "github.com/TRON-US/go-unixfs/io"
	ftutil "github.com/TRON-US/go-unixfs/util"
	coreiface "github.com/TRON-US/interface-go-btfs-core"
	options "github.com/TRON-US/interface-go-btfs-core/options"
	path "github.com/TRON-US/interface-go-btfs-core/path"

	blockservice "github.com/ipfs/go-blockservice"
	cid "github.com/ipfs/go-cid"
	cidutil "github.com/ipfs/go-cidutil"
	filestore "github.com/ipfs/go-filestore"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	ipld "github.com/ipfs/go-ipld-format"
	dag "github.com/ipfs/go-merkledag"
	dagtest "github.com/ipfs/go-merkledag/test"
	"github.com/ipfs/go-path/resolver"
	"github.com/libp2p/go-libp2p/core/peer"
	// "github.com/prometheus/common/log"
)

type UnixfsAPI CoreAPI

var nilNode *core.IpfsNode
var once sync.Once

func getOrCreateNilNode() (*core.IpfsNode, error) {
	once.Do(func() {
		if nilNode != nil {
			return
		}
		node, err := core.NewNode(context.Background(), &core.BuildCfg{
			//TODO: need this to be true or all files
			// hashed will be stored in memory!
			NilRepo: true,
		})
		if err != nil {
			panic(err)
		}
		nilNode = node
	})

	return nilNode, nil
}

// Add builds a merkledag node from a reader, adds it to the blockstore,
// and returns the key representing that node.
func (api *UnixfsAPI) Add(ctx context.Context, filesNode files.Node, opts ...options.UnixfsAddOption) (path.Resolved, error) {
	settings, prefix, err := options.UnixfsAddOptions(opts...)
	if err != nil {
		return nil, err
	}

	cfg, err := api.repo.Config()
	if err != nil {
		return nil, err
	}

	// check if repo will exceed storage limit if added
	// TODO: this doesn't handle the case if the hashed file is already in blocks (deduplicated)
	// TODO: conditional GC is disabled due to it is somehow not possible to pass the size to the daemon
	//if err := corerepo.ConditionalGC(req.Context(), n, uint64(size)); err != nil {
	//	res.SetError(err, cmds.ErrNormal)
	//	return
	//}

	if settings.NoCopy && !(cfg.Experimental.FilestoreEnabled || cfg.Experimental.UrlstoreEnabled) {
		return nil, fmt.Errorf(`either the filestore or the urlstore must be enabled to use nocopy
		How to enable:
		1.btfs config --json Experimental.FilestoreEnabled true
		2.Then restart your BTFS node to reload your config.
		3.Finally, when adding files with btfs add, pass the --nocopy flag to use the filestore instead of copying the files into your local BTFS repo.`)
	}

	addblockstore := api.blockstore
	if !(settings.FsCache || settings.NoCopy) {
		addblockstore = bstore.NewGCBlockstore(api.baseBlocks, api.blockstore)
	}
	exch := api.exchange
	pinning := api.pinning

	if settings.OnlyHash {
		node, err := getOrCreateNilNode()
		if err != nil {
			return nil, err
		}
		addblockstore = node.Blockstore
		exch = node.Exchange
		pinning = node.Pinning
	}

	bserv := blockservice.New(addblockstore, exch) // hash security 001
	dserv := dag.NewDAGService(bserv)

	// add a sync call to the DagService
	// this ensures that data written to the DagService is persisted to the underlying datastore
	// TODO: propagate the Sync function from the datastore through the blockstore, blockservice and dagservice
	var syncDserv *syncDagService
	if settings.OnlyHash {
		syncDserv = &syncDagService{
			DAGService: dserv,
			syncFn:     func() error { return nil },
		}
	} else {
		syncDserv = &syncDagService{
			DAGService: dserv,
			syncFn: func() error {
				ds := api.repo.Datastore()
				if err := ds.Sync(ctx, bstore.BlockPrefix); err != nil {
					return err
				}
				return ds.Sync(ctx, filestore.FilestorePrefix)
			},
		}
	}

	fileAdder, err := coreunix.NewAdder(ctx, pinning, addblockstore, syncDserv)
	if err != nil {
		return nil, err
	}

	fileAdder.Chunker = settings.Chunker
	if settings.Events != nil {
		fileAdder.Out = settings.Events
		fileAdder.Progress = settings.Progress
	}
	fileAdder.Pin = settings.Pin && !settings.OnlyHash
	fileAdder.Silent = settings.Silent
	fileAdder.RawLeaves = settings.RawLeaves
	fileAdder.NoCopy = settings.NoCopy
	fileAdder.CidBuilder = prefix

	switch settings.Layout {
	case options.BalancedLayout:
		// Default
	case options.TrickleLayout:
		fileAdder.Trickle = true
	default:
		return nil, fmt.Errorf("unknown layout: %d", settings.Layout)
	}

	if settings.Inline {
		fileAdder.CidBuilder = cidutil.InlineBuilder{
			Builder: fileAdder.CidBuilder,
			Limit:   settings.InlineLimit,
		}
	}

	if settings.OnlyHash {
		md := dagtest.Mock()
		emptyDirNode := ft.EmptyDirNode()
		// Use the same prefix for the "empty" MFS root as for the file adder.
		emptyDirNode.SetCidBuilder(fileAdder.CidBuilder)
		mr, err := mfs.NewRoot(ctx, md, emptyDirNode, nil)
		if err != nil {
			return nil, err
		}

		fileAdder.SetMfsRoot(mr)
	}

	if settings.Encrypt {
		pubKey := settings.Pubkey
		if pubKey == "" {
			peerId := settings.PeerId
			if peerId == "" {
				peerId = api.identity.Pretty()
			}
			pubKey, err = peerId2pubkey(peerId)
			if err != nil {
				return nil, err
			}
		}
		// log.Infof("The file will be encrypted with pubkey: %s", settings.Pubkey)
		switch f := filesNode.(type) {
		case files.File:
			bytes, err := ioutil.ReadAll(f)
			if err != nil {
				return nil, err
			}

			ciphertext, metadata, err := ecies.Encrypt(pubKey, bytes)
			if err != nil {
				return nil, err
			}
			m := make(map[string]interface{})
			m["Iv"] = metadata.Iv
			m["EphemPublicKey"] = metadata.EphemPublicKey
			m["Mac"] = metadata.Mac
			m["Mode"] = metadata.Mode
			if err != nil {
				return nil, err
			}

			settings.TokenMetadata, err = api.appendMetaMap(settings.TokenMetadata, m)
			if err != nil {
				return nil, err
			}
			filesNode = files.NewBytesFile([]byte(ciphertext))
		default:
			return nil, notSupport(f)
		}
	}

	if settings.PinDuration != 0 {
		fileAdder.PinDuration = settings.PinDuration
	}
	// This block is intentionally placed here so that
	// any execution case can append metadata
	if settings.TokenMetadata != "" {
		fileAdder.TokenMetadata = settings.TokenMetadata
		if _, ok := filesNode.(files.Directory); ok {
			fileAdder.MetaForDirectory = true
		}
	}

	// if chunker is reed-solomon and the given `node` is directory,
	// create ReedSolomonAdder. Otherwise use Adder.
	var nd ipld.Node
	dir, ok := filesNode.(files.Directory)
	if ok && chunker.IsReedSolomon(settings.Chunker) {
		rsfileAdder, err := coreunix.NewReedSolomonAdder(fileAdder)
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
		nd, err = rsfileAdder.AddAllAndPin(ctx, filesNode)
	} else {
		nd, err = fileAdder.AddAllAndPin(ctx, filesNode)
	}
	if err != nil {
		return nil, err
	}
	if nd == nil {
		return nil, errors.New("unexpected nil value for ipld.Node")
	}

	if !settings.OnlyHash {
		if err := api.provider.Provide(nd.Cid()); err != nil {
			return nil, err
		}
	}

	return path.IpfsPath(nd.Cid()), nil
}

func notSupport(f interface{}) error {
	return fmt.Errorf("not support: %v", f)
}

func peerId2pubkey(peerId string) (string, error) {
	id, err := peer.Decode(peerId)
	if err != nil {
		return "", err
	}
	publicKey, err := id.ExtractPublicKey()
	if err != nil {
		return "", err
	}
	bytes, err := ci.MarshalPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes[4:]), nil
}

func (api *UnixfsAPI) Get(ctx context.Context, p path.Path, opts ...options.UnixfsGetOption) (files.Node, error) {
	settings, err := options.UnixfsGetOptions(opts...)
	if err != nil {
		return nil, err
	}

	ses := api.core().getSession(ctx)

	nd, err := ses.ResolveNode(ctx, p)
	if err != nil {
		if e, ok := err.(resolver.ErrNoLink); ok {
			file, err := api.getReedSolomonFile(ctx, settings, p)
			if err != nil {
				return nil, e
			}
			return file, nil
		}
		return nil, err
	}

	var node files.Node
	if !settings.Metadata && settings.Decrypt {
		node, err = unixfile.NewUnixfsFile(ctx, ses.dag, nd,
			unixfile.UnixfsFileOptions{RepairShards: settings.Repairs})
		switch f := node.(type) {
		case files.File:
			privKey, err := api.getPrivateKey(settings.PrivateKey)
			if err != nil {
				return nil, err
			}

			mbytes, err := api.GetMetadata(ctx, p)
			if err != nil {
				return nil, err
			}

			t := &ecies.EciesMetadata{}
			err = json.Unmarshal(mbytes, t)
			if err != nil {
				return nil, err
			}

			bytes, err := ioutil.ReadAll(f)
			if err != nil {
				return nil, err
			}

			s, err := ecies.Decrypt(privKey, string(bytes), t)
			if err != nil {
				return nil, err
			}
			node = files.NewBytesFile([]byte(s))
		default:
			return nil, notSupport(f)
		}
	} else {
		var ds ipld.DAGService
		// If needing repair, the dag has to be writable
		if settings.Repairs != nil {
			ds = api.dag
		} else {
			ds = ses.dag
		}
		node, err = unixfile.NewUnixfsFile(ctx, ds, nd,
			unixfile.UnixfsFileOptions{Meta: settings.Metadata, RepairShards: settings.Repairs})
		if settings.Metadata {
			f, ok := node.(files.File)
			if !ok {
				return nil, notSupport(f)
			}
			inb, err := ftutil.ReadMetadataElementFromDag(ctx, nd, api.dag, settings.Metadata)
			if err != nil {
				return nil, err
			}

			m := make(map[string]interface{})
			err = json.Unmarshal(inb, &m)
			if err != nil {
				return nil, err
			}

			// Delete SuperMeta entries from the map.
			var superMeta helpers.SuperMeta
			typ := reflect.TypeOf(superMeta)
			for i := 0; i < typ.NumField(); i++ {
				delete(m, typ.Field(i).Name)
			}

			outb, err := json.Marshal(m)
			if err != nil {
				return nil, err
			}
			node = files.NewBytesFile(outb)
		}
	}

	return node, err
}

func (api *UnixfsAPI) getReedSolomonFile(ctx context.Context, settings *options.UnixfsGetSettings,
	p path.Path) (n files.Node, e error) {
	defer func() {
		if err := recover(); err != nil {
			n, e = nil, fmt.Errorf("panic:%v", err)
		}
	}()
	root, paths, err := splitPath(p)
	if err != nil {
		return nil, err
	}
	ses := api.core().getSession(ctx)
	var ds ipld.DAGService
	// If needing repair, the dag has to be writable
	if settings.Repairs != nil {
		ds = api.dag
	} else {
		ds = ses.dag
	}
	nd, err := ses.ResolveNode(ctx, root)
	if err != nil {
		return nil, err
	}
	node, err := unixfile.NewUnixfsFile(ctx, ds, nd,
		unixfile.UnixfsFileOptions{Meta: settings.Metadata, RepairShards: settings.Repairs})
	dir := node.(*unixfile.RsDirectory)
	dit := dir.Entries()
	if !dit.Next() {
		return nil, errors.New("not found")
	}
	for i := 0; i < len(paths); i++ {
		if i == len(paths)-1 {
			for dit.Name() != paths[i] {
				if !dit.Next() {
					return nil, errors.New("not found")
				}
			}
			return dit.Node(), nil
		}
		for dit.Name() != paths[i] {
			if !dit.Next() {
				return nil, errors.New("not found")
			}
		}
		_t := dit.Node().(*unixfile.RsDirectory)
		dit = _t.Entries()
		if !dit.Next() {
			return nil, errors.New("not found")
		}
	}
	return dit.Node(), nil
}

func splitPath(p path.Path) (root path.Path, paths []string, err error) {
	split := strings.Split(p.String(), "/")
	// /btfs/Qm/a/b => "","btfs","Qm
	if len(split) <= 3 || split[1] != "btfs" || split[3] == "" {
		err = errors.New("invalid path")
		return
	}
	root = path.New("/btfs/" + split[2])
	for i := 3; i < len(split); i++ {
		paths = append(paths, split[i])
	}
	return
}

// compatible with base64, 32bits-hex and 36bits-hex
func (api *UnixfsAPI) getPrivateKey(input string) (string, error) {
	privKey := input
	var bytes []byte
	var err error
	if privKey == "" {
		bytes, err = ci.MarshalPrivateKey(api.privateKey)
		if err != nil {
			return "", err
		}
	} else if isPath(privKey) {
		bytes, err = ioutil.ReadFile(privKey)
		if err != nil {
			return "", err
		}
		bytes, err = decodePrivateKey(string(bytes))
		if err != nil {
			return "", err
		}
	} else {
		bytes, err = decodePrivateKey(privKey)
		if err != nil {
			return "", err
		}
	}
	if len(bytes) < 32 {
		return "", fmt.Errorf("invalid private key")
	}
	privKey = hex.EncodeToString(bytes[4:])
	return privKey, nil
}

func isPath(path string) bool {
	if path == "" {
		return false
	}
	fc := path[0]
	return fc == '.' || fc == '/' || fc == '~'
}

func decodePrivateKey(privKey string) ([]byte, error) {
	bytes, err := base64.StdEncoding.DecodeString(privKey)
	if err != nil {
		bytes, err = hex.DecodeString(privKey)
		if err != nil {
			return nil, err
		}
	}
	return bytes, nil
}

// Ls returns the contents of an BTFS or BTNS object(s) at path p, with the format:
// `<link base58 hash> <link size in bytes> <link name>`
func (api *UnixfsAPI) Ls(ctx context.Context, p path.Path, opts ...options.UnixfsLsOption) (<-chan coreiface.DirEntry, error) {
	settings, err := options.UnixfsLsOptions(opts...)
	if err != nil {
		return nil, err
	}

	ses := api.core().getSession(ctx)
	uses := (*UnixfsAPI)(ses)

	dagnode, err := ses.ResolveNode(ctx, p)
	if err != nil {
		return nil, err
	}

	dir, err := uio.NewDirectoryFromNode(ses.dag, dagnode)
	if err == uio.ErrNotADir {
		return uses.lsFromLinks(ctx, dagnode.Links(), settings)
	}
	if err != nil {
		return nil, err
	}

	return uses.lsFromLinksAsync(ctx, dir, settings)
}

func (api *UnixfsAPI) AddMetadata(ctx context.Context, p path.Path, m string, opts ...options.UnixfsAddMetaOption) (path.Resolved, error) {
	settings, err := options.UnixfsAddMetaOptions(opts...)
	if err != nil {
		return nil, err
	}

	addblockstore := api.blockstore
	exch := api.exchange
	pinning := api.pinning

	bserv := blockservice.New(addblockstore, exch) // hash security 001
	dserv := dag.NewDAGService(bserv)

	metaModifier, err := coreunix.NewMetaModifier(ctx, pinning, addblockstore, dserv)
	if err != nil {
		return nil, err
	}

	if m != "" {
		metaModifier.TokenMetadata = m
	}

	metaModifier.Pin = settings.Pin
	metaModifier.Overwrite = settings.Overwrite

	ses := api.core().getSession(ctx)

	nd, err := ses.ResolveNode(ctx, p)
	if err != nil {
		return nil, err
	}

	newNode, rp, err := metaModifier.AddMetaAndPin(nd)
	if err != nil {
		return nil, err
	}

	if err := api.provider.Provide(newNode.Cid()); err != nil {
		return nil, err
	}

	// Return a new /btfs resolved path from the CID.
	return rp, nil
}

func (api *UnixfsAPI) RemoveMetadata(ctx context.Context, p path.Path, m string, opts ...options.UnixfsRemoveMetaOption) (path.Resolved, error) {
	settings, err := options.UnixfsRemoveMetaOptions(opts...)
	if err != nil {
		return nil, err
	}

	addblockstore := api.blockstore
	exch := api.exchange
	pinning := api.pinning

	bserv := blockservice.New(addblockstore, exch) // hash security 001
	dserv := dag.NewDAGService(bserv)

	metaModifier, err := coreunix.NewMetaModifier(ctx, pinning, addblockstore, dserv)
	if err != nil {
		return nil, err
	}

	if m != "" {
		metaModifier.TokenMetadata = m
	}

	metaModifier.Pin = settings.Pin

	ses := api.core().getSession(ctx)

	nd, err := ses.ResolveNode(ctx, p)
	if err != nil {
		return nil, err
	}

	newNode, rp, err := metaModifier.RemoveMetaAndPin(nd)
	if err != nil {
		return nil, err
	}

	if err := api.provider.Provide(newNode.Cid()); err != nil {
		return nil, err
	}

	// Return a new /btfs resolved path from the CID.
	return rp, nil
}

func (api *UnixfsAPI) processLink(ctx context.Context, linkres ft.LinkResult, settings *options.UnixfsLsSettings) coreiface.DirEntry {
	lnk := coreiface.DirEntry{
		Name: linkres.Link.Name,
		Cid:  linkres.Link.Cid,
		Err:  linkres.Err,
	}
	if lnk.Err != nil {
		return lnk
	}

	switch lnk.Cid.Type() {
	case cid.Raw:
		// No need to check with raw leaves
		lnk.Type = coreiface.TFile
		lnk.Size = linkres.Link.Size
	case cid.DagProtobuf:
		if !settings.ResolveChildren {
			break
		}

		linkNode, err := linkres.Link.GetNode(ctx, api.dag)
		if err != nil {
			lnk.Err = err
			break
		}

		if pn, ok := linkNode.(*dag.ProtoNode); ok {
			d, err := ft.FSNodeFromBytes(pn.Data())
			if err != nil {
				lnk.Err = err
				break
			}
			switch d.Type() {
			case ft.TFile, ft.TRaw:
				lnk.Type = coreiface.TFile
			case ft.THAMTShard, ft.TDirectory, ft.TMetadata:
				lnk.Type = coreiface.TDirectory
			case ft.TSymlink:
				lnk.Type = coreiface.TSymlink
				lnk.Target = string(d.Data())
			}
			lnk.Size = d.FileSize()
		}
	}

	return lnk
}

func (api *UnixfsAPI) lsFromLinksAsync(ctx context.Context, dir uio.Directory, settings *options.UnixfsLsSettings) (<-chan coreiface.DirEntry, error) {
	out := make(chan coreiface.DirEntry)

	go func() {
		defer close(out)
		for l := range dir.EnumLinksAsync(ctx) {
			select {
			case out <- api.processLink(ctx, l, settings): //TODO: perf: processing can be done in background and in parallel
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

func (api *UnixfsAPI) lsFromLinks(ctx context.Context, ndlinks []*ipld.Link, settings *options.UnixfsLsSettings) (<-chan coreiface.DirEntry, error) {
	links := make(chan coreiface.DirEntry, len(ndlinks))
	for _, l := range ndlinks {
		lr := ft.LinkResult{Link: &ipld.Link{Name: l.Name, Size: l.Size, Cid: l.Cid}}

		links <- api.processLink(ctx, lr, settings) //TODO: can be parallel if settings.Async
	}
	close(links)
	return links, nil
}

func (api *UnixfsAPI) core() *CoreAPI {
	return (*CoreAPI)(api)
}

func (api *UnixfsAPI) GetMetadata(ctx context.Context, p path.Path) ([]byte, error) {
	f, err := api.core().Unixfs().Get(ctx, p, options.Unixfs.Metadata(true))
	if err != nil {
		return nil, err
	}
	switch f.(type) {
	case files.File:
		bytes, err := ioutil.ReadAll(f.(files.File))
		if err != nil {
			return nil, err
		}
		return bytes, nil
	default:
		return nil, notSupport(f)
	}
}

func (api *UnixfsAPI) appendMetaMap(tokenMetadata string, metaMap map[string]interface{}) (string, error) {
	if metaMap == nil {
		return "", nil
	}

	var b []byte = nil
	var err error = nil
	if tokenMetadata != "" {
		tmp := make(map[string]interface{})
		err = json.Unmarshal([]byte(tokenMetadata), &tmp)
		if err != nil {
			return "", err
		}
		for k, v := range metaMap {
			tmp[k] = v
		}
		b, err = json.Marshal(tmp)
		if err != nil {
			return "", err
		}
	} else {
		b, err = json.Marshal(metaMap)
		if err != nil {
			return "", err
		}
	}

	tokenMetadata = string(b)
	return tokenMetadata, nil
}

// syncDagService is used by the Adder to ensure blocks get persisted to the underlying datastore
type syncDagService struct {
	ipld.DAGService
	syncFn func() error
}

func (s *syncDagService) Sync() error {
	return s.syncFn()
}
