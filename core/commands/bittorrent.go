package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/anacrolix/torrent/metainfo"
	cmds "github.com/bittorrent/go-btfs-cmds"
	cid "github.com/ipfs/go-cid"
	mbase "github.com/multiformats/go-multibase"
)

var bittorrentCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Convert and discover properties of CIDs",
	},
	Subcommands: map[string]*cmds.Command{
		"metainfo": metainfoBTCmd,
		"scrape":   downloadBTCmd,
		"bencode":  downloadBTCmd,
		"download": downloadBTCmd,
		"serve":    downloadBTCmd,
	},
	Extra: CreateCmdExtras(SetDoesNotUseRepo(true)),
}

var metainfoBTCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Download a bittorrent file from the bittorrent seed or a magnet URL.",
		ShortDescription: "Download a bittorrent file from the bittorrent seed or a magnet URL.",
	},
	Arguments: []cmds.Argument{
		// cmds.FileArg("path", true, true, "The path to a file in which you want to get metainfo.").EnableRecursive().EnableStdin(),
		cmds.StringArg("path", true, true, "The path to a bittorrent file in which you want to get metainfo.").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.StringOption(cidFormatOptionName, "Printf style format string.").WithDefault("%s"),
		cmds.StringOption(cidVersionOptionName, "CID version to convert to."),
	},
	Run: func(req *cmds.Request, resp cmds.ResponseEmitter, env cmds.Environment) error {
		btFilePath := req.Arguments[0]
		mi, err := metainfo.LoadFromFile(btFilePath)
		if err != nil {
			return err
		}
		info, err := mi.UnmarshalInfo()
		if err != nil {
			return fmt.Errorf("error unmarshalling info: %s", err)
		}
		d := map[string]interface{}{
			"Name":         info.Name,
			"Name.Utf8":    info.NameUtf8,
			"NumPieces":    info.NumPieces(),
			"PieceLength":  info.PieceLength,
			"InfoHash":     mi.HashInfoBytes().HexString(),
			"NumFiles":     len(info.UpvertedFiles()),
			"TotalLength":  info.TotalLength(),
			"Announce":     mi.Announce,
			"AnnounceList": mi.AnnounceList,
			"UrlList":      mi.UrlList,
			"Files":        info.UpvertedFiles(),
		}
		if len(mi.Nodes) > 0 {
			d["Nodes"] = mi.Nodes
		}
		return resp.Emit(d)
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(func(req *cmds.Request, w io.Writer, out interface{}) error {
			marshaled, err := json.MarshalIndent(out, "", "	")
			if err != nil {
				return err
			}
			marshaled = append(marshaled, byte('\n'))
			fmt.Fprintln(w, string(marshaled))
			return nil
		}),
	},
}

var downloadBTCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:         "Download a bittorrent file from the bittorrent seed or a magnet URL.",
		LongDescription: "Download a bittorrent file from the bittorrent seed or a magnet URL.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("cid", true, true, "Cids to format.").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.StringOption(cidFormatOptionName, "Printf style format string.").WithDefault("%s"),
		cmds.StringOption(cidVersionOptionName, "CID version to convert to."),
		cmds.StringOption(cidCodecOptionName, "CID codec to convert to."),
		cmds.StringOption(cidMultibaseOptionName, "Multibase to display CID in."),
	},
	Run: func(req *cmds.Request, resp cmds.ResponseEmitter, env cmds.Environment) error {
		fmtStr, _ := req.Options[cidFormatOptionName].(string)
		verStr, _ := req.Options[cidVersionOptionName].(string)
		codecStr, _ := req.Options[cidCodecOptionName].(string)
		baseStr, _ := req.Options[cidMultibaseOptionName].(string)

		opts := cidFormatOpts{}

		if strings.IndexByte(fmtStr, '%') == -1 {
			return fmt.Errorf("invalid format string: %s", fmtStr)
		}
		opts.fmtStr = fmtStr

		if codecStr != "" {
			codec, ok := cid.Codecs[codecStr]
			if !ok {
				return fmt.Errorf("unknown IPLD codec: %s", codecStr)
			}
			opts.newCodec = codec
		} // otherwise, leave it as 0 (not a valid IPLD codec)

		switch verStr {
		case "":
			// noop
		case "0":
			if opts.newCodec != 0 && opts.newCodec != cid.DagProtobuf {
				return fmt.Errorf("cannot convert to CIDv0 with any codec other than DagPB")
			}
			opts.verConv = toCidV0
		case "1":
			opts.verConv = toCidV1
		default:
			return fmt.Errorf("invalid cid version: %s", verStr)
		}

		if baseStr != "" {
			encoder, err := mbase.EncoderByName(baseStr)
			if err != nil {
				return err
			}
			opts.newBase = encoder.Encoding()
		} else {
			opts.newBase = mbase.Encoding(-1)
		}

		return emitCids(req, resp, opts)
	},
	PostRun: cmds.PostRunMap{
		cmds.CLI: streamResult(func(v interface{}, out io.Writer) nonFatalError {
			r := v.(*CidFormatRes)
			if r.ErrorMsg != "" {
				return nonFatalError(fmt.Sprintf("%s: %s", r.CidStr, r.ErrorMsg))
			}
			fmt.Fprintf(out, "%s\n", r.Formatted)
			return ""
		}),
	},
	Type: CidFormatRes{},
}
