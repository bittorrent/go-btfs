package commands

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	cmds "github.com/TRON-US/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/e"

	"github.com/whyrusleeping/tar-utils"
	"gopkg.in/cheggaaa/pb.v1"
)

const (
	outputOptionName           = "output"
	archiveOptionName          = "archive"
	compressOptionName         = "compress"
	compressionLevelOptionName = "compression-level"
	getMetaDisplayOptionName   = "meta"
	decryptName                = "decrypt"
	privateKeyName             = "private-key"
	repairShardsName           = "repair-shards"
)

var GetCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Download BTFS objects.",
		ShortDescription: `
Stores to disk the data contained an BTFS or BTNS object(s) at the given path.

By default, the output will be stored at './<btfs-path>', but an alternate
path can be specified with '--output=<path>' or '-o=<path>'.

To output a TAR archive instead of unpacked files, use '--archive' or '-a'.

To compress the output with GZIP compression, use '--compress' or '-C'. You
may also specify the level of compression by specifying '-l=<1-9>'.

To output just the metadata of this BTFS node, use '--meta' or '-m'.

To repair missing shards of a Reed-Solomon encoded file, use '--repair-shards' or '-rs'.
If '--meta' or '-m' is enabled, this option is ignored.
`,
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("btfs-path", true, false, "The path to the BTFS object(s) to be outputted.").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.StringOption(outputOptionName, "o", "The path where the output should be stored."),
		cmds.BoolOption(archiveOptionName, "a", "Output a TAR archive."),
		cmds.BoolOption(compressOptionName, "C", "Compress the output with GZIP compression."),
		cmds.IntOption(compressionLevelOptionName, "l", "The level of compression (1-9)."),
		cmds.BoolOption(getMetaDisplayOptionName, "m", "Display token metadata"),
		cmds.BoolOption(decryptName, "d", "Decrypt the file."),
		cmds.StringOption(privateKeyName, "pk", "The private key to decrypt file."),
		cmds.StringOption(repairShardsName, "rs", "Repair the list of shards. Multihashes separated by ','."),
		cmds.BoolOption(quietOptionName, "q", "Quiet mode: perform get operation without writing to anywhere. Same as using -o /dev/null."),
	},
	PreRun: func(req *cmds.Request, env cmds.Environment) error {
		_, err := cmdenv.GetCompressLevel(getCompressOptions(req))
		return err
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		api, err := cmdenv.GetApi(env, req)
		if err != nil {
			return err
		}

		btfsPath := req.Arguments[0]
		decrypt, _ := req.Options[decryptName].(bool)
		privateKey, _ := req.Options[privateKeyName].(string)
		meta, _ := req.Options[getMetaDisplayOptionName].(bool)
		repairShards, _ := req.Options[repairShardsName].(string)
		quiet, _ := req.Options[quietOptionName].(bool)
		archive, _ := req.Options[archiveOptionName].(bool)
		cmprs, cmplvl := getCompressOptions(req)

		reader, err := cmdenv.GetFile(req, res, api, btfsPath, decrypt, privateKey, meta, repairShards, quiet, archive, cmprs, cmplvl)
		if err != nil {
			return err
		}

		return res.Emit(reader)
	},
	PostRun: cmds.PostRunMap{
		cmds.CLI: func(res cmds.Response, re cmds.ResponseEmitter) error {
			req := res.Request()

			v, err := res.Next()
			if err != nil {
				return err
			}

			// quiet mode return
			if v == nil {
				return nil
			}

			outReader, ok := v.(io.Reader)
			if !ok {
				return e.New(e.TypeErr(outReader, v))
			}

			outPath := getOutPath(req)

			cmplvl, err := cmdenv.GetCompressLevel(getCompressOptions(req))
			if err != nil {
				return err
			}

			archive, _ := req.Options[archiveOptionName].(bool)
			gw := getWriter{
				Out:         os.Stdout,
				Err:         os.Stderr,
				Archive:     archive,
				Compression: cmplvl,
				Size:        int64(res.Length()),
			}

			return gw.Write(outReader, outPath)
		},
	},
}

type clearlineReader struct {
	io.Reader
	out io.Writer
}

func (r *clearlineReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	if err == io.EOF {
		// callback
		fmt.Fprintf(r.out, "\033[2K\r") // clear progress bar line on EOF
	}
	return
}

func progressBarForReader(out io.Writer, r io.Reader, l int64) (*pb.ProgressBar, io.Reader) {
	bar := makeProgressBar(out, l)
	barR := bar.NewProxyReader(r)
	return bar, &clearlineReader{barR, out}
}

func makeProgressBar(out io.Writer, l int64) *pb.ProgressBar {
	// setup bar reader
	// TODO: get total length of files
	bar := pb.New64(l).SetUnits(pb.U_BYTES)
	bar.Output = out

	// the progress bar lib doesn't give us a way to get the width of the output,
	// so as a hack we just use a callback to measure the output, then get rid of it
	bar.Callback = func(line string) {
		terminalWidth := len(line)
		bar.Callback = nil
		log.Infof("terminal width: %v\n", terminalWidth)
	}
	return bar
}

func getOutPath(req *cmds.Request) string {
	outPath, _ := req.Options[outputOptionName].(string)
	if outPath == "" {
		trimmed := strings.TrimRight(req.Arguments[0], "/")
		_, outPath = filepath.Split(trimmed)
		outPath = filepath.Clean(outPath)
	}
	return outPath
}

type getWriter struct {
	Out io.Writer // for output to user
	Err io.Writer // for progress bar output

	Archive     bool
	Compression int
	Size        int64
}

func (gw *getWriter) Write(r io.Reader, fpath string) error {
	if gw.Archive || gw.Compression != gzip.NoCompression {
		return gw.writeArchive(r, fpath)
	}
	return gw.writeExtracted(r, fpath)
}

func (gw *getWriter) writeArchive(r io.Reader, fpath string) error {
	// adjust file name if tar
	if gw.Archive {
		if !strings.HasSuffix(fpath, ".tar") && !strings.HasSuffix(fpath, ".tar.gz") {
			fpath += ".tar"
		}
	}

	// adjust file name if gz
	if gw.Compression != gzip.NoCompression {
		if !strings.HasSuffix(fpath, ".gz") {
			fpath += ".gz"
		}
	}

	// create file
	file, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(gw.Out, "Saving archive to %s\n", fpath)
	bar, barR := progressBarForReader(gw.Err, r, gw.Size)
	bar.Start()
	defer bar.Finish()

	_, err = io.Copy(file, barR)
	return err
}

func (gw *getWriter) writeExtracted(r io.Reader, fpath string) error {
	fmt.Fprintf(gw.Out, "Saving file(s) to %s\n", fpath)
	bar := makeProgressBar(gw.Err, gw.Size)
	bar.Start()
	defer bar.Finish()
	defer bar.Set64(gw.Size)

	extractor := &tar.Extractor{Path: fpath, Progress: bar.Add64}
	return extractor.Extract(r)
}

func getCompressOptions(req *cmds.Request) (bool, int) {
	cmprs, _ := req.Options[compressOptionName].(bool)
	cmplvl, _ := req.Options[compressionLevelOptionName].(int)
	return cmprs, cmplvl
}
