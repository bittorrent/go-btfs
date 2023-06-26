package cmdenv

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	gopath "path"
	"strings"

	cmds "github.com/bittorrent/go-btfs-cmds"
	config "github.com/bittorrent/go-btfs-config"
	files "github.com/bittorrent/go-btfs-files"
	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/corerepo"

	coreiface "github.com/bittorrent/interface-go-btfs-core"
	"github.com/bittorrent/interface-go-btfs-core/options"
	"github.com/bittorrent/interface-go-btfs-core/path"
	"github.com/ipfs/go-cid"
)

// DefaultBufSize is the buffer size for gets. for now, 1MB, which is ~4 blocks.
// TODO: does this need to be configurable?
var DefaultBufSize = 1048576
var ErrInvalidCompressionLevel = errors.New("compression level must be between 1 and 9")

// GetFileArg returns the next file from the directory or an error
func GetFileArg(it files.DirIterator) (files.File, error) {
	if !it.Next() {
		err := it.Err()
		if err == nil {
			err = fmt.Errorf("expected a file argument")
		}
		return nil, err
	}
	file := files.FileFromEntry(it)
	if file == nil {
		return nil, fmt.Errorf("file argument was nil")
	}
	return file, nil
}

type PinOutput struct {
	Pins []string
}

func DownloadAndRebuildFile(req *cmds.Request, res cmds.ResponseEmitter, api coreiface.CoreAPI, fileHash string, lostShards string) error {
	_, err := GetFile(req, res, api, fileHash, false, "", false, lostShards, false, false, false, 0)
	return err
}

func GetFile(req *cmds.Request, res cmds.ResponseEmitter, api coreiface.CoreAPI, btfsPath string, decrypt bool,
	privateKey string, meta bool, repairShards string, quiet bool, archive bool, cmprs bool, cmplvl int) (io.Reader, error) {

	var repairs []cid.Cid
	if repairShards != "" {
		rshards := strings.Split(repairShards, ",")
		for _, rshard := range rshards {
			rcid, err := cid.Parse(rshard)
			if err != nil {
				return nil, err
			}
			repairs = append(repairs, rcid)
		}
	}

	opts := []options.UnixfsGetOption{
		options.Unixfs.Decrypt(decrypt),
		options.Unixfs.PrivateKey(privateKey),
		options.Unixfs.Metadata(meta),
		options.Unixfs.Repairs(repairs),
	}

	p := path.New(btfsPath)
	file, err := api.Unixfs().Get(req.Context, p, opts...)
	if err != nil {
		return nil, err
	}

	if quiet {
		return nil, nil
	}

	size, err := file.Size()
	if err != nil {
		return nil, err
	}
	res.SetLength(uint64(size))

	cmplvl, err = GetCompressLevel(cmprs, cmplvl)
	if err != nil {
		return nil, err
	}
	reader, err := fileArchive(file, p.String(), archive, cmplvl)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func GetCompressLevel(cmprs bool, cmplvl int) (int, error) {
	switch {
	case !cmprs:
		return gzip.NoCompression, nil
	case cmprs && cmplvl == 0:
		return gzip.DefaultCompression, nil
	case cmprs && (cmplvl < 1 || cmplvl > 9):
		return gzip.NoCompression, ErrInvalidCompressionLevel
	}
	return cmplvl, nil
}

type identityWriteCloser struct {
	w io.Writer
}

func (i *identityWriteCloser) Write(p []byte) (int, error) {
	return i.w.Write(p)
}

func (i *identityWriteCloser) Close() error {
	return nil
}

func fileArchive(f files.Node, name string, archive bool, compression int) (io.Reader, error) {
	cleaned := gopath.Clean(name)
	_, filename := gopath.Split(cleaned)

	// need to connect a writer to a reader
	piper, pipew := io.Pipe()
	checkErrAndClosePipe := func(err error) bool {
		if err != nil {
			_ = pipew.CloseWithError(err)
			return true
		}
		return false
	}

	// use a buffered writer to parallelize task
	bufw := bufio.NewWriterSize(pipew, DefaultBufSize)

	// compression determines whether to use gzip compression.
	maybeGzw, err := newMaybeGzWriter(bufw, compression)
	if checkErrAndClosePipe(err) {
		return nil, err
	}

	closeGzwAndPipe := func() {
		if err := maybeGzw.Close(); checkErrAndClosePipe(err) {
			return
		}
		if err := bufw.Flush(); checkErrAndClosePipe(err) {
			return
		}
		pipew.Close() // everything seems to be ok.
	}

	if !archive && compression != gzip.NoCompression {
		// the case when the node is a file
		r := files.ToFile(f)
		if r == nil {
			return nil, errors.New("file is not regular")
		}

		go func() {
			if _, err := io.Copy(maybeGzw, r); checkErrAndClosePipe(err) {
				return
			}
			closeGzwAndPipe() // everything seems to be ok
		}()
	} else {
		// the case for 1. archive, and 2. not archived and not compressed, in which tar is used anyway as a transport format

		// construct the tar writer
		w, err := files.NewTarWriter(maybeGzw)
		if checkErrAndClosePipe(err) {
			return nil, err
		}

		go func() {
			// write all the nodes recursively
			if err := w.WriteFile(f, filename); checkErrAndClosePipe(err) {
				return
			}
			w.Close()         // close tar writer
			closeGzwAndPipe() // everything seems to be ok
		}()
	}

	return piper, nil
}

func newMaybeGzWriter(w io.Writer, compression int) (io.WriteCloser, error) {
	if compression != gzip.NoCompression {
		return gzip.NewWriterLevel(w, compression)
	}
	return &identityWriteCloser{w}, nil
}

func RmPin(req *cmds.Request, res cmds.ResponseEmitter, n *core.IpfsNode, api coreiface.CoreAPI, cfg *config.Config, arguments []string, recursive bool, force bool) error {
	enc, err := GetCidEncoder(req)
	if err != nil {
		return err
	}

	pins := make([]string, 0, len(req.Arguments))
	for _, b := range arguments {
		rp, err := api.ResolvePath(req.Context, path.New(b))
		if err != nil {
			return err
		}

		id := enc.Encode(rp.Cid())
		pins = append(pins, id)
		if err := api.Pin().Rm(req.Context, rp,
			options.Pin.RmRecursive(recursive), options.Pin.RmForce(force)); err != nil {
			return err
		}
	}

	if err := cmds.EmitOnce(res, &PinOutput{pins}); err != nil {
		return err
	}

	if cfg.Experimental.RemoveOnUnpin {
		corerepo.RepoGc(req, res, n, false)
	}

	return nil
}
