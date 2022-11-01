package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/tracker/udp"
	cmds "github.com/bittorrent/go-btfs-cmds"
	cid "github.com/ipfs/go-cid"
	mbase "github.com/multiformats/go-multibase"
)

var bittorrentCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "A tool command to integrate with bittorrent net(support bittorrent seed or a magnet URI scheme).",
	},
	Subcommands: map[string]*cmds.Command{
		"metainfo": metainfoBTCmd,
		"scrape":   scrapeBTCmd,
		"bencode":  bencodeBTCmd,
		"download": downloadBTCmd,
		"serve":    downloadBTCmd,
	},
	Extra: CreateCmdExtras(SetDoesNotUseRepo(true)),
}

var metainfoBTCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Print the metainfo of a bittorrent file from a bittorrent seed file.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("path", true, true, "The path to a bittorrent file in which you want to get metainfo."),
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
		// d["PieceHashes"] = func() (ret []string) {
		// 	for i := range iter.N(info.NumPieces()) {
		// 		ret = append(ret, hex.EncodeToString(info.Pieces[i*20:(i+1)*20]))
		// 	}
		// 	return
		// }()
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

var scrapeBTCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Fetch swarm metrics for info-hashes from tracker.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("tracker-url", true, false, "The tracker url."),
		cmds.StringArg("info-hash", true, true, "The path to a bittorrent file in which you want to get metainfo."),
	},
	Run: func(req *cmds.Request, resp cmds.ResponseEmitter, env cmds.Environment) error {
		tracker := req.Arguments[0]

		trackerUrl, err := url.Parse(tracker)
		if err != nil {
			return fmt.Errorf("parsing tracker url: %w", err)
		}
		cc, err := udp.NewConnClient(udp.NewConnClientOpts{
			Network: trackerUrl.Scheme,
			Host:    trackerUrl.Host,
		})
		if err != nil {
			return fmt.Errorf("creating new udp tracker conn client: %w", err)
		}
		defer cc.Close()
		var ihs []udp.InfoHash
		for _, hashStr := range req.Arguments[1:] {
			ih := metainfo.NewHashFromHex(hashStr)
			ihs = append(ihs, ih)
		}
		scrapeOut, err := cc.Client.Scrape(context.TODO(), ihs)
		if err != nil {
			return fmt.Errorf("scraping: %w", err)
		}
		return resp.Emit(scrapeOut)
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

var bencodeBTCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Print the bencoded info person-friendly of a bittorrent file from a bittorrent seed file.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("path", true, true, "The path to a bittorrent file in which bencoded data stored."),
	},
	Run: func(req *cmds.Request, resp cmds.ResponseEmitter, env cmds.Environment) error {
		btFilePath := req.Arguments[0]
		f, err := os.Open(btFilePath)
		if err != nil {
			return err
		}
		defer f.Close()
		d := bencode.NewDecoder(f)
		var v interface{}
		err = d.Decode(&v)
		if err != nil {
			return fmt.Errorf("decoding message : %w", err)
		}
		resp.Emit(v)
		return nil
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
