package coreunix

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	gopath "path"

	"container/list"
	"encoding/json"

	files "github.com/TRON-US/go-btfs-files"
	"github.com/TRON-US/go-unixfs"
	uio "github.com/TRON-US/go-unixfs/io"
	ipld "github.com/ipfs/go-ipld-format"
	dag "github.com/ipfs/go-merkledag"
)

const (
	SerialFile = iota + 1
	MultipartFile
	SliceFile
)

type ReedSolomonAdder struct {
	Adder
	FileType      int
	InfileReaders []io.Reader
	IsDir         bool
	CurrentOffset uint64
}

// NewReedSolomonAdder returns a new ReedSolomonAdder used for a file/directory add operation.
func NewReedSolomonAdder(adder *Adder) (*ReedSolomonAdder, error) {
	rsa := &ReedSolomonAdder{
		InfileReaders: nil,
		IsDir:         false,
		CurrentOffset: 0,
	}
	rsa.Adder = *adder

	return rsa, nil
}

// AddAllAndPin adds the given request's files and pin them.
func (rsadder *ReedSolomonAdder) AddAllAndPin(ctx context.Context, file files.Node) (ipld.Node, error) {
	if rsadder.Pin {
		rsadder.unlocker = rsadder.gcLocker.PinLock(ctx)
	}
	defer func() {
		if rsadder.unlocker != nil {
			rsadder.unlocker.Unlock(ctx)
		}
	}()

	// Create a directory tree from the input `file` directory.
	fileList := list.New()
	defer func() {
		for ent := fileList.Front(); ent != nil; ent = ent.Next() {
			if node, ok := ent.Value.(files.Node); ok {
				node.Close()
			}
		}
	}()

	n, err := rsadder.addFileNode(ctx, "", file, fileList, true)
	if err != nil {
		return nil, err
	}
	var byts []byte
	switch node := n.(type) {
	case *uio.DirNode, *uio.SymlinkNode, *uio.FileNode:
		byts, err = json.Marshal(node)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unexpected node type")
	}

	var nd ipld.Node
	if n.Path() == "" && n.NodeSize() == 0 {
		nd, err = rsadder.addToMfs(ctx, file)
		if err != nil {
			return nil, err
		}
	} else {
		var reader io.Reader = io.MultiReader(rsadder.InfileReaders...)
		if rsadder.Progress {
			rdr := &progressReader{file: reader, path: "", out: rsadder.Out}
			if file, ok := file.(files.Node); ok {
				reader = &progressReader3{rdr, file}
			} else {
				reader = rdr
			}
		}

		// Create a DAG with the above directory tree as metadata and
		// the data from the given `file` directory or file.
		nd, err = rsadder.add(reader, byts)
		if err != nil {
			return nil, err
		}

		// output directory and file events
		err = rsadder.outputDirs("", nd)
		if err != nil {
			return nil, err
		}
	}

	// Pin the newly created DAG.
	if !rsadder.Pin {
		return nd, nil
	}

	if err := rsadder.PinRoot(ctx, nd); err != nil {
		return nil, err
	}

	return nd, nil
}

// addFileNode traverses the directory tree under
// the given `file` in a bottom up DFS way.
func (rsadder *ReedSolomonAdder) addFileNode(ctx context.Context, path string, file files.Node, fList *list.List, toplevel bool) (uio.Node, error) {
	defer func() {
		fList.PushFront(file)
	}()
	err := rsadder.maybePauseForGC(ctx)
	if err != nil {
		return nil, err
	}

	if rsadder.liveNodes >= liveCacheSize {
		// TODO: flush free memory rsadder.mfsRoot()'s FlushMemFree() to flush free memory
		log.Info("rsadder liveNodes >= liveCacheSize, needs flushing")
	}
	rsadder.liveNodes++

	switch f := file.(type) {
	case files.Directory:
		return rsadder.addDir(ctx, path, f, fList, toplevel)
	case *files.Symlink:
		return rsadder.addSymlink(path, f)
	case files.File:
		return rsadder.addFile(path, f)
	default:
		return nil, errors.New("unknown file type")
	}
}

func (rsadder *ReedSolomonAdder) addDir(ctx context.Context, path string, dir files.Directory, fList *list.List, toplevel bool) (uio.Node, error) {
	log.Infof("adding directory: %s", path)

	_, dstName := gopath.Split(path)
	node := &uio.DirNode{
		BaseNode: uio.BaseNode{
			NodeType: uio.DirNodeType,
			NodePath: path,
			NodeName: dstName},
		// Kept the following line for possible use cases than `btfs get`
		//Directory: dir,
	}

	it := dir.Entries()
	var size uint64
	for it.Next() {
		fpath := gopath.Join(path, it.Name())
		child, err := rsadder.addFileNode(ctx, fpath, it.Node(), fList, false)
		if err != nil {
			return nil, err
		}
		node.Links = append(node.Links, child)
		size += uint64(child.NodeSize())
	}
	node.Siz = size
	if toplevel && rsadder.FileType == MultipartFile {
		err := dir.SetSize(int64(size))
		if err != nil {
			return nil, err
		}
	}
	return node, it.Err()
}

func (rsadder *ReedSolomonAdder) addSymlink(path string, l *files.Symlink) (uio.Node, error) {
	_, dstName := gopath.Split(path)
	return &uio.SymlinkNode{
		BaseNode: uio.BaseNode{
			NodeType: uio.SymlinkNodeType,
			NodePath: path,
			NodeName: dstName},
		Data: l.Target,
		// Kept the following line for the same reason as DirNode.
		//Symlink:   *l,
	}, nil
}

func getReader(r io.Reader) (io.Reader, int64, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, 0, err
	}

	fSize := int64(len(buf))
	r = bytes.NewReader(buf)
	return r, fSize, nil
}

func (rsadder *ReedSolomonAdder) addFile(path string, file files.File) (uio.Node, error) {
	_, dstName := gopath.Split(path)
	var (
		reader io.Reader = file
		fSize  int64
		err    error
	)

	switch rsadder.FileType {
	case MultipartFile:
		rf, ok := file.(*files.ReaderFile)
		if !ok {
			return nil, errors.New("expected multipartFile file type")
		}

		reader, fSize, err = getReader(rf.Reader())
		if err != nil {
			return nil, err
		}
	case SerialFile, SliceFile:
		fSize, err = file.Size()
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported file type")
	}

	rsadder.InfileReaders = append(rsadder.InfileReaders, reader)
	currentOffset := rsadder.CurrentOffset
	rsadder.CurrentOffset += uint64(fSize)

	node := &uio.FileNode{
		BaseNode: uio.BaseNode{
			NodeType:    uio.FileNodeType,
			NodePath:    path,
			NodeName:    dstName,
			Siz:         uint64(fSize),
			StartOffset: currentOffset,
		},
		// Kept the following line for the same reason as DirNode.
		//File: file,
	}
	return node, nil
}

// outputDirs outputs file and directory dagnodes in a postorder DFS fashion.
func (rsadder *ReedSolomonAdder) outputDirs(path string, n ipld.Node) error {
	dn, ok := n.(*dag.ProtoNode)
	if !ok {
		return errors.New("expected protobuf node type")
	}
	fsn, err := unixfs.FSNodeFromBytes(dn.Data())
	if err != nil {
		return err
	}

	if fsn.IsDir() {
		for _, link := range dn.Links() {
			child, err := link.GetNode(rsadder.ctx, rsadder.dagService)
			if err != nil {
				return err
			}
			childpath := gopath.Join(path, link.Name)
			err = rsadder.outputDirs(childpath, child)
			if err != nil {
				return err
			}
		}
	} else {
		return outputDagnode(rsadder.Out, path, n)
	}
	return nil
}
