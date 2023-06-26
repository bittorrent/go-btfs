package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/bittorrent/go-btfs/core/commands/cmdenv"

	cmds "github.com/bittorrent/go-btfs-cmds"
	files "github.com/bittorrent/go-btfs-files"
	iface "github.com/bittorrent/interface-go-btfs-core"
	"github.com/bittorrent/interface-go-btfs-core/options"
	"github.com/bittorrent/interface-go-btfs-core/path"
)

const (
	progressBarMinSize       = 1024 * 1024 * 8 // show progress bar for outputs > 8MiB
	offsetOptionName         = "offset"
	lengthOptionName         = "length"
	catMetaDisplayOptionName = "meta"
)

var CatCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Show BTFS object data.",
		ShortDescription: "Displays the data contained by an BTFS or BTNS object(s) at the given path.",
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("btfs-path", true, true, "The path to the BTFS object(s) to be outputted.").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.Int64Option(offsetOptionName, "o", "Byte offset to begin reading from."),
		cmds.Int64Option(lengthOptionName, "l", "Maximum number of bytes to read."),
		cmds.BoolOption(catMetaDisplayOptionName, "m", "Display token metadata"),
		cmds.BoolOption(decryptName, "d", "Decrypt the file."),
		cmds.StringOption(privateKeyName, "pk", "The private key to decrypt file."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		api, err := cmdenv.GetApi(env, req)
		if err != nil {
			return err
		}

		offset, _ := req.Options[offsetOptionName].(int64)
		if offset < 0 {
			return fmt.Errorf("cannot specify negative offset")
		}

		max, found := req.Options[lengthOptionName].(int64)

		if max < 0 {
			return fmt.Errorf("cannot specify negative length")
		}
		if !found {
			max = -1
		}

		meta, _ := req.Options[catMetaDisplayOptionName].(bool)

		err = req.ParseBodyArgs()
		if err != nil {
			return err
		}

		readers, length, err := cat(req.Context, api, req.Arguments, int64(offset), int64(max), meta, req.Options)
		if err != nil {
			return err
		}

		/*
			if err := corerepo.ConditionalGC(req.Context, node, length); err != nil {
				re.SetError(err, cmds.ErrNormal)
				return
			}
		*/

		res.SetLength(length)
		reader := io.MultiReader(readers...)

		// Since the reader returns the error that a block is missing, and that error is
		// returned from io.Copy inside Emit, we need to take Emit errors and send
		// them to the client. Usually we don't do that because it means the connection
		// is broken or we supplied an illegal argument etc.
		return res.Emit(reader)
	},
	PostRun: cmds.PostRunMap{
		cmds.CLI: func(res cmds.Response, re cmds.ResponseEmitter) error {
			if res.Length() > 0 && res.Length() < progressBarMinSize {
				return cmds.Copy(re, res)
			}

			for {
				v, err := res.Next()
				if err != nil {
					if err == io.EOF {
						return nil
					}
					return err
				}

				switch val := v.(type) {
				case io.Reader:
					bar, reader := progressBarForReader(os.Stderr, val, int64(res.Length()))
					bar.Start()

					err = re.Emit(reader)
					if err != nil {
						return err
					}
				default:
					log.Warnf("cat postrun: received unexpected type %T", val)
				}
			}
		},
	},
}

func cat(ctx context.Context, api iface.CoreAPI, paths []string, offset int64, max int64, meta bool, opts cmds.OptMap) ([]io.Reader, uint64, error) {
	readers := make([]io.Reader, 0, len(paths))
	length := uint64(0)
	if max == 0 {
		return nil, 0, nil
	}

	if opts[decryptName] == nil {
		opts[decryptName] = false
	}
	if opts[privateKeyName] == nil {
		opts[privateKeyName] = ""
	}
	getOptions := []options.UnixfsGetOption{
		options.Unixfs.Decrypt(opts[decryptName].(bool)),
		options.Unixfs.PrivateKey(opts[privateKeyName].(string)),
		options.Unixfs.Metadata(meta),
	}

	for _, p := range paths {
		f, err := api.Unixfs().Get(ctx, path.New(p), getOptions...)
		if err != nil {
			return nil, 0, err
		}

		var file files.File
		switch f := f.(type) {
		case files.File:
			file = f
		case files.Directory:
			return nil, 0, iface.ErrIsDir
		default:
			return nil, 0, iface.ErrNotSupported
		}

		fsize, err := file.Size()
		if err != nil {
			return nil, 0, err
		}

		if offset > fsize {
			offset = offset - fsize
			continue
		}

		var count int64
		var seekNotAllowed bool
		if opts[decryptName] != nil && opts[decryptName].(bool) {
			seekNotAllowed = true
		}
		if !seekNotAllowed && !meta {
			count, err = file.Seek(offset, io.SeekStart)
			if err != nil {
				return nil, 0, err
			}
		}

		offset = 0

		fsize, err = file.Size()
		if err != nil {
			return nil, 0, err
		}

		size := uint64(fsize - count)
		length += size
		if max > 0 && length >= uint64(max) {
			var r io.Reader = file
			if overshoot := int64(length - uint64(max)); overshoot != 0 {
				r = io.LimitReader(file, int64(size)-overshoot)
				length = uint64(max)
			}
			readers = append(readers, r)
			break
		}
		readers = append(readers, file)
	}
	return readers, length, nil
}
