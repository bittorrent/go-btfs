package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/tracker/udp"
	cmds "github.com/bittorrent/go-btfs-cmds"
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
		cmds.StringArg("magnet uri", false, false, "Magnet uri if your seed is coming from magnet."),
	},
	Options: []cmds.Option{
		cmds.StringOption("t", "Bittorrent seed file."),
	},
	Run: func(req *cmds.Request, resp cmds.ResponseEmitter, env cmds.Environment) error {
		magnet := req.Arguments[0]
		btFilePath, _ := req.Options["t"].(string)
		clientConfig := torrent.NewDefaultClientConfig()

		_, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		client, err := torrent.NewClient(clientConfig)
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		defer client.Close()
		var t *torrent.Torrent
		if btFilePath != "" {
			metaInfo, err := metainfo.LoadFromFile(btFilePath)
			if err != nil {
				return fmt.Errorf("error loading torrent file %s: %w", btFilePath, err)
			}
			t, err = client.AddTorrent(metaInfo)
			if err != nil {
				return fmt.Errorf("adding torrent: %w", err)
			}
		} else if magnet != "" {
			t, err = client.AddMagnet(magnet)
			if err != nil {
				return fmt.Errorf("error adding magnet: %w", err)
			}
		} else {
			return fmt.Errorf("your must provide a magnet uri or a torrent file path")
		}
		go func() {
			client.WriteStatus(os.Stdout)
		}()
		<-t.GotInfo()
		t.DownloadAll()
		client.WaitAll()
		return nil
	},
}
