package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	_ "expvar"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bittorrent/go-btfs/chain/tokencfg"

	"github.com/bittorrent/go-btfs/guide"

	version "github.com/bittorrent/go-btfs"
	cmds "github.com/bittorrent/go-btfs-cmds"
	config "github.com/bittorrent/go-btfs-config"
	cserial "github.com/bittorrent/go-btfs-config/serialize"
	"github.com/bittorrent/go-btfs/bindata"
	"github.com/bittorrent/go-btfs/chain"
	cc "github.com/bittorrent/go-btfs/chain/config"
	chainconfig "github.com/bittorrent/go-btfs/chain/config"
	utilmain "github.com/bittorrent/go-btfs/cmd/btfs/util"
	oldcmds "github.com/bittorrent/go-btfs/commands"
	"github.com/bittorrent/go-btfs/core"
	commands "github.com/bittorrent/go-btfs/core/commands"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/path"
	corehttp "github.com/bittorrent/go-btfs/core/corehttp"
	httpremote "github.com/bittorrent/go-btfs/core/corehttp/remote"
	corerepo "github.com/bittorrent/go-btfs/core/corerepo"
	libp2p "github.com/bittorrent/go-btfs/core/node/libp2p"
	nodeMount "github.com/bittorrent/go-btfs/fuse/node"
	"github.com/bittorrent/go-btfs/repo"
	fsrepo "github.com/bittorrent/go-btfs/repo/fsrepo"
	"github.com/bittorrent/go-btfs/reportstatus"
	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	"github.com/bittorrent/go-btfs/spin"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/bittorrent/go-btfs/transaction/crypto"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/ethereum/go-ethereum/common"

	cp "github.com/bittorrent/go-btfs-common/crypto"
	nodepb "github.com/bittorrent/go-btfs-common/protos/node"
	multierror "github.com/hashicorp/go-multierror"
	util "github.com/ipfs/go-ipfs-util"
	mprome "github.com/ipfs/go-metrics-prometheus"
	goprocess "github.com/jbenet/goprocess"
	sockets "github.com/libp2p/go-socket-activation"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	prometheus "github.com/prometheus/client_golang/prometheus"
	promauto "github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	adjustFDLimitKwd          = "manage-fdlimit"
	enableGCKwd               = "enable-gc"
	initOptionKwd             = "init"
	initConfigOptionKwd       = "init-config"
	initProfileOptionKwd      = "init-profile"
	ipfsMountKwd              = "mount-btfs"
	ipnsMountKwd              = "mount-btns"
	migrateKwd                = "migrate"
	mountKwd                  = "mount"
	offlineKwd                = "offline" // global option
	routingOptionKwd          = "routing"
	routingOptionSupernodeKwd = "supernode"
	routingOptionDHTClientKwd = "dhtclient"
	routingOptionDHTKwd       = "dht"
	routingOptionDHTServerKwd = "dhtserver"
	routingOptionNoneKwd      = "none"
	routingOptionCustomKwd    = "custom"
	routingOptionDefaultKwd   = "default"
	unencryptTransportKwd     = "disable-transport-encryption"
	unrestrictedApiAccessKwd  = "unrestricted-api"
	writableKwd               = "writable"
	enablePubSubKwd           = "enable-pubsub-experiment"
	enableIPNSPubSubKwd       = "enable-namesys-pubsub"
	enableMultiplexKwd        = "enable-mplex-experiment"
	hValueKwd                 = "hval"
	enableDataCollection      = "dc"
	enableStartupTest         = "enable-startup-test"
	swarmPortKwd              = "swarm-port"
	deploymentGasPrice        = "deployment-gasPrice"
	chainID                   = "chain-id"
	// apiAddrKwd    = "address-api"
	// swarmAddrKwd  = "address-swarm"
)

// BTFS daemon test exit error code
const (
	findBTFSBinaryFailed = 100
	getFileTestFailed    = 101
	addFileTestFailed    = 102
)

var daemonCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Run a network-connected BTFS node.",
		ShortDescription: `
'btfs daemon' runs a persistent btfs daemon that can serve commands
over the network. Most applications that use BTFS will do so by
communicating with a daemon over the HTTP API. While the daemon is
running, calls to 'btfs' commands will be sent over the network to
the daemon.
`,
		LongDescription: `
The daemon will start listening on ports on the network, which are
documented in (and can be modified through) 'btfs config Addresses'.
For example, to change the 'Gateway' port:

  btfs config Addresses.Gateway /ip4/127.0.0.1/tcp/8082

The API address can be changed the same way:

  btfs config Addresses.API /ip4/127.0.0.1/tcp/5002

Make sure to restart the daemon after changing addresses.

By default, the gateway is only accessible locally. To expose it to
other computers in the network, use 0.0.0.0 as the ip address:

  btfs config Addresses.Gateway /ip4/0.0.0.0/tcp/8080

Be careful if you expose the API. It is a security risk, as anyone could
control your node remotely. If you need to control the node remotely,
make sure to protect the port as you would other services or database
(firewall, authenticated proxy, etc).

HTTP Headers

btfs supports passing arbitrary headers to the API and Gateway. You can
do this by setting headers on the API.HTTPHeaders and Gateway.HTTPHeaders
keys:

  btfs config --json API.HTTPHeaders.X-Special-Header '["so special :)"]'
  btfs config --json Gateway.HTTPHeaders.X-Special-Header '["so special :)"]'

Note that the value of the keys is an _array_ of strings. This is because
headers can have more than one value, and it is convenient to pass through
to other libraries.

CORS Headers (for API)

You can setup CORS headers the same way:

  btfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["example.com"]'
  btfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "GET", "POST"]'
  btfs config --json API.HTTPHeaders.Access-Control-Allow-Credentials '["true"]'

Shutdown

To shut down the daemon, send a SIGINT signal to it (e.g. by pressing 'Ctrl-C')
or send a SIGTERM signal to it (e.g. with 'kill'). It may take a while for the
daemon to shutdown gracefully, but it can be killed forcibly by sending a
second signal.

BTFS_PATH environment variable

btfs uses a repository in the local file system. By default, the repo is
located at ~/.btfs. To change the repo location, set the $BTFS_PATH
environment variable:

  export BTFS_PATH=/path/to/btfsrepo

Routing

BTFS by default will use a DHT for content routing. There is a highly
experimental alternative that operates the DHT in a 'client only' mode that
can be enabled by running the daemon as:

  btfs daemon --routing=dhtclient

This will later be transitioned into a config option once it gets out of the
'experimental' stage.

DEPRECATION NOTICE

Previously, btfs used an environment variable as seen below:

  export API_ORIGIN="http://localhost:8888/"

This is deprecated. It is still honored in this version, but will be removed
in a future version, along with this notice. Please move to setting the HTTP
Headers.
`,
	},

	Options: []cmds.Option{
		cmds.BoolOption(initOptionKwd, "Initialize btfs with default settings if not already initialized"),
		cmds.StringOption(initConfigOptionKwd, "Path to existing configuration file to be loaded during --init"),
		cmds.StringOption(initProfileOptionKwd, "Configuration profiles to apply for --init. See btfs init --help for more"),
		cmds.StringOption(routingOptionKwd, "Overrides the routing option").WithDefault(routingOptionDefaultKwd),
		cmds.BoolOption(mountKwd, "Mounts BTFS to the filesystem"),
		cmds.BoolOption(writableKwd, "Enable writing objects (with POST, PUT and DELETE)"),
		cmds.StringOption(ipfsMountKwd, "Path to the mountpoint for BTFS (if using --mount). Defaults to config setting."),
		cmds.StringOption(ipnsMountKwd, "Path to the mountpoint for BTNS (if using --mount). Defaults to config setting."),
		cmds.BoolOption(unrestrictedApiAccessKwd, "Allow API access to unlisted hashes"),
		cmds.BoolOption(unencryptTransportKwd, "Disable transport encryption (for debugging protocols)"),
		cmds.BoolOption(enableGCKwd, "Enable automatic periodic repo garbage collection"),
		cmds.BoolOption(adjustFDLimitKwd, "Check and raise file descriptor limits if needed").WithDefault(true),
		cmds.BoolOption(migrateKwd, "If true, assume yes at the migrate prompt. If false, assume no.").WithDefault(true),
		cmds.BoolOption(enablePubSubKwd, "Instantiate the btfs daemon with the experimental pubsub feature enabled."),
		cmds.BoolOption(enableIPNSPubSubKwd, "Enable BTNS record distribution through pubsub; enables pubsub."),
		cmds.BoolOption(enableMultiplexKwd, "DEPRECATED"),
		cmds.StringOption(hValueKwd, "H-value identifies the BitTorrent client this daemon is started by. None if not started by a BitTorrent client."),
		cmds.BoolOption(enableDataCollection, "Allow BTFS to collect and send out node statistics."),
		cmds.BoolOption(enableStartupTest, "Allow BTFS to perform start up test.").WithDefault(false),
		cmds.StringOption(swarmPortKwd, "Override existing announced swarm address with external port in the format of [WAN:LAN]."),
		cmds.StringOption(deploymentGasPrice, "gas price in unit to use for deployment and funding."),
		cmds.StringOption(chainID, "The ID of blockchain to deploy."),
		// TODO: add way to override addresses. tricky part: updating the config if also --init.
		// cmds.StringOption(apiAddrKwd, "Address for the daemon rpc API (overrides config)"),
		// cmds.StringOption(swarmAddrKwd, "Address for the swarm socket (overrides config)"),
	},
	Subcommands: map[string]*cmds.Command{},
	NoRemote:    true,
	Extra:       commands.CreateCmdExtras(commands.SetDoesNotUseConfigAsInput(true)),
	Run:         wrapDaemonFunc,
}

// defaultMux tells mux to serve path using the default muxer. This is
// mostly useful to hook up things that register in the default muxer,
// and don't provide a convenient http.Handler entry point, such as
// expvar and http/pprof.
func defaultMux(path string) corehttp.ServeOption {
	return func(node *core.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		mux.Handle(path, http.DefaultServeMux)
		return mux, nil
	}
}

func wrapDaemonFunc(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) (_err error) {
	_err = daemonFunc(req, re, env)
	commands.NotifyAndWaitIfOnRestarting()
	return
}

func daemonFunc(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) (_err error) {
	cctx := env.(*oldcmds.Context)
	_, b := os.LookupEnv(path.BtfsPathKey)
	if !b {
		c := cctx.ConfigRoot
		if bs, err := ioutil.ReadFile(path.PropertiesFileName); err == nil && len(bs) > 0 {
			c = string(bs)
		}
		cctx.ConfigRoot = c
	}

	// Inject metrics before we do anything
	err := mprome.Inject()
	if err != nil {
		log.Errorf("Injecting prometheus handler for metrics failed with message: %s\n", err.Error())
	}

	// let the user know we're going.
	fmt.Printf("Initializing daemon...\n")

	defer func() {
		if _err != nil {
			// Print an extra line before any errors. This could go
			// in the commands lib but doesn't really make sense for
			// all commands.
			fmt.Println()
		}
	}()

	// print the btfs version
	printVersion()

	managefd, _ := req.Options[adjustFDLimitKwd].(bool)
	if managefd {
		if _, _, err := utilmain.ManageFdLimit(); err != nil {
			log.Errorf("setting file descriptor limit: %s", err)
		}
	}

	// check transport encryption flag.
	unencrypted, _ := req.Options[unencryptTransportKwd].(bool)
	if unencrypted {
		log.Warnf(`Running with --%s: All connections are UNENCRYPTED.
		You will not be able to connect to regular encrypted networks.`, unencryptTransportKwd)
	}

	// first, whether user has provided the initialization flag. we may be
	// running in an uninitialized state.
	initialize, _ := req.Options[initOptionKwd].(bool)
	hValue, hasHval := req.Options[hValueKwd].(string)
	inited := false
	if initialize {
		cfg := cctx.ConfigRoot
		if !fsrepo.IsInitialized(cfg) {
			cfgLocation, _ := req.Options[initConfigOptionKwd].(string)
			profiles, _ := req.Options[initProfileOptionKwd].(string)
			var conf *config.Config

			if cfgLocation != "" {
				if conf, err = cserial.Load(cfgLocation); err != nil {
					return err
				}
			}
			if hasHval && profiles == "" {
				profiles = "storage-host"
			}

			if err = doInit(os.Stdout, cfg, false, utilmain.NBitsForKeypairDefault, profiles, conf,
				keyTypeDefault, "", "", false, false); err != nil {
				return err
			}

			inited = true
		}
	}

	// acquire the repo lock _before_ constructing a node. we need to make
	// sure we are permitted to access the resources (datastore, etc.)
	repo, err := fsrepo.Open(cctx.ConfigRoot)
	switch err {
	default:
		if strings.Contains(err.Error(), "someone else has the lock") {
			fmt.Println(`Error:Someone else has the lock;
What causes this error: there is already one daemon process running in background
Solution: kill it first and run btfs daemon again.
If the user need to start multiple nodes on the same machine, the configuration needs to be modified to a new place.`)
		}
		return err
	case nil:
		break
	}

	// The node will also close the repo but there are many places we could
	// fail before we get to that. It can't hurt to close it twice.
	defer repo.Close()

	cfg, err := cctx.GetConfig()
	if err != nil {
		return err
	}

	if !inited {
		migrated := config.MigrateConfig(cfg, false, hasHval)
		if cfg.ChainInfo.PriceOracleAddress != "" {
			cfg.ChainInfo.PriceOracleAddress = ""
			migrated = true
		}
		if migrated {
			// Flush changes if migrated
			err = repo.SetConfig(cfg)
			if err != nil {
				return err
			}
		}
	}

	// Print self information for logging and debugging purposes
	fmt.Printf("Repo location: %s\n", cctx.ConfigRoot)
	fmt.Printf("Peer identity: %s\n", cfg.Identity.PeerID)

	privKey, err := cp.ToPrivKey(cfg.Identity.PrivKey)
	if err != nil {
		return err
	}

	keys, err := cp.FromIcPrivateKey(privKey)
	if err != nil {
		return err
	}

	// decode from string
	pkbytesOri, err := base64.StdEncoding.DecodeString(cfg.Identity.PrivKey)
	if err != nil {
		return err
	}
	//new singer
	pk := crypto.Secp256k1PrivateKeyFromBytes(pkbytesOri[4:])
	singer := crypto.NewDefaultSigner(pk)

	address0x, _ := singer.EthereumAddress()

	fmt.Println("the address of Bttc format is: ", address0x)
	fmt.Println("the address of Tron format is: ", keys.Base58Address)

	SimpleMode := cfg.SimpleMode
	if SimpleMode == false {
		// guide server init
		optionApiAddr, _ := req.Options[commands.ApiOption].(string)
		guide.SetServerAddr(cfg.Addresses.API, optionApiAddr)
		guide.SetInfo(&guide.Info{
			BtfsVersion: version.CurrentVersionNumber,
			HostID:      cfg.Identity.PeerID,
			BttcAddress: address0x.String(),
			PrivateKey:  hex.EncodeToString(pkbytesOri[4:]),
		})
		guide.StartServer()
		defer guide.TryShutdownServer()
	}

	//chain init
	configRoot := cctx.ConfigRoot
	statestore, err := chain.InitStateStore(configRoot)
	if err != nil {
		fmt.Println("init statestore err: ", err)
		return err
	}
	defer func() {
		statestore.Close()
	}()

	if SimpleMode == false {
		chainid, stored, err := getChainID(req, cfg, statestore)
		if err != nil {
			return err
		}
		chainCfg, err := chainconfig.InitChainConfig(cfg, stored, chainid)
		if err != nil {
			return err
		}

		// upgrade factory to v2 if necessary
		needUpdateFactory := false
		needUpdateFactory, err = doIfNeedUpgradeFactoryToV2(chainid, chainCfg, statestore, repo, cfg, configRoot)
		if err != nil {
			fmt.Printf("upgrade vault contract failed, err=%s\n", err)
			return err
		}
		if needUpdateFactory { // no error means upgrade preparation done, re-init the statestore
			statestore, err = chain.InitStateStore(configRoot)
			if err != nil {
				fmt.Println("init statestore err: ", err)
				return err
			}
			err = chain.StoreChainIdIfNotExists(chainid, statestore)
			if err != nil {
				fmt.Printf("save chainid failed, err: %s\n", err)
				return
			}
		}

		tokencfg.InitToken(chainid)

		//endpoint
		chainInfo, err := chain.InitChain(context.Background(), statestore, singer, time.Duration(1000000000),
			chainid, cfg.Identity.PeerID, chainCfg)
		if err != nil {
			return err
		}

		// Sync the with the given Ethereum backend:
		isSynced, _, err := transaction.IsSynced(context.Background(), chainInfo.Backend, chain.MaxDelay)
		if err != nil {
			return fmt.Errorf("is synced: %w", err)
		}

		if !isSynced {
			log.Infof("waiting to sync with the Ethereum backend")

			err := transaction.WaitSynced(context.Background(), chainInfo.Backend, chain.MaxDelay)
			if err != nil {
				return fmt.Errorf("waiting backend sync: %w", err)
			}
		}

		deployGasPrice, found := req.Options[deploymentGasPrice].(string)
		if !found {
			deployGasPrice = chainInfo.Chainconfig.DeploymentGas
		}

		/*settleinfo*/
		settleInfo, err := chain.InitSettlement(context.Background(), statestore, chainInfo, deployGasPrice, chainInfo.ChainID)
		if err != nil {
			fmt.Println("init settlement err: ", err)
			if strings.Contains(err.Error(), "insufficient funds") {
				fmt.Println("Please recharge BTT to your address to solve this error")
			}
			if strings.Contains(err.Error(), "contract deployment failed") {
				fmt.Println(`Solution1: It is recommended to check if the balance is sufficient. If the balance is low, it is recommended to top up.`)
				fmt.Println(`Solution2: Suggest to redeploy.`)
			}

			return err
		}

		/*upgrade vault implementation*/
		oldImpl, newImpl, err := settleInfo.VaultService.UpgradeTo(context.Background(), chainInfo.Chainconfig.VaultLogicAddress)
		if err != nil {
			emsg := err.Error()
			if strings.Contains(emsg, "already upgraded") {
				fmt.Printf("vault implementation is updated: %s\n", chainInfo.Chainconfig.VaultLogicAddress)
				err = nil
			} else {
				fmt.Println("upgrade vault implementation err: ", err)
				return err
			}
		} else {
			fmt.Printf("vault logic implementation upgrade from %s to %s\n", oldImpl, newImpl)
		}

		// init report online info
		err = CheckExistLastOnlineReportV2(cfg, configRoot, chainid)
		if err != nil {
			fmt.Println("check report status, err: ", err)
			return err
		}

		err = CheckHubDomainConfig(cfg, configRoot, chainid)
		if err != nil {
			fmt.Println("check report status, err: ", err)
			return err
		}
	}

	// init ip2location db
	if err := bindata.Init(); err != nil {
		// log init ip2location err
		fmt.Println("init ip2location err: ", err)
		log.Errorf("init ip2location err:%+v", err)
	}

	offline, _ := req.Options[offlineKwd].(bool)
	ipnsps, _ := req.Options[enableIPNSPubSubKwd].(bool)
	pubsub, _ := req.Options[enablePubSubKwd].(bool)
	if _, hasMplex := req.Options[enableMultiplexKwd]; hasMplex {
		log.Errorf("The mplex multiplexer has been enabled by default and the experimental %s flag has been removed.")
		log.Errorf("To disable this multiplexer, please configure `Swarm.Transports.Multiplexers'.")
	}

	// Start assembling node config
	ncfg := &core.BuildCfg{
		Repo:                        repo,
		Permanent:                   true, // It is temporary way to signify that node is permanent
		Online:                      !offline,
		DisableEncryptedConnections: unencrypted,
		ExtraOpts: map[string]bool{
			"pubsub": pubsub,
			"ipnsps": ipnsps,
		},
		//TODO(Kubuxu): refactor Online vs Offline by adding Permanent vs Ephemeral
	}

	// Check if swarm port is overriden
	// Continue with default setup on error
	swarmPort, ok := req.Options[swarmPortKwd].(string)
	if ok {
		sp := strings.Split(swarmPort, ":")
		if len(sp) != 2 {
			log.Errorf("Invalid swarm-port: %v", swarmPort)
		} else {
			wanPort, err1 := strconv.Atoi(sp[0])
			lanPort, err2 := strconv.Atoi(sp[1])
			if err1 != nil || err2 != nil {
				log.Errorf("Invalid swarm-port: %v", swarmPort)
			} else {
				err = ncfg.AnnouncePublicIpWithPort(wanPort, lanPort)
				if err != nil {
					log.Errorf("Failed to announce new swarm address: %v", err)
				}
			}
		}
	}

	routingOption, _ := req.Options[routingOptionKwd].(string)
	if routingOption == routingOptionDefaultKwd {
		routingOption = cfg.Routing.Type.String()
		if routingOption == "" {
			routingOption = routingOptionDHTKwd
		}
	}
	switch routingOption {
	case routingOptionSupernodeKwd:
		return errors.New("supernode routing was never fully implemented and has been removed")
	case routingOptionDHTClientKwd:
		ncfg.Routing = libp2p.DHTClientOption
	case routingOptionDHTKwd:
		ncfg.Routing = libp2p.DHTOption
	case routingOptionDHTServerKwd:
		ncfg.Routing = libp2p.DHTServerOption
	case routingOptionNoneKwd:
		ncfg.Routing = libp2p.NilRouterOption
	case routingOptionCustomKwd:
		ncfg.Routing = libp2p.ConstructDelegatedRouting(
			cfg.Routing.Routers,
			cfg.Routing.Methods,
			cfg.Identity.PeerID,
			cfg.Addresses.Swarm,
			cfg.Identity.PrivKey,
		)
	default:
		return fmt.Errorf("unrecognized routing option: %s", routingOption)
	}

	node, err := core.NewNode(req.Context, ncfg)
	if err != nil {
		log.Error("error from node construction: ", err)
		return err
	}
	node.IsDaemon = true

	//Check if there is a swarm.key at btfs loc. This would still print fingerprint if they created a swarm.key with the same values
	spath := filepath.Join(cctx.ConfigRoot, "swarm.key")
	if node.PNetFingerprint != nil && util.FileExists(spath) {
		fmt.Println("Swarm is limited to private network of peers with the swarm key")
		fmt.Printf("Swarm key fingerprint: %x\n", node.PNetFingerprint)
	}

	printSwarmAddrs(node)

	defer func() {
		// We wait for the node to close first, as the node has children
		// that it will wait for before closing, such as the API server.
		node.Close()

		select {
		case <-req.Context.Done():
			log.Info("Gracefully shut down daemon")
		default:
		}
	}()

	cctx.ConstructNode = func() (*core.IpfsNode, error) {
		return node, nil
	}

	// Start "core" plugins. We want to do this *before* starting the HTTP
	// API as the user may be relying on these plugins.
	err = cctx.Plugins.Start(node)
	if err != nil {
		return err
	}
	node.Process.AddChild(goprocess.WithTeardown(cctx.Plugins.Close))

	if SimpleMode == false {
		// if the guide server was started, shutdown it
		guide.TryShutdownServer()
	}

	// construct api endpoint - every time
	apiErrc, err := serveHTTPApi(req, cctx, SimpleMode)
	if err != nil {
		return err
	}

	// construct fuse mountpoints - if the user provided the --mount flag
	mount, _ := req.Options[mountKwd].(bool)
	if mount && offline {
		return cmds.Errorf(cmds.ErrClient, "mount is not currently supported in offline mode")
	}
	if mount {
		if err := mountFuse(req, cctx); err != nil {
			return err
		}
	}

	// repo blockstore GC - if --enable-gc flag is present
	gcErrc, err := maybeRunGC(req, node)
	if err != nil {
		return err
	}

	// construct http gateway
	gwErrc, err := serveHTTPGateway(req, cctx)
	if err != nil {
		return err
	}

	// construct http remote api - if it is set in the config
	var rapiErrc <-chan error
	if len(cfg.Addresses.RemoteAPI) > 0 {
		var err error
		rapiErrc, err = serveHTTPRemoteApi(req, cctx)
		if err != nil {
			return err
		}
	}

	// Add btfs version info to prometheus metrics
	var btfsInfoMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "btfs_info",
		Help: "BTFS version information.",
	}, []string{"version", "commit"})

	// Setting to 1 lets us multiply it with other stats to add the version labels
	btfsInfoMetric.With(prometheus.Labels{
		"version": version.CurrentVersionNumber,
		"commit":  version.CurrentCommit,
	}).Set(1)

	// initialize metrics collector
	prometheus.MustRegister(&corehttp.IpfsNodeCollector{Node: node})

	// The daemon is *finally* ready.
	fmt.Printf("Daemon is ready\n")
	notifyReady()

	runStartupTest, _ := req.Options[enableStartupTest].(bool)

	// BTFS functional test
	if runStartupTest {
		functest(cfg.Services.OnlineServerDomain, cfg.Identity.PeerID, hValue)
	}

	if SimpleMode == false {
		// set Analytics flag if specified
		if dc, ok := req.Options[enableDataCollection]; ok == true {
			node.Repo.SetConfigKey("Experimental.Analytics", dc)
		}
		// Spin jobs in the background
		spin.RenterSessions(req, env)
		api, err := cmdenv.GetApi(env, req)
		if err != nil {
			return err
		}

		spin.Analytics(api, cctx.ConfigRoot, node, version.CurrentVersionNumber, hValue)
		spin.Hosts(node, env)
		spin.Contracts(node, req, env, nodepb.ContractStat_HOST.String())
	}

	// Give the user some immediate feedback when they hit C-c
	go func() {
		<-req.Context.Done()
		notifyStopping()
		fmt.Println("Received interrupt signal, shutting down...")
		fmt.Println("(Hit ctrl-c again to force-shutdown the daemon.)")
	}()

	// collect long-running errors and block for shutdown
	// TODO(cryptix): our fuse currently doesn't follow this pattern for graceful shutdown
	var errs error
	for err := range merge(apiErrc, gwErrc, rapiErrc, gcErrc) {
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

// serveHTTPApi collects options, creates listener, prints status message and starts serving requests
func serveHTTPApi(req *cmds.Request, cctx *oldcmds.Context, SimpleMode bool) (<-chan error, error) {
	cfg, err := cctx.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("serveHTTPApi: GetConfig() failed: %s", err)
	}

	listeners, err := sockets.TakeListeners("io.ipfs.api")
	if err != nil {
		return nil, fmt.Errorf("serveHTTPApi: socket activation failed: %s", err)
	}

	apiAddrs := make([]string, 0, 2)
	apiAddr, _ := req.Options[commands.ApiOption].(string)
	if apiAddr == "" {
		apiAddrs = cfg.Addresses.API
	} else {
		apiAddrs = append(apiAddrs, apiAddr)
	}

	listenerAddrs := make(map[string]bool, len(listeners))
	for _, listener := range listeners {
		listenerAddrs[string(listener.Multiaddr().Bytes())] = true
	}

	for _, addr := range apiAddrs {
		apiMaddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return nil, fmt.Errorf("serveHTTPApi: invalid API address: %q (err: %s)", addr, err)
		}
		if listenerAddrs[string(apiMaddr.Bytes())] {
			continue
		}

		apiLis, err := manet.Listen(apiMaddr)
		if err != nil {
			return nil, fmt.Errorf("serveHTTPApi: manet.Listen(%s) failed: %s", apiMaddr, err)
		}

		listenerAddrs[string(apiMaddr.Bytes())] = true
		listeners = append(listeners, apiLis)
	}

	for _, listener := range listeners {
		// we might have listened to /tcp/0 - let's see what we are listing on
		fmt.Printf("API server listening on %s\n", listener.Multiaddr())
		// Browsers require TCP.
		switch listener.Addr().Network() {
		case "tcp", "tcp4", "tcp6":
			if SimpleMode == false {
				fmt.Printf("Dashboard: http://%s/dashboard\n", listener.Addr())
			}
		}
	}

	// by default, we don't let you load arbitrary btfs objects through the api,
	// because this would open up the api to scripting vulnerabilities.
	// only the webui objects are allowed.
	// if you know what you're doing, go ahead and pass --unrestricted-api.
	unrestricted, _ := req.Options[unrestrictedApiAccessKwd].(bool)
	gatewayOpt := corehttp.GatewayOption(corehttp.WebUIPaths...)
	if unrestricted {
		gatewayOpt = corehttp.GatewayOption("/btfs", "/btns")
	}

	var opts = []corehttp.ServeOption{
		corehttp.MetricsCollectionOption("api"),
		corehttp.CheckVersionOption(),
		corehttp.CommandsOption(*cctx),
		corehttp.WebUIOption,
		corehttp.DashboardOption,
		corehttp.HostUIOption,
		gatewayOpt,
		corehttp.VersionOption(),
		defaultMux("/debug/vars"),
		defaultMux("/debug/pprof/"),
		defaultMux("/debug/stack"),
		corehttp.MutexFractionOption("/debug/pprof-mutex/"),
		corehttp.BlockProfileRateOption("/debug/pprof-block/"),
		corehttp.MetricsScrapingOption("/debug/metrics/prometheus"),
		corehttp.LogOption(),
	}

	if len(cfg.Gateway.RootRedirect) > 0 {
		opts = append(opts, corehttp.RedirectOption("", cfg.Gateway.RootRedirect))
	}

	node, err := cctx.ConstructNode()
	if err != nil {
		return nil, fmt.Errorf("serveHTTPApi: ConstructNode() failed: %s", err)
	}

	if err := node.Repo.SetAPIAddr(listeners[0].Multiaddr()); err != nil {
		return nil, fmt.Errorf("serveHTTPApi: SetAPIAddr() failed: %s", err)
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	for _, apiLis := range listeners {
		wg.Add(1)
		go func(lis manet.Listener) {
			defer wg.Done()
			errc <- corehttp.Serve(node, manet.NetListener(lis), opts...)
		}(apiLis)
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	return errc, nil
}

func getInputChainID(req *cmds.Request) (chainid int64, err error) {
	inputChainIdStr, found := req.Options[chainID].(string)
	if found {
		inputChainid, err := strconv.ParseInt(inputChainIdStr, 10, 64)
		if err != nil {
			return 0, err
		}

		return inputChainid, nil
	}
	return 0, nil
}

func getChainID(req *cmds.Request, cfg *config.Config, stateStorer storage.StateStorer) (chainId int64, stored bool, err error) {
	cfgChainId := cfg.ChainInfo.ChainId
	inputChainId, err := getInputChainID(req)
	if err != nil {
		return 0, stored, err
	}
	storeChainid, err := chain.GetChainIdFromDisk(stateStorer)
	if err != nil {
		return 0, stored, err
	}

	chainId = cc.DefaultChain
	//config chain version, must be have cfgChainId
	if storeChainid > 0 {
		// on moving, from old version to new version, later must be sync info
		if cfgChainId == 0 && cfg.ChainInfo.Endpoint == "" {
			chainId = storeChainid
			stored = false
		} else { // not move, new version start.
			// compare cfg chain id and leveldb chain id
			if storeChainid != cfgChainId {
				return 0, stored, errors.New(
					fmt.Sprintf("current chainId=%d is different from config chainId=%d, "+
						"you can not change chain id in config file", storeChainid, cfgChainId))
			}

			// compare input chain id and leveldb chain id
			if inputChainId > 0 && storeChainid != inputChainId {
				return 0, stored, errors.New(
					fmt.Sprintf("current chainId=%d is different from input chainId=%d, "+
						"you can not change chain id with --chain-id when node start", storeChainid, inputChainId))
			}

			chainId = storeChainid
			stored = true
		}
	} else { // not move, old version start.
		// old version, should be inputChainId first, DefaultChainId second.
		if inputChainId > 0 {
			chainId = inputChainId
		}
		stored = false
	}

	return chainId, stored, nil
}

// printSwarmAddrs prints the addresses of the host
func printSwarmAddrs(node *core.IpfsNode) {
	if !node.IsOnline {
		fmt.Println("Swarm not listening, running in offline mode.")
		return
	}

	var lisAddrs []string
	ifaceAddrs, err := node.PeerHost.Network().InterfaceListenAddresses()
	if err != nil {
		log.Errorf("failed to read listening addresses: %s", err)
	}
	for _, addr := range ifaceAddrs {
		lisAddrs = append(lisAddrs, addr.String())
	}
	sort.Strings(lisAddrs)
	for _, addr := range lisAddrs {
		fmt.Printf("Swarm listening on %s\n", addr)
	}

	var addrs []string
	for _, addr := range node.PeerHost.Addrs() {
		addrs = append(addrs, addr.String())
	}
	sort.Strings(addrs)
	for _, addr := range addrs {
		fmt.Printf("Swarm announcing %s\n", addr)
	}
}

// serveHTTPGateway collects options, creates listener, prints status message and starts serving requests
func serveHTTPGateway(req *cmds.Request, cctx *oldcmds.Context) (<-chan error, error) {
	cfg, err := cctx.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("serveHTTPGateway: GetConfig() failed: %s", err)
	}

	writable, writableOptionFound := req.Options[writableKwd].(bool)
	if !writableOptionFound {
		writable = cfg.Gateway.Writable
	}
	if writable {
		log.Errorf("Support for Gateway.Writable and --writable has been REMOVED. Please remove it from your config file or CLI.")
	}
	listeners, err := sockets.TakeListeners("io.ipfs.gateway")
	if err != nil {
		return nil, fmt.Errorf("serveHTTPGateway: socket activation failed: %s", err)
	}

	listenerAddrs := make(map[string]bool, len(listeners))
	for _, listener := range listeners {
		listenerAddrs[string(listener.Multiaddr().Bytes())] = true
	}

	gatewayAddrs := cfg.Addresses.Gateway
	for _, addr := range gatewayAddrs {
		gatewayMaddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return nil, fmt.Errorf("serveHTTPGateway: invalid gateway address: %q (err: %s)", addr, err)
		}

		if listenerAddrs[string(gatewayMaddr.Bytes())] {
			continue
		}

		gwLis, err := manet.Listen(gatewayMaddr)
		if err != nil {
			return nil, fmt.Errorf("serveHTTPGateway: manet.Listen(%s) failed: %s", gatewayMaddr, err)
		}
		listenerAddrs[string(gatewayMaddr.Bytes())] = true
		listeners = append(listeners, gwLis)
	}

	for _, listener := range listeners {
		fmt.Printf("Gateway server listening on %s\n", listener.Multiaddr())
	}

	cmdctx := *cctx
	cmdctx.Gateway = true

	var opts = []corehttp.ServeOption{
		corehttp.MetricsCollectionOption("gateway"),
		corehttp.HostnameOption(),
		// TODO: rm writable
		corehttp.GatewayOption("/btfs", "/btns"),
		corehttp.VersionOption(),
		corehttp.CheckVersionOption(),
		corehttp.CommandsROOption(cmdctx),
	}

	if cfg.Experimental.P2pHttpProxy {
		opts = append(opts, corehttp.P2PProxyOption())
	}

	if len(cfg.Gateway.RootRedirect) > 0 {
		opts = append(opts, corehttp.RedirectOption("", cfg.Gateway.RootRedirect))
	}

	node, err := cctx.ConstructNode()
	if err != nil {
		return nil, fmt.Errorf("serveHTTPGateway: ConstructNode() failed: %s", err)
	}
	if len(cfg.Gateway.PathPrefixes) > 0 {
		log.Errorf("Support for custom Gateway.PathPrefixes was removed")
	}
	errc := make(chan error)
	var wg sync.WaitGroup
	for _, lis := range listeners {
		wg.Add(1)
		go func(lis manet.Listener) {
			defer wg.Done()
			errc <- corehttp.Serve(node, manet.NetListener(lis), opts...)
		}(lis)
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	return errc, nil
}

// serveHTTPRemoteApi collects options, creates listener, prints status message and starts serving requests
func serveHTTPRemoteApi(req *cmds.Request, cctx *oldcmds.Context) (<-chan error, error) {
	cfg, err := cctx.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("serveHTTPRemoteApi: GetConfig() failed: %s", err)
	}

	if !cfg.Experimental.Libp2pStreamMounting {
		return nil, fmt.Errorf("serveHTTPRemoteApi: libp2p stream mounting must be enabled")
	}

	rapiAddrs := cfg.Addresses.RemoteAPI
	listeners := make([]manet.Listener, 0, len(rapiAddrs))
	for _, addr := range rapiAddrs {
		rapiMaddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return nil, fmt.Errorf("serveHTTPRemoteApi: invalid remote api address: %q (err: %s)", addr, err)
		}

		rapiLis, err := manet.Listen(rapiMaddr)
		if err != nil {
			return nil, fmt.Errorf("serveHTTPRemoteApi: manet.Listen(%s) failed: %s", rapiMaddr, err)
		}
		// we might have listened to /tcp/0 - lets see what we are listing on
		rapiMaddr = rapiLis.Multiaddr()
		fmt.Printf("Remote API server listening on %s\n", rapiMaddr)

		listeners = append(listeners, rapiLis)
	}

	var opts = []corehttp.ServeOption{
		corehttp.MetricsCollectionOption("remote_api"),
		corehttp.VersionOption(),
		corehttp.CheckVersionOption(),
		corehttp.CommandsRemoteOption(*cctx),
	}

	node, err := cctx.ConstructNode()
	if err != nil {
		return nil, fmt.Errorf("serveHTTPRemoteApi: ConstructNode() failed: %s", err)
	}

	// set default listener to remote api endpoint
	if _, err := node.P2P.ForwardRemote(node.Context(),
		httpremote.P2PRemoteCallProto, listeners[0].Multiaddr(), false); err != nil {
		return nil, fmt.Errorf("serveHTTPRemoteApi: ForwardRemote() failed: %s", err)
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	for _, lis := range listeners {
		wg.Add(1)
		go func(lis manet.Listener) {
			defer wg.Done()
			errc <- corehttp.Serve(node, manet.NetListener(lis), opts...)
		}(lis)
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	return errc, nil
}

// collects options and opens the fuse mountpoint
func mountFuse(req *cmds.Request, cctx *oldcmds.Context) error {
	cfg, err := cctx.GetConfig()
	if err != nil {
		return fmt.Errorf("mountFuse: GetConfig() failed: %s", err)
	}

	fsdir, found := req.Options[ipfsMountKwd].(string)
	if !found {
		fsdir = cfg.Mounts.IPFS
	}

	nsdir, found := req.Options[ipnsMountKwd].(string)
	if !found {
		nsdir = cfg.Mounts.IPNS
	}

	node, err := cctx.ConstructNode()
	if err != nil {
		return fmt.Errorf("mountFuse: ConstructNode() failed: %s", err)
	}

	err = nodeMount.Mount(node, fsdir, nsdir)
	if err != nil {
		return err
	}
	fmt.Printf("BTFS mounted at: %s\n", fsdir)
	fmt.Printf("BTNS mounted at: %s\n", nsdir)
	return nil
}

func maybeRunGC(req *cmds.Request, node *core.IpfsNode) (<-chan error, error) {
	enableGC, _ := req.Options[enableGCKwd].(bool)
	if !enableGC {
		return nil, nil
	}

	errc := make(chan error)
	go func() {
		errc <- corerepo.PeriodicGC(req.Context, node)
		close(errc)
	}()
	return errc, nil
}

// merge does fan-in of multiple read-only error channels
// taken from http://blog.golang.org/pipelines
func merge(cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	for _, c := range cs {
		if c != nil {
			wg.Add(1)
			go output(c)
		}
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func YesNoPrompt(prompt string) bool {
	var s string
	for i := 0; i < 3; i++ {
		fmt.Printf("%s ", prompt)
		fmt.Scanf("%s", &s)
		switch s {
		case "y", "Y":
			return true
		case "n", "N":
			return false
		case "":
			return false
		}
		fmt.Println("Please press either 'y' or 'n'")
	}

	return false
}

func printVersion() {
	v := version.CurrentVersionNumber
	if version.CurrentCommit != "" {
		v += "-" + version.CurrentCommit
	}
	fmt.Printf("go-btfs version: %s\n", v)
	fmt.Printf("Repo version: %d\n", fsrepo.RepoVersion)
	fmt.Printf("System version: %s\n", runtime.GOARCH+"/"+runtime.GOOS)
	fmt.Printf("Golang version: %s\n", runtime.Version())
}

func getBtfsBinaryPath() (string, error) {
	defaultBtfsPath, err := getCurrentPath()
	if err != nil {
		log.Errorf("Get current program execution path error, reasons: [%v]", err)
		return "", err
	}

	ext := ""
	sep := "/"
	if runtime.GOOS == "windows" {
		ext = ".exe"
		sep = "\\"
	}
	latestBtfsBinary := "btfs" + ext
	latestBtfsBinaryPath := fmt.Sprint(defaultBtfsPath, sep, latestBtfsBinary)

	return latestBtfsBinaryPath, nil
}

func functest(onlineServerDomain, peerId, hValue string) {
	btfsBinaryPath, err := getBtfsBinaryPath()
	if err != nil {
		fmt.Printf("Get btfs path failed, BTFS daemon test skipped\n")
		os.Exit(findBTFSBinaryFailed)
	}

	// prepare functional test before start btfs daemon
	ready_to_test := prepare_test(btfsBinaryPath, onlineServerDomain, peerId, hValue)
	// start btfs functional test
	if ready_to_test {
		test_success := false
		// try up to two times
		for i := 0; i < 2; i++ {
			err := get_functest(btfsBinaryPath)
			if err != nil {
				fmt.Printf("BTFS daemon get file test failed! Reason: %v\n", err)
				SendError(err.Error(), onlineServerDomain, peerId, hValue)
			} else {
				fmt.Printf("BTFS daemon get file test succeeded!\n")
				test_success = true
				break
			}
		}
		if !test_success {
			fmt.Printf("BTFS daemon get file test failed twice! exiting\n")
			os.Exit(getFileTestFailed)
		}
		test_success = false
		// try up to two times
		for i := 0; i < 2; i++ {
			if err := add_functest(btfsBinaryPath, peerId); err != nil {
				fmt.Printf("BTFS daemon add file test failed! Reason: %v\n", err)
				SendError(err.Error(), onlineServerDomain, peerId, hValue)
			} else {
				fmt.Printf("BTFS daemon add file test succeeded!\n")
				test_success = true
				break
			}
		}
		if !test_success {
			fmt.Printf("BTFS daemon add file test failed twice! exiting\n")
			os.Exit(addFileTestFailed)
		}
	} else {
		fmt.Printf("BTFS daemon test skipped\n")
	}
}

// VaultFactory upgraded to V2, we need re-deploy a vault for user
func doIfNeedUpgradeFactoryToV2(chainid int64, chainCfg *chainconfig.ChainConfig, statestore storage.StateStorer, repo repo.Repo, cfg *config.Config, configRoot string) (need bool, err error) {

	currChainCfg, ok := chainconfig.GetChainConfig(chainid)
	if !ok {
		err = errors.New(fmt.Sprintf("chain %d is not supported yet", chainid))
		return
	}

	confFactory := cfg.ChainInfo.CurrentFactory      // Factory address from config file, may be ""
	currFactory := currChainCfg.CurrentFactory.Hex() // Factory address read from source code
	if !chainconfig.IsV2FactoryAddr(currFactory) {
		return
	}

	// calculate whether need to upgrade factory to v2
	if confFactory == "" {
		need = true
	} else {
		if confFactory == currFactory {
			need = false
		} else {
			need = chainconfig.IsV1FactoryAddr(confFactory)
		}
	}
	if !need {
		return
	}

	fmt.Println("prepare upgrading your vault contract")

	oldVault, err := vault.GetStoredVaultAddr(statestore)
	if err != nil {
		return
	}

	// backup `statestore` folder
	err = statestore.Close()
	if err != nil {
		return
	}

	bkSuffix := fmt.Sprintf("backup%d", rand.Intn(100))
	err = chain.BackUpStateStore(configRoot, bkSuffix)
	if err != nil {
		return
	}

	// backup `config` file
	var bkConfig string
	bkConfig, err = repo.BackUpConfigV2(bkSuffix)
	if err != nil {
		fmt.Printf("backup config file failed, err: %s\n", err)
		return
	}
	fmt.Printf("backup config file successfully to %s\n", bkConfig)

	// update factory address and other chain info to config file.
	// note that we only changed the `CurrentFactory`, so we won't overide other chaininfo field in the config file.
	chainCfg.CurrentFactory = common.HexToAddress(currFactory)

	cfg.ChainInfo.ChainId = chainid // set these fields here is a little tricky
	cfg.ChainInfo.Endpoint = chainCfg.Endpoint
	cfg.ChainInfo.CurrentFactory = chainCfg.CurrentFactory.Hex()
	cfg.ChainInfo.PriceOracleAddress = chainCfg.PriceOracleAddress.Hex()

	err = commands.SyncConfigChainInfoV2(configRoot, chainid, chainCfg.Endpoint, chainCfg.CurrentFactory, chainCfg.PriceOracleAddress)
	if err != nil {
		return
	}

	zeroaddr := common.Address{}
	if oldVault != zeroaddr {
		fmt.Printf("your old vault address is %s\n", oldVault)
	}
	fmt.Println("will re-deploy a vault contract for you")
	return
}

// CheckExistLastOnlineReport sync conf and lastOnlineInfo
func CheckExistLastOnlineReport(cfg *config.Config, configRoot string, chainId int64, reportStatusServ reportstatus.Service) error {
	lastOnline, err := chain.GetLastOnline()
	if err != nil {
		return err
	}

	// if nil, set config online status config
	if lastOnline == nil {
		var reportOnline bool
		var reportStatusContract bool
		if cfg.Experimental.StorageHostEnabled {
			reportOnline = true
			reportStatusContract = true
		}

		var onlineServerDomain string
		if chainId == 199 {
			onlineServerDomain = config.DefaultServicesConfig().OnlineServerDomain
		} else {
			onlineServerDomain = config.DefaultServicesConfigTestnet().OnlineServerDomain
		}

		err = commands.SyncConfigOnlineCfg(configRoot, onlineServerDomain, reportOnline, reportStatusContract)
		if err != nil {
			return err
		}
	}

	// if nil, set last online info
	if lastOnline == nil {
		err = reportStatusServ.CheckLastOnlineInfo(cfg.Identity.PeerID, cfg.Identity.BttcAddr)
		if err != nil {
			return err
		}
	}
	return nil
}

// CheckExistLastOnlineReport sync conf and lastOnlineInfo
func CheckHubDomainConfig(cfg *config.Config, configRoot string, chainId int64) error {
	var hubServerDomain string
	if chainId == 199 {
		hubServerDomain = config.DefaultServicesConfig().HubDomain
	} else {
		hubServerDomain = config.DefaultServicesConfigTestnet().HubDomain
	}

	if hubServerDomain != cfg.Services.HubDomain {
		err := commands.SyncHubDomainConfig(configRoot, hubServerDomain)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckExistLastOnlineReportV2 sync conf and lastOnlineInfo
func CheckExistLastOnlineReportV2(cfg *config.Config, configRoot string, chainId int64) error {
	lastOnline, err := chain.GetLastOnline()
	if err != nil {
		return err
	}

	// if nil, set config online status config
	if lastOnline == nil {
		var reportOnline bool
		if cfg.Experimental.StorageHostEnabled {
			reportOnline = true
		}

		var onlineServerDomain string
		if chainId == 199 {
			onlineServerDomain = config.DefaultServicesConfig().OnlineServerDomain
		} else {
			onlineServerDomain = config.DefaultServicesConfigTestnet().OnlineServerDomain
		}

		err = commands.SyncConfigOnlineCfgV2(configRoot, onlineServerDomain, reportOnline)
		if err != nil {
			return err
		}
	}

	// if nil, set last online info
	if lastOnline == nil {
		if err != spin.GetLastOnlineInfoWhenNodeMigration(context.Background(), cfg) {
			return err
		}
	}
	return nil
}
