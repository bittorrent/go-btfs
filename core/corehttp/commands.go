package corehttp

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	shell "github.com/bittorrent/go-btfs-api"

	version "github.com/bittorrent/go-btfs"
	oldcmds "github.com/bittorrent/go-btfs/commands"
	"github.com/bittorrent/go-btfs/core"
	corecommands "github.com/bittorrent/go-btfs/core/commands"

	cmds "github.com/bittorrent/go-btfs-cmds"
	cmdsHttp "github.com/bittorrent/go-btfs-cmds/http"
	config "github.com/bittorrent/go-btfs-config"
	path "github.com/ipfs/go-path"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	errAPIVersionMismatch = errors.New("api version mismatch")
)

const originEnvKey = "API_ORIGIN"
const originEnvKeyDeprecate = `You are using the ` + originEnvKey + `ENV Variable.
This functionality is deprecated, and will be removed in future versions.
Instead, try either adding headers to the config, or passing them via
cli arguments:

	btfs config API.HTTPHeaders --json '{"Access-Control-Allow-Origin": ["*"]}'
	btfs daemon
`

// APIPath is the path at which the API is mounted.
const APIPath = "/api/" + shell.API_VERSION

// redirectPaths contain the paths to be redirectd to the current APIPath
var redirectPaths = []string{
	"/api/v0",
	"/api/latest",
}

var defaultLocalhostOrigins = []string{
	"http://127.0.0.1:<port>",
	"https://127.0.0.1:<port>",
	"http://[::1]:<port>",
	"https://[::1]:<port>",
	"http://localhost:<port>",
	"https://localhost:<port>",
}

func addCORSFromEnv(c *cmdsHttp.ServerConfig) {
	origin := os.Getenv(originEnvKey)
	if origin != "" {
		log.Warn(originEnvKeyDeprecate)
		c.AppendAllowedOrigins(origin)
	}
}

func addHeadersFromConfig(c *cmdsHttp.ServerConfig, nc *config.Config) {
	log.Info("Using API.HTTPHeaders:", nc.API.HTTPHeaders)

	if acao := nc.API.HTTPHeaders[cmdsHttp.ACAOrigin]; acao != nil {
		c.SetAllowedOrigins(acao...)
	}
	if acam := nc.API.HTTPHeaders[cmdsHttp.ACAMethods]; acam != nil {
		c.SetAllowedMethods(acam...)
	}
	for _, v := range nc.API.HTTPHeaders[cmdsHttp.ACACredentials] {
		c.SetAllowCredentials(strings.ToLower(v) == "true")
	}

	c.Headers = make(map[string][]string, len(nc.API.HTTPHeaders)+1)

	// Copy these because the config is shared and this function is called
	// in multiple places concurrently. Updating these in-place *is* racy.
	for h, v := range nc.API.HTTPHeaders {
		h = http.CanonicalHeaderKey(h)
		switch h {
		case cmdsHttp.ACAOrigin, cmdsHttp.ACAMethods, cmdsHttp.ACACredentials:
			// these are handled by the CORs library.
		default:
			c.Headers[h] = v
		}
	}
	c.Headers["Server"] = []string{"go-btfs/" + version.CurrentVersionNumber}
}

func addCORSDefaults(c *cmdsHttp.ServerConfig) {
	// by default use localhost origins
	if len(c.AllowedOrigins()) == 0 {
		c.SetAllowedOrigins(defaultLocalhostOrigins...)
	}

	// by default, use GET, PUT, POST
	if len(c.AllowedMethods()) == 0 {
		c.SetAllowedMethods(http.MethodGet, http.MethodPost, http.MethodPut)
	}
}

func patchCORSVars(c *cmdsHttp.ServerConfig, addr net.Addr) {

	// we have to grab the port from an addr, which may be an ip6 addr.
	// TODO: this should take multiaddrs and derive port from there.
	port := ""
	if tcpaddr, ok := addr.(*net.TCPAddr); ok {
		port = strconv.Itoa(tcpaddr.Port)
	} else if udpaddr, ok := addr.(*net.UDPAddr); ok {
		port = strconv.Itoa(udpaddr.Port)
	}

	// we're listening on tcp/udp with ports. ("udp!?" you say? yeah... it happens...)
	oldOrigins := c.AllowedOrigins()
	newOrigins := make([]string, len(oldOrigins))
	for i, o := range oldOrigins {
		// TODO: allow replacing <host>. tricky, ip4 and ip6 and hostnames...
		if port != "" {
			o = strings.Replace(o, "<port>", port, -1)
		}
		newOrigins[i] = o
	}
	c.SetAllowedOrigins(newOrigins...)
}

func commandsOption(cctx oldcmds.Context, command *cmds.Command, allowGet bool) ServeOption {
	return func(n *core.IpfsNode, l net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {

		cfg := cmdsHttp.NewServerConfig()
		cfg.AllowGet = allowGet
		corsAllowedMethods := []string{http.MethodPost}
		if allowGet {
			corsAllowedMethods = append(corsAllowedMethods, http.MethodGet)
		}

		cfg.SetAllowedMethods(corsAllowedMethods...)
		cfg.APIPath = APIPath
		cfg.RedirectPaths = redirectPaths
		rcfg, err := n.Repo.Config()
		if err != nil {
			return nil, err
		}

		addHeadersFromConfig(cfg, rcfg)
		addCORSFromEnv(cfg)
		addCORSDefaults(cfg)
		patchCORSVars(cfg, l.Addr())

		cmdHandler := cmdsHttp.NewHandler(&cctx, command, cfg)
		cmdHandler = otelhttp.NewHandler(cmdHandler, "corehttp.cmdsHandler")
		mux.Handle(APIPath+"/", cmdHandler)
		for _, rp := range redirectPaths {
			mux.Handle(rp+"/", cmdHandler)
		}
		return mux, nil
	}
}

// CommandsOption constructs a ServerOption for hooking the commands into the
// HTTP server. It will NOT allow GET requests.
func CommandsOption(cctx oldcmds.Context) ServeOption {
	return commandsOption(cctx, corecommands.Root, false)
}

// CommandsROOption constructs a ServerOption for hooking the read-only commands
// into the HTTP server. It will allow GET requests.
func CommandsROOption(cctx oldcmds.Context) ServeOption {
	return commandsOption(cctx, corecommands.RootRO, true)
}

// CommandsRemoteOption constructs a ServerOption for hooking the public-facing,
// remote commands into the HTTP server.
func CommandsRemoteOption(cctx oldcmds.Context) ServeOption {
	return commandsOption(cctx, corecommands.RootRemote, false)
}

// CheckVersionOption returns a ServeOption that checks whether the client btfs version matches. Does nothing when the user agent string does not contain `/go-btfs/`
func CheckVersionOption() ServeOption {
	daemonVersion := version.ApiVersion

	return ServeOption(func(n *core.IpfsNode, l net.Listener, parent *http.ServeMux) (*http.ServeMux, error) {
		mux := http.NewServeMux()
		parent.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, APIPath) {
				cmdqry := r.URL.Path[len(APIPath):]
				pth := path.SplitList(cmdqry)

				// backwards compatibility to previous version check
				if len(pth) >= 2 && pth[1] != "version" {
					clientVersion := r.UserAgent()
					// skips check if client is not go-btfs
					if strings.Contains(clientVersion, "/go-btfs/") && daemonVersion != clientVersion {
						http.Error(w, fmt.Sprintf("%s (%s != %s)", errAPIVersionMismatch, daemonVersion, clientVersion), http.StatusBadRequest)
						return
					}
				}
			}

			mux.ServeHTTP(w, r)
		})

		return mux, nil
	})
}
