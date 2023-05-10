package commands

import (
	"errors"

	"github.com/bittorrent/go-btfs/core/commands/bttc"
	"github.com/bittorrent/go-btfs/core/commands/cheque"
	cmdenv "github.com/bittorrent/go-btfs/core/commands/cmdenv"
	dag "github.com/bittorrent/go-btfs/core/commands/dag"
	name "github.com/bittorrent/go-btfs/core/commands/name"
	ocmd "github.com/bittorrent/go-btfs/core/commands/object"
	settlement "github.com/bittorrent/go-btfs/core/commands/settlements"
	"github.com/bittorrent/go-btfs/core/commands/storage"
	"github.com/bittorrent/go-btfs/core/commands/storage/challenge"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/upload"
	unixfs "github.com/bittorrent/go-btfs/core/commands/unixfs"
	"github.com/bittorrent/go-btfs/core/commands/vault"

	cmds "github.com/bittorrent/go-btfs-cmds"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("core/commands")

var ErrNotOnline = errors.New("this command must be run in online mode. Try running 'btfs daemon' first")

const (
	ConfigOption  = "config"
	DebugOption   = "debug"
	LocalOption   = "local" // DEPRECATED: use OfflineOption
	OfflineOption = "offline"
	ApiOption     = "api"
)

var Root = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:  "Global p2p merkle-dag filesystem.",
		Synopsis: "btfs [--config=<config> | -c] [--debug | -D] [--help] [-h] [--api=<api>] [--offline] [--cid-base=<base>] [--upgrade-cidv0-in-output] [--encoding=<encoding> | --enc] [--timeout=<timeout>] <command> ...",
		Subcommands: `
BASIC COMMANDS
  init          Initialize btfs local configuration
  add <path>    Add a file to BTFS
  cat <ref>     Show BTFS object data
  get <ref>     Download BTFS objects
  ls <ref>      List links from an object
  refs <ref>    List hashes of links from an object

BTFS COMMANDS
  storage       Manage client and host storage features
  rm            Clean up locally stored files and objects

DATA STRUCTURE COMMANDS
  block         Interact with raw blocks in the datastore
  object        Interact with raw dag nodes
  files         Interact with objects as if they were a unix filesystem
  dag           Interact with IPLD documents (experimental)
  metadata      Interact with metadata for BTFS files

ASSETS COMMANDS
  bttc          BTTC service related commands, e.g. swap between BTT and WBTT
  vault         Vault service reladted commands, e.g. deposit WBTT to your vault
  cheque        Cheque service related commands, e.g. list cheques you have received
  settlement    Show cheque settlement info

ADVANCED COMMANDS
  daemon        Start a long-running daemon process
  mount         Mount an BTFS read-only mount point
  resolve       Resolve any type of name
  name          Publish and resolve BTNS names
  key           Create and list BTNS name keypairs
  dns           Resolve DNS links
  pin           Pin objects to local storage
  repo          Manipulate the BTFS repository
  stats         Various operational stats
  p2p           Libp2p stream mounting
  filestore     Manage the filestore (experimental)

NETWORK COMMANDS
  id            Show info about BTFS peers
  bootstrap     Add or remove bootstrap peers
  swarm         Manage connections to the p2p network
  dht           Query the DHT for values or peers
  ping          Measure the latency of a connection
  diag          Print diagnostics

TOOL COMMANDS
  config        Manage configuration
  version       Show btfs version information
  commands      List all available commands
  cid           Convert and discover properties of CIDs
  log           Manage and show logs of running daemon

Use 'btfs <command> --help' to learn more about each command.

btfs uses a repository in the local file system. By default, the repo is
located at ~/.btfs. To change the repo location, set the $BTFS_PATH
environment variable:

  export BTFS_PATH=/path/to/btfsrepo

EXIT STATUS

The CLI will exit with one of the following values:

0     Successful execution.
1     Failed executions.
`,
	},
	Options: []cmds.Option{
		cmds.StringOption(ConfigOption, "c", "Path to the configuration file to use."),
		cmds.BoolOption(DebugOption, "D", "Operate in debug mode."),
		cmds.BoolOption(cmds.OptLongHelp, "Show the full command help text."),
		cmds.BoolOption(cmds.OptShortHelp, "Show a short version of the command help text."),
		cmds.BoolOption(LocalOption, "L", "Run the command locally, instead of using the daemon. DEPRECATED: use --offline."),
		cmds.BoolOption(OfflineOption, "Run the command offline."),
		cmds.StringOption(ApiOption, "Use a specific API instance (defaults to /ip4/127.0.0.1/tcp/5001)"),

		// global options, added to every command
		cmdenv.OptionCidBase,
		cmdenv.OptionUpgradeCidV0InOutput,

		cmds.OptionEncodingType,
		cmds.OptionStreamChannels,
		cmds.OptionTimeout,
	},
}

// commandsDaemonCmd is the "btfs commands" command for daemon
var CommandsDaemonCmd = CommandsCmd(Root)

var rootSubcommands = map[string]*cmds.Command{
	"add":        AddCmd,
	"bitswap":    BitswapCmd,
	"block":      BlockCmd,
	"cat":        CatCmd,
	"commands":   CommandsDaemonCmd,
	"files":      FilesCmd,
	"filestore":  FileStoreCmd,
	"get":        GetCmd,
	"pubsub":     PubsubCmd,
	"repo":       RepoCmd,
	"stats":      StatsCmd,
	"bootstrap":  BootstrapCmd,
	"test":       TestCmd,
	"config":     ConfigCmd,
	"dag":        dag.DagCmd,
	"dht":        DhtCmd,
	"diag":       DiagCmd,
	"dns":        DNSCmd,
	"id":         IDCmd,
	"key":        KeyCmd,
	"log":        LogCmd,
	"ls":         LsCmd,
	"mount":      MountCmd,
	"name":       name.NameCmd,
	"object":     ocmd.ObjectCmd,
	"pin":        PinCmd,
	"ping":       PingCmd,
	"p2p":        P2PCmd,
	"refs":       RefsCmd,
	"resolve":    ResolveCmd,
	"swarm":      SwarmCmd,
	"tar":        TarCmd,
	"file":       unixfs.UnixFSCmd,
	"urlstore":   urlStoreCmd,
	"version":    VersionCmd,
	"shutdown":   daemonShutdownCmd,
	"restart":    restartCmd,
	"cid":        CidCmd,
	"rm":         RmCmd,
	"storage":    storage.StorageCmd,
	"metadata":   MetadataCmd,
	"guard":      GuardCmd,
	"cheque":     cheque.ChequeCmd,
	"vault":      vault.VaultCmd,
	"bttc":       bttc.BttcCmd,
	"settlement": settlement.SettlementCmd,
	//"update":    ExternalBinary(),
	"network":        NetworkCmd,
	"statuscontract": StatusContractCmd,
	"bittorrent":     bittorrentCmd,
}

// RootRO is the readonly version of Root
var RootRO = &cmds.Command{}

var CommandsDaemonROCmd = CommandsCmd(RootRO)

// RefsROCmd is `btfs refs` command
var RefsROCmd = &cmds.Command{}

// VersionROCmd is `btfs version` command (without deps).
var VersionROCmd = &cmds.Command{}

var rootROSubcommands = map[string]*cmds.Command{
	"commands": CommandsDaemonROCmd,
	"cat":      CatCmd,
	"block": &cmds.Command{
		Subcommands: map[string]*cmds.Command{
			"stat": blockStatCmd,
			"get":  blockGetCmd,
		},
	},
	"get": GetCmd,
	"dns": DNSCmd,
	"ls":  LsCmd,
	"name": {
		Subcommands: map[string]*cmds.Command{
			"resolve": name.IpnsCmd,
		},
	},
	"object": {
		Subcommands: map[string]*cmds.Command{
			"data":  ocmd.ObjectDataCmd,
			"links": ocmd.ObjectLinksCmd,
			"get":   ocmd.ObjectGetCmd,
			"stat":  ocmd.ObjectStatCmd,
		},
	},
	"dag": {
		Subcommands: map[string]*cmds.Command{
			"get":     dag.DagGetCmd,
			"resolve": dag.DagResolveCmd,
			"stat":    dag.DagStatCmd,
		},
	},
	"resolve": ResolveCmd,
}

// RootRemote is the remote-facing version of Root
var RootRemote = &cmds.Command{}

var rootRemoteSubcommands = map[string]*cmds.Command{
	"storage": &cmds.Command{
		Subcommands: map[string]*cmds.Command{
			"challenge": &cmds.Command{
				Subcommands: map[string]*cmds.Command{
					"response": challenge.StorageChallengeResponseCmd,
				},
			},
			"upload": &cmds.Command{
				Subcommands: map[string]*cmds.Command{
					"init":          upload.StorageUploadInitCmd,
					"supporttokens": upload.StorageUploadSupportTokensCmd,
					"recvcontract":  upload.StorageUploadRecvContractCmd,
					"cheque":        upload.StorageUploadChequeCmd,
				},
			},
			"dcrepair": &cmds.Command{
				Subcommands: map[string]*cmds.Command{
					"response": upload.HostRepairResponseCmd,
				},
			},
		},
	},
	"p2p": &cmds.Command{
		Subcommands: map[string]*cmds.Command{
			"handshake": P2phandshakeCmd,
		},
	},
}

func init() {
	Root.ProcessHelp()
	*RootRO = *Root
	*RootRemote = *Root

	// this was in the big map definition above before,
	// but if we leave it there lgc.NewCommand will be executed
	// before the value is updated (:/sanitize readonly refs command/)

	// sanitize readonly refs command
	*RefsROCmd = *RefsCmd
	RefsROCmd.Subcommands = map[string]*cmds.Command{}
	rootROSubcommands["refs"] = RefsROCmd

	// sanitize readonly version command (no need to expose precise deps)
	*VersionROCmd = *VersionCmd
	VersionROCmd.Subcommands = map[string]*cmds.Command{}
	rootROSubcommands["version"] = VersionROCmd
	// also sanitize remote version command
	rootRemoteSubcommands["version"] = VersionROCmd

	Root.Subcommands = rootSubcommands
	RootRO.Subcommands = rootROSubcommands
	RootRemote.Subcommands = rootRemoteSubcommands
}

type MessageOutput struct {
	Message string
}
