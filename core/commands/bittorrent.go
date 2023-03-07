package commands

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"github.com/anacrolix/torrent/tracker/udp"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bradfitz/iter"
	humanize "github.com/dustin/go-humanize"
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
		"serve":    serveBTCmd,
	},
	Extra: CreateCmdExtras(SetDoesNotUseRepo(true)),
}

var metainfoBTCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Print the metainfo of a bittorrent file from a bittorrent seed file.",
	},
	Options: []cmds.Option{
		cmds.StringOption("t", "Bittorrent seed file."),
		cmds.StringOption("m", "Magnet uri."),
	},
	Run: func(req *cmds.Request, resp cmds.ResponseEmitter, env cmds.Environment) error {
		btFilePath, _ := req.Options["t"].(string)
		magnetUri, _ := req.Options["m"].(string)
		var mi *metainfo.MetaInfo
		var err error
		if btFilePath != "" {
			mi, err = metainfo.LoadFromFile(btFilePath)
			if err != nil {
				return fmt.Errorf("LoadFromFile bittorrent seed file: %w", err)
			}
		} else if magnetUri != "" {
			clientConfig := torrent.NewDefaultClientConfig()
			clientConfig.ListenPort = 0
			client, err := torrent.NewClient(clientConfig)
			if err != nil {
				return fmt.Errorf("creating client: %w", err)
			}
			defer client.Close()
			t, err := client.AddMagnet(magnetUri)
			if err != nil {
				return fmt.Errorf("client.AddMagnet: %w", err)
			}
			select {
			case <-t.GotInfo():
			case <-time.After(5 * time.Minute):
				return fmt.Errorf("get magnet from bt network timeout, may be this seed cannot find")
			}
			m := t.Metainfo()
			mi = &m
		} else {
			return fmt.Errorf("you must specify a -t to represent the path of a bt seed file or -m to represent a magnet uri")
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
		d["PieceHashes"] = func() (ret []string) {
			for i := range iter.N(info.NumPieces()) {
				ret = append(ret, hex.EncodeToString(info.Pieces[i*20:(i+1)*20]))
			}
			return
		}()
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
	Options: []cmds.Option{
		cmds.StringOption("t", "The tracker url."),
		cmds.StringOption("i", "the hash list of a file or piece which is separated by commas"),
	},
	Run: func(req *cmds.Request, resp cmds.ResponseEmitter, env cmds.Environment) error {
		tracker, _ := req.Options["t"].(string)
		if tracker == "" {
			return fmt.Errorf("tracker(-t) must be specified")
		}
		hashListStr, _ := req.Options["i"].(string)
		if hashListStr == "" {
			return fmt.Errorf("hash list(-i) must be specified and is separated by commas")
		}
		hashList := strings.Split(hashListStr, ",")
		var ihs []udp.InfoHash
		for _, hashStr := range hashList {
			if len(hashStr) != 2*metainfo.HashSize {
				return fmt.Errorf("hash info must be %d characters", 2*metainfo.HashSize)
			}
			ih := metainfo.NewHashFromHex(hashStr)
			ihs = append(ihs, ih)
		}
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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		scrapeOut, err := cc.Client.Scrape(ctx, ihs)
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
		Tagline: "Download a bittorrent file from the bittorrent seed or a magnet URL.",
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
		clientConfig.ListenPort = 0
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
		select {
		case <-t.GotInfo():
			fmt.Println("Got metainfo done.Begin to download files...")
		case <-time.After(5 * time.Minute):
			log.Error("Get metainfo timeout,exceed two minutes, we can't find the metainfo for this torrent.")
			return fmt.Errorf("get metainfo timeout")
		}
		t.DownloadAll()
		// print the progress of the download.
		fmt.Printf("This torrent needs storage space about: %s\n", humanize.Bytes(uint64(t.Length())))
		torrentBar(t, false)
		isCompleted := client.WaitAll()
		if !isCompleted {
			log.Error("download error because of the closed of the client")
			return fmt.Errorf("download error because of the closed of the client")
		}

		btfsBinaryPath := "btfs"
		cmd := exec.Command(btfsBinaryPath, "add", "-r", "--pin", t.Name())

		go func() {
			time.Sleep(10 * time.Minute)
			_, err := os.FindProcess(int(cmd.Process.Pid))
			if err != nil {
				log.Info("process already finished\n")
			} else {
				err := cmd.Process.Kill()
				if err != nil {
					if !strings.Contains(err.Error(), "process already finished") {
						log.Errorf("cannot kill process: [%v] \n", err)
					}
				}
			}
		}()
		// No matter what, delete the files after the download completed.
		defer func() {
			err := os.RemoveAll(t.Name())
			if err != nil {
				log.Errorf("cannot remove the downloaded files: [%s], [%v] \n", t.Name(), err)
			}
		}()
		var errbuf bytes.Buffer
		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("btfs add [%s] failed: [%v], [%s]", t.Name(), err, errbuf.String())
		}
		s := strings.Split(string(out), " ")
		if len(s) < 2 {
			return fmt.Errorf("btfs add test failed: invalid add result[%s]", string(out))
		}
		fmt.Println(string(out))
		return nil
	},
}

var minBTListenPort = 30000
var maxBTListenPort = 31000

var serveBTCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:         "Serve as a bittorrent client with the specified files.",
		LongDescription: "Serve as a bittorrent client with the specified files.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("path", true, true, "the paths of some files that you want to serve as seeds"),
	},
	Run: func(req *cmds.Request, resp cmds.ResponseEmitter, env cmds.Environment) error {
		filePaths := req.Arguments
		if len(filePaths) == 0 {
			return fmt.Errorf("you must provide the paths of some files that you want to serve as seeds")
		}
		cfg := torrent.NewDefaultClientConfig()
		cfg.ListenPort = minBTListenPort
		cfg.Seed = true
	retry:
		cl, err := torrent.NewClient(cfg)
		if err != nil {
			if strings.Contains(err.Error(), "address already in use") {
				fmt.Println(err)
				cfg.ListenPort = cfg.ListenPort + 1
				if cfg.ListenPort > maxBTListenPort {
					return fmt.Errorf("we have try all the port between %d and %d ,but they are all in used", minBTListenPort, maxBTListenPort)
				}
				goto retry
			} else {
				return fmt.Errorf("new torrent client: %w", err)
			}
		}
		defer cl.Close()

		for _, filePath := range filePaths {
			totalLength, err := totalLength(filePath)
			if err != nil {
				return fmt.Errorf("calculating total length of %q: %v", filePath, err)
			}
			pieceLength := metainfo.ChoosePieceLength(totalLength)
			info := metainfo.Info{
				PieceLength: pieceLength,
			}
			err = info.BuildFromFilePath(filePath)
			if err != nil {
				return fmt.Errorf("building info from path %q: %w", filePath, err)
			}
			for _, fi := range info.Files {
				fmt.Printf("added %q\n", fi.Path)
			}
			mi := metainfo.MetaInfo{
				InfoBytes: bencode.MustMarshal(info),
			}
			pc, err := storage.NewDefaultPieceCompletionForDir(".")
			if err != nil {
				return fmt.Errorf("new piece completion: %w", err)
			}
			defer pc.Close()
			ih := mi.HashInfoBytes()
			to, _ := cl.AddTorrentOpt(torrent.AddTorrentOpts{
				InfoHash: ih,
				Storage: storage.NewFileOpts(storage.NewFileClientOpts{
					ClientBaseDir: filePath,
					FilePathMaker: func(opts storage.FilePathMakerOpts) string {
						return filepath.Join(opts.File.Path...)
					},
					TorrentDirMaker: nil,
					PieceCompletion: pc,
				}),
			})
			defer to.Drop()
			err = to.MergeSpec(&torrent.TorrentSpec{
				InfoBytes: mi.InfoBytes,
				Trackers: [][]string{{
					`wss://tracker.btorrent.xyz`,
					`wss://tracker.openwebtorrent.com`,
					"http://p4p.arenabg.com:1337/announce",
					"udp://tracker.opentrackr.org:1337/announce",
					"udp://tracker.openbittorrent.com:6969/announce",
				}},
			})
			if err != nil {
				return fmt.Errorf("setting trackers: %w", err)
			}
			fmt.Printf("%v: %v\n", to, to.Metainfo().Magnet(&ih, &info))
		}
		select {}
	},
}

func torrentBar(t *torrent.Torrent, pieceStates bool) {
	go func() {
		start := time.Now()
		if t.Info() == nil {
			fmt.Printf("%v: getting torrent info for %q\n", time.Since(start), t.Name())
			<-t.GotInfo()
		}
		lastStats := t.Stats()
		var lastLine string
		interval := 10 * time.Second
		tick := time.NewTicker(interval)
		for range tick.C {
			var completedPieces, partialPieces int
			psrs := t.PieceStateRuns()
			for _, r := range psrs {
				if r.Complete {
					completedPieces += r.Length
				}
				if r.Partial {
					partialPieces += r.Length
				}
			}
			stats := t.Stats()
			byteRate := int64(time.Second)
			byteRate *= stats.BytesReadUsefulData.Int64() - lastStats.BytesReadUsefulData.Int64()
			byteRate /= int64(interval)
			line := fmt.Sprintf(
				"%v: downloading %q: %s/%s, %d/%d pieces completed (%d partial): %v/s\n",
				time.Since(start),
				t.Name(),
				humanize.Bytes(uint64(t.BytesCompleted())),
				humanize.Bytes(uint64(t.Length())),
				completedPieces,
				t.NumPieces(),
				partialPieces,
				humanize.Bytes(uint64(byteRate)),
			)
			if line != lastLine {
				lastLine = line
				fmt.Println(line)
			}
			if pieceStates {
				fmt.Println(psrs)
			}
			lastStats = stats
			if t.Complete.Bool() {
				fmt.Println("Download completed!!!")
				tick.Stop()
				return
			}
		}
	}()
}

func totalLength(path string) (totalLength int64, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		totalLength += info.Size()
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("walking path, %w", err)
	}
	return totalLength, nil
}
