package libp2p

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/benbjohnson/clock"
	serialize "github.com/bittorrent/go-btfs-config/serialize"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	rcmgrObs "github.com/libp2p/go-libp2p/p2p/host/resource-manager/obs"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/fx"

	config "github.com/bittorrent/go-btfs-config"
	"github.com/bittorrent/go-btfs/core/node/helpers"
	"github.com/bittorrent/go-btfs/repo"
)

const NetLimitDefaultFilename = "limit.json"
const NetLimitTraceFilename = "rcmgr.json.gz"

var ErrNoResourceMgr = fmt.Errorf("missing ResourceMgr: make sure the daemon is running with Swarm.ResourceMgr.Enabled")

func ResourceManager(cfg config.SwarmConfig) interface{} {
	return func(mctx helpers.MetricsCtx, lc fx.Lifecycle, repo repo.Repo) (network.ResourceManager, Libp2pOpts, error) {
		var manager network.ResourceManager
		var opts Libp2pOpts

		enabled := cfg.ResourceMgr.Enabled.WithDefault(true)

		//  ENV overrides Config (if present)
		switch os.Getenv("LIBP2P_RCMGR") {
		case "0", "false":
			enabled = false
		case "1", "true":
			enabled = true
		}

		if enabled {
			log.Debug("libp2p resource manager is enabled")

			repoPath, err := config.PathRoot()
			if err != nil {
				return nil, opts, fmt.Errorf("opening BTFS_PATH: %w", err)
			}

			limitConfig, err := createDefaultLimitConfig(cfg)
			if err != nil {
				return nil, opts, err
			}

			partialLimitConfig := rcmgr.PartialLimitConfig{}
			err = serialize.ReadConfigFile(filepath.Join(repoPath, "libp2p-resource-limit-overrides.json"), &partialLimitConfig)
			if errors.Is(err, serialize.ErrNotInitialized) {
				err = nil
			} else {
				if err != nil {
					return nil, opts, err
				}
			}

			limitConfig = partialLimitConfig.Build(limitConfig)

			// The logic for defaults and overriding with specified SwarmConfig.ResourceMgr.Limits
			// is documented in docs/config.md.
			// Any changes here should be reflected there.
			// if cfg.ResourceMgr.Limits != nil {
			// 	l := *cfg.ResourceMgr.Limits
			// 	// This effectively overrides the computed default LimitConfig with any vlues from cfg.ResourceMgr.Limits
			// 	l.Apply(limitConfig)
			// 	limitConfig = l
			// }

			limiter := rcmgr.NewFixedLimiter(limitConfig)

			str, err := rcmgrObs.NewStatsTraceReporter()
			if err != nil {
				return nil, opts, err
			}

			ropts := []rcmgr.Option{rcmgr.WithMetrics(createRcmgrMetrics()), rcmgr.WithTraceReporter(str)}

			if len(cfg.ResourceMgr.Allowlist) > 0 {
				var mas []multiaddr.Multiaddr
				for _, maStr := range cfg.ResourceMgr.Allowlist {
					ma, err := multiaddr.NewMultiaddr(maStr)
					if err != nil {
						log.Errorf("failed to parse multiaddr=%v for allowlist, skipping. err=%v", maStr, err)
						continue
					}
					mas = append(mas, ma)
				}
				ropts = append(ropts, rcmgr.WithAllowlistedMultiaddrs(mas))
				log.Infof("Setting allowlist to: %v", mas)
			}

			// err = view.Register(rcmgrObs.DefaultViews...)
			// if err != nil {
			// 	return nil, opts, fmt.Errorf("registering rcmgr obs views: %w", err)
			// }

			if os.Getenv("LIBP2P_DEBUG_RCMGR") != "" {
				traceFilePath := filepath.Join(repoPath, NetLimitTraceFilename)
				ropts = append(ropts, rcmgr.WithTrace(traceFilePath))
			}

			manager, err = rcmgr.NewResourceManager(limiter, ropts...)
			if err != nil {
				return nil, opts, fmt.Errorf("creating libp2p resource manager: %w", err)
			}
			lrm := &loggingResourceManager{
				clock:    clock.New(),
				logger:   &logging.Logger("resourcemanager").SugaredLogger,
				delegate: manager,
			}
			lrm.start(helpers.LifecycleCtx(mctx, lc))
			manager = lrm
		} else {
			log.Error("libp2p resource manager is disabled")
			// manager = network.NullResourceManager
		}

		opts.Opts = append(opts.Opts, libp2p.ResourceManager(manager))

		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return manager.Close()
			}})

		return manager, opts, nil
	}
}

type NetStatOut struct {
	System    *rcmgr.BaseLimit           `json:",omitempty"`
	Transient *rcmgr.BaseLimit           `json:",omitempty"`
	Services  map[string]rcmgr.BaseLimit `json:",omitempty"`
	Protocols map[string]rcmgr.BaseLimit `json:",omitempty"`
	Peers     map[string]rcmgr.BaseLimit `json:",omitempty"`
}

func NetStat(mgr network.ResourceManager, scope string, percentage int) (NetStatOut, error) {
	var err error
	var result NetStatOut
	switch {
	case scope == "all":
		rapi, ok := mgr.(rcmgr.ResourceManagerState)
		if !ok { // NullResourceManager
			return result, ErrNoResourceMgr
		}

		limits, err := NetLimitAll(mgr)
		if err != nil {
			return result, err
		}

		stat := rapi.Stat()
		result.System = compareLimits(scopeToLimit(&stat.System), limits.System, percentage)
		result.Transient = compareLimits(scopeToLimit(&stat.Transient), limits.Transient, percentage)
		if len(stat.Services) > 0 {
			result.Services = make(map[string]rcmgr.BaseLimit, len(stat.Services))
			for srv, stat := range stat.Services {
				ls := limits.Services[srv]
				fstat := compareLimits(scopeToLimit(&stat), &ls, percentage)
				if fstat != nil {
					result.Services[srv] = *fstat
				}
			}
		}
		if len(stat.Protocols) > 0 {
			result.Protocols = make(map[string]rcmgr.BaseLimit, len(stat.Protocols))
			for proto, stat := range stat.Protocols {
				ls := limits.Protocols[string(proto)]
				fstat := compareLimits(scopeToLimit(&stat), &ls, percentage)
				if fstat != nil {
					result.Protocols[string(proto)] = *fstat
				}
			}
		}
		if len(stat.Peers) > 0 {
			result.Peers = make(map[string]rcmgr.BaseLimit, len(stat.Peers))
			for p, stat := range stat.Peers {
				ls := limits.Peers[p.String()]
				fstat := compareLimits(scopeToLimit(&stat), &ls, percentage)
				if fstat != nil {
					result.Peers[p.String()] = *fstat
				}
			}
		}

		return result, nil

	case scope == config.ResourceMgrSystemScope:
		err = mgr.ViewSystem(func(s network.ResourceScope) error {
			stat := s.Stat()
			result.System = scopeToLimit(&stat)
			return nil
		})
		return result, err

	case scope == config.ResourceMgrTransientScope:
		err = mgr.ViewTransient(func(s network.ResourceScope) error {
			stat := s.Stat()
			result.Transient = scopeToLimit(&stat)
			return nil
		})
		return result, err

	case strings.HasPrefix(scope, config.ResourceMgrServiceScopePrefix):
		svc := strings.TrimPrefix(scope, config.ResourceMgrServiceScopePrefix)
		err = mgr.ViewService(svc, func(s network.ServiceScope) error {
			stat := s.Stat()
			result.Services = map[string]rcmgr.BaseLimit{
				svc: *scopeToLimit(&stat),
			}
			return nil
		})
		return result, err

	case strings.HasPrefix(scope, config.ResourceMgrProtocolScopePrefix):
		proto := strings.TrimPrefix(scope, config.ResourceMgrProtocolScopePrefix)
		err = mgr.ViewProtocol(protocol.ID(proto), func(s network.ProtocolScope) error {
			stat := s.Stat()
			result.Protocols = map[string]rcmgr.BaseLimit{
				proto: *scopeToLimit(&stat),
			}
			return nil
		})
		return result, err

	case strings.HasPrefix(scope, config.ResourceMgrPeerScopePrefix):
		p := strings.TrimPrefix(scope, config.ResourceMgrPeerScopePrefix)
		pid, err := peer.Decode(p)
		if err != nil {
			return result, fmt.Errorf("invalid peer ID: %q: %w", p, err)
		}
		err = mgr.ViewPeer(pid, func(s network.PeerScope) error {
			stat := s.Stat()
			result.Peers = map[string]rcmgr.BaseLimit{
				p: *scopeToLimit(&stat),
			}
			return nil
		})
		return result, err

	default:
		return result, fmt.Errorf("invalid scope %q", scope)
	}
}

var scopes = []string{
	config.ResourceMgrSystemScope,
	config.ResourceMgrTransientScope,
	config.ResourceMgrServiceScopePrefix,
	config.ResourceMgrProtocolScopePrefix,
	config.ResourceMgrPeerScopePrefix,
}

func scopeToLimit(s *network.ScopeStat) *rcmgr.BaseLimit {
	return &rcmgr.BaseLimit{
		Streams:         s.NumStreamsInbound + s.NumStreamsOutbound,
		StreamsInbound:  s.NumStreamsInbound,
		StreamsOutbound: s.NumStreamsOutbound,
		Conns:           s.NumConnsInbound + s.NumConnsOutbound,
		ConnsInbound:    s.NumConnsInbound,
		ConnsOutbound:   s.NumConnsOutbound,
		FD:              s.NumFD,
		Memory:          s.Memory,
	}
}

// compareLimits compares stat and limit.
// If any of the stats value are equals or above the specified percentage,
// stat object is returned.
func compareLimits(stat, limit *rcmgr.BaseLimit, percentage int) *rcmgr.BaseLimit {
	if stat == nil || limit == nil {
		return nil
	}
	if abovePercentage(int(stat.Memory), int(limit.Memory), percentage) {
		return stat
	}
	if abovePercentage(stat.ConnsInbound, limit.ConnsInbound, percentage) {
		return stat
	}
	if abovePercentage(stat.ConnsOutbound, limit.ConnsOutbound, percentage) {
		return stat
	}
	if abovePercentage(stat.Conns, limit.Conns, percentage) {
		return stat
	}
	if abovePercentage(stat.FD, limit.FD, percentage) {
		return stat
	}
	if abovePercentage(stat.StreamsInbound, limit.StreamsInbound, percentage) {
		return stat
	}
	if abovePercentage(stat.StreamsOutbound, limit.StreamsOutbound, percentage) {
		return stat
	}
	if abovePercentage(stat.Streams, limit.Streams, percentage) {
		return stat
	}

	return nil
}

func abovePercentage(v1, v2, percentage int) bool {
	if percentage == 0 {
		return true
	}

	if v2 == 0 {
		return false
	}

	return int((v1/v2))*100 >= percentage
}

func NetLimitAll(mgr network.ResourceManager) (*NetStatOut, error) {
	var result = &NetStatOut{}
	lister, ok := mgr.(rcmgr.ResourceManagerState)
	if !ok { // NullResourceManager
		return result, ErrNoResourceMgr
	}

	for _, s := range scopes {
		switch s {
		case config.ResourceMgrSystemScope:
			s, err := NetLimit(mgr, config.ResourceMgrSystemScope)
			if err != nil {
				return nil, err
			}
			result.System = &s
		case config.ResourceMgrTransientScope:
			s, err := NetLimit(mgr, config.ResourceMgrSystemScope)
			if err != nil {
				return nil, err
			}
			result.Transient = &s
		case config.ResourceMgrServiceScopePrefix:
			result.Services = make(map[string]rcmgr.BaseLimit)
			for _, serv := range lister.ListServices() {
				s, err := NetLimit(mgr, config.ResourceMgrServiceScopePrefix+serv)
				if err != nil {
					return nil, err
				}
				result.Services[serv] = s
			}
		case config.ResourceMgrProtocolScopePrefix:
			result.Protocols = make(map[string]rcmgr.BaseLimit)
			for _, prot := range lister.ListProtocols() {
				ps := string(prot)
				s, err := NetLimit(mgr, config.ResourceMgrProtocolScopePrefix+ps)
				if err != nil {
					return nil, err
				}
				result.Protocols[ps] = s
			}
		case config.ResourceMgrPeerScopePrefix:
			result.Peers = make(map[string]rcmgr.BaseLimit)
			for _, peer := range lister.ListPeers() {
				ps := peer.String()
				s, err := NetLimit(mgr, config.ResourceMgrPeerScopePrefix+ps)
				if err != nil {
					return nil, err
				}
				result.Peers[ps] = s
			}
		}
	}

	return result, nil
}

func NetLimit(mgr network.ResourceManager, scope string) (rcmgr.BaseLimit, error) {
	var result rcmgr.BaseLimit
	getLimit := func(s network.ResourceScope) error {
		limiter, ok := s.(rcmgr.ResourceScopeLimiter)
		if !ok { // NullResourceManager
			return ErrNoResourceMgr
		}
		limit := limiter.Limit()
		switch l := limit.(type) {
		case *rcmgr.BaseLimit:
			result.Memory = l.Memory
			result.Streams = l.Streams
			result.StreamsInbound = l.StreamsInbound
			result.StreamsOutbound = l.StreamsOutbound
			result.Conns = l.Conns
			result.ConnsInbound = l.ConnsInbound
			result.ConnsOutbound = l.ConnsOutbound
			result.FD = l.FD
		default:
			return fmt.Errorf("unknown limit type %T", limit)
		}

		return nil
	}

	switch {
	case scope == config.ResourceMgrSystemScope:
		return result, mgr.ViewSystem(func(s network.ResourceScope) error { return getLimit(s) })
	case scope == config.ResourceMgrTransientScope:
		return result, mgr.ViewTransient(func(s network.ResourceScope) error { return getLimit(s) })
	case strings.HasPrefix(scope, config.ResourceMgrServiceScopePrefix):
		svc := strings.TrimPrefix(scope, config.ResourceMgrServiceScopePrefix)
		return result, mgr.ViewService(svc, func(s network.ServiceScope) error { return getLimit(s) })
	case strings.HasPrefix(scope, config.ResourceMgrProtocolScopePrefix):
		proto := strings.TrimPrefix(scope, config.ResourceMgrProtocolScopePrefix)
		return result, mgr.ViewProtocol(protocol.ID(proto), func(s network.ProtocolScope) error { return getLimit(s) })
	case strings.HasPrefix(scope, config.ResourceMgrPeerScopePrefix):
		p := strings.TrimPrefix(scope, config.ResourceMgrPeerScopePrefix)
		pid, err := peer.Decode(p)
		if err != nil {
			return result, fmt.Errorf("invalid peer ID: %q: %w", p, err)
		}
		return result, mgr.ViewPeer(pid, func(s network.PeerScope) error { return getLimit(s) })
	default:
		return result, fmt.Errorf("invalid scope %q", scope)
	}
}

// NetSetLimit sets new ResourceManager limits for the given scope. The limits take effect immediately, and are also persisted to the repo config.
// func NetSetLimit(mgr network.ResourceManager, repo repo.Repo, scope string, limit rcmgr.BaseLimit) error {
// 	setLimit := func(s network.ResourceScope) error {
// 		limiter, ok := s.(rcmgr.ResourceScopeLimiter)
// 		if !ok { // NullResourceManager
// 			return ErrNoResourceMgr
// 		}

// 		limiter.SetLimit(&limit)
// 		return nil
// 	}

// 	cfg, err := repo.Config()
// 	if err != nil {
// 		return fmt.Errorf("reading config to set limit: %w", err)
// 	}

// 	if cfg.Swarm.ResourceMgr.Limits == nil {
// 		cfg.Swarm.ResourceMgr.Limits = &rcmgr.ConcreteLimitConfig{}
// 	}
// 	configLimits := cfg.Swarm.ResourceMgr.Limits

// 	var setConfigFunc func()
// 	switch {
// 	case scope == config.ResourceMgrSystemScope:
// 		err = mgr.ViewSystem(func(s network.ResourceScope) error { return setLimit(s) })
// 		setConfigFunc = func() { configLimits.System = limit }
// 	case scope == config.ResourceMgrTransientScope:
// 		err = mgr.ViewTransient(func(s network.ResourceScope) error { return setLimit(s) })
// 		setConfigFunc = func() { configLimits.Transient = limit }
// 	case strings.HasPrefix(scope, config.ResourceMgrServiceScopePrefix):
// 		svc := strings.TrimPrefix(scope, config.ResourceMgrServiceScopePrefix)
// 		err = mgr.ViewService(svc, func(s network.ServiceScope) error { return setLimit(s) })
// 		setConfigFunc = func() {
// 			if configLimits.Service == nil {
// 				configLimits.Service = map[string]rcmgr.BaseLimit{}
// 			}
// 			configLimits.Service[svc] = limit
// 		}
// 	case strings.HasPrefix(scope, config.ResourceMgrProtocolScopePrefix):
// 		proto := strings.TrimPrefix(scope, config.ResourceMgrProtocolScopePrefix)
// 		err = mgr.ViewProtocol(protocol.ID(proto), func(s network.ProtocolScope) error { return setLimit(s) })
// 		setConfigFunc = func() {
// 			if configLimits.Protocol == nil {
// 				configLimits.Protocol = map[protocol.ID]rcmgr.BaseLimit{}
// 			}
// 			configLimits.Protocol[protocol.ID(proto)] = limit
// 		}
// 	case strings.HasPrefix(scope, config.ResourceMgrPeerScopePrefix):
// 		p := strings.TrimPrefix(scope, config.ResourceMgrPeerScopePrefix)
// 		var pid peer.ID
// 		pid, err = peer.Decode(p)
// 		if err != nil {
// 			return fmt.Errorf("invalid peer ID: %q: %w", p, err)
// 		}
// 		err = mgr.ViewPeer(pid, func(s network.PeerScope) error { return setLimit(s) })
// 		setConfigFunc = func() {
// 			if configLimits.Peer == nil {
// 				configLimits.Peer = map[peer.ID]rcmgr.BaseLimit{}
// 			}
// 			configLimits.Peer[pid] = limit
// 		}
// 	default:
// 		return fmt.Errorf("invalid scope %q", scope)
// 	}

// 	if err != nil {
// 		return fmt.Errorf("setting new limits on resource manager: %w", err)
// 	}

// 	if cfg.Swarm.ResourceMgr.Limits == nil {
// 		cfg.Swarm.ResourceMgr.Limits = &rcmgr.ConcreteLimitConfig{}
// 	}
// 	setConfigFunc()

// 	if err := repo.SetConfig(cfg); err != nil {
// 		return fmt.Errorf("writing new limits to repo config: %w", err)
// 	}

// 	return nil
// }

func LimitConfig(cfg config.SwarmConfig, userResourceOverrides rcmgr.PartialLimitConfig) (limitConfig rcmgr.ConcreteLimitConfig, err error) {
	limitConfig, err = createDefaultLimitConfig(cfg)
	if err != nil {
		return rcmgr.ConcreteLimitConfig{}, err
	}

	// The logic for defaults and overriding with specified userResourceOverrides
	// is documented in docs/libp2p-resource-management.md.
	// Any changes here should be reflected there.

	// This effectively overrides the computed default LimitConfig with any non-"useDefault" values from the userResourceOverrides file.
	// Because of how how Build works, any rcmgr.Default value in userResourceOverrides
	// will be overridden with a computed default value.
	limitConfig = userResourceOverrides.Build(limitConfig)

	return limitConfig, nil
}

type ResourceLimitsAndUsage struct {
	// This is duplicated from rcmgr.ResourceResourceLimits but adding *Usage fields.
	Memory               rcmgr.LimitVal64
	MemoryUsage          int64
	FD                   rcmgr.LimitVal
	FDUsage              int
	Conns                rcmgr.LimitVal
	ConnsUsage           int
	ConnsInbound         rcmgr.LimitVal
	ConnsInboundUsage    int
	ConnsOutbound        rcmgr.LimitVal
	ConnsOutboundUsage   int
	Streams              rcmgr.LimitVal
	StreamsUsage         int
	StreamsInbound       rcmgr.LimitVal
	StreamsInboundUsage  int
	StreamsOutbound      rcmgr.LimitVal
	StreamsOutboundUsage int
}

type LimitsConfigAndUsage struct {
	// This is duplicated from rcmgr.ResourceManagerStat but using ResourceLimitsAndUsage
	// instead of network.ScopeStat.
	System    ResourceLimitsAndUsage                 `json:",omitempty"`
	Transient ResourceLimitsAndUsage                 `json:",omitempty"`
	Services  map[string]ResourceLimitsAndUsage      `json:",omitempty"`
	Protocols map[protocol.ID]ResourceLimitsAndUsage `json:",omitempty"`
	Peers     map[peer.ID]ResourceLimitsAndUsage     `json:",omitempty"`
}

func MergeLimitsAndStatsIntoLimitsConfigAndUsage(l rcmgr.ConcreteLimitConfig, stats rcmgr.ResourceManagerStat) LimitsConfigAndUsage {
	limits := l.ToPartialLimitConfig()

	return LimitsConfigAndUsage{
		System:    mergeResourceLimitsAndScopeStatToResourceLimitsAndUsage(limits.System, stats.System),
		Transient: mergeResourceLimitsAndScopeStatToResourceLimitsAndUsage(limits.Transient, stats.Transient),
		Services:  mergeLimitsAndStatsMapIntoLimitsConfigAndUsageMap(limits.Service, stats.Services),
		Protocols: mergeLimitsAndStatsMapIntoLimitsConfigAndUsageMap(limits.Protocol, stats.Protocols),
		Peers:     mergeLimitsAndStatsMapIntoLimitsConfigAndUsageMap(limits.Peer, stats.Peers),
	}
}

func mergeResourceLimitsAndScopeStatToResourceLimitsAndUsage(rl rcmgr.ResourceLimits, ss network.ScopeStat) ResourceLimitsAndUsage {
	return ResourceLimitsAndUsage{
		Memory:               rl.Memory,
		MemoryUsage:          ss.Memory,
		FD:                   rl.FD,
		FDUsage:              ss.NumFD,
		Conns:                rl.Conns,
		ConnsUsage:           ss.NumConnsOutbound + ss.NumConnsInbound,
		ConnsOutbound:        rl.ConnsOutbound,
		ConnsOutboundUsage:   ss.NumConnsOutbound,
		ConnsInbound:         rl.ConnsInbound,
		ConnsInboundUsage:    ss.NumConnsInbound,
		Streams:              rl.Streams,
		StreamsUsage:         ss.NumStreamsOutbound + ss.NumStreamsInbound,
		StreamsOutbound:      rl.StreamsOutbound,
		StreamsOutboundUsage: ss.NumStreamsOutbound,
		StreamsInbound:       rl.StreamsInbound,
		StreamsInboundUsage:  ss.NumStreamsInbound,
	}
}

func mergeLimitsAndStatsMapIntoLimitsConfigAndUsageMap[K comparable](limits map[K]rcmgr.ResourceLimits, stats map[K]network.ScopeStat) map[K]ResourceLimitsAndUsage {
	r := make(map[K]ResourceLimitsAndUsage, maxInt(len(limits), len(stats)))
	for p, s := range stats {
		var l rcmgr.ResourceLimits
		if limits != nil {
			if rl, ok := limits[p]; ok {
				l = rl
			}
		}
		r[p] = mergeResourceLimitsAndScopeStatToResourceLimitsAndUsage(l, s)
	}
	for p, s := range limits {
		if _, ok := stats[p]; ok {
			continue // we already processed this element in the loop above
		}

		r[p] = mergeResourceLimitsAndScopeStatToResourceLimitsAndUsage(s, network.ScopeStat{})
	}
	return r
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// LimitConfigsToInfo gets limits and stats and generates a list of scopes and limits to be printed.
func LimitConfigsToInfo(stats LimitsConfigAndUsage) ResourceInfos {
	result := ResourceInfos{}

	result = append(result, resourceLimitsAndUsageToResourceInfo(config.ResourceMgrSystemScope, stats.System)...)
	result = append(result, resourceLimitsAndUsageToResourceInfo(config.ResourceMgrTransientScope, stats.Transient)...)

	for i, s := range stats.Services {
		result = append(result, resourceLimitsAndUsageToResourceInfo(
			config.ResourceMgrServiceScopePrefix+i,
			s,
		)...)
	}

	for i, p := range stats.Protocols {
		result = append(result, resourceLimitsAndUsageToResourceInfo(
			config.ResourceMgrProtocolScopePrefix+string(i),
			p,
		)...)
	}

	for i, p := range stats.Peers {
		result = append(result, resourceLimitsAndUsageToResourceInfo(
			config.ResourceMgrPeerScopePrefix+i.String(),
			p,
		)...)
	}

	return result
}

type ResourceInfo struct {
	ScopeName    string
	LimitName    string
	LimitValue   rcmgr.LimitVal64
	CurrentUsage int64
}

type ResourceInfos []ResourceInfo

const (
	limitNameMemory          = "Memory"
	limitNameFD              = "FD"
	limitNameConns           = "Conns"
	limitNameConnsInbound    = "ConnsInbound"
	limitNameConnsOutbound   = "ConnsOutbound"
	limitNameStreams         = "Streams"
	limitNameStreamsInbound  = "StreamsInbound"
	limitNameStreamsOutbound = "StreamsOutbound"
)

var limits = []string{
	limitNameMemory,
	limitNameFD,
	limitNameConns,
	limitNameConnsInbound,
	limitNameConnsOutbound,
	limitNameStreams,
	limitNameStreamsInbound,
	limitNameStreamsOutbound,
}

func resourceLimitsAndUsageToResourceInfo(scopeName string, stats ResourceLimitsAndUsage) ResourceInfos {
	result := ResourceInfos{}
	for _, l := range limits {
		ri := ResourceInfo{
			ScopeName: scopeName,
		}
		switch l {
		case limitNameMemory:
			ri.LimitName = limitNameMemory
			ri.LimitValue = stats.Memory
			ri.CurrentUsage = stats.MemoryUsage
		case limitNameFD:
			ri.LimitName = limitNameFD
			ri.LimitValue = rcmgr.LimitVal64(stats.FD)
			ri.CurrentUsage = int64(stats.FDUsage)
		case limitNameConns:
			ri.LimitName = limitNameConns
			ri.LimitValue = rcmgr.LimitVal64(stats.Conns)
			ri.CurrentUsage = int64(stats.ConnsUsage)
		case limitNameConnsInbound:
			ri.LimitName = limitNameConnsInbound
			ri.LimitValue = rcmgr.LimitVal64(stats.ConnsInbound)
			ri.CurrentUsage = int64(stats.ConnsInboundUsage)
		case limitNameConnsOutbound:
			ri.LimitName = limitNameConnsOutbound
			ri.LimitValue = rcmgr.LimitVal64(stats.ConnsOutbound)
			ri.CurrentUsage = int64(stats.ConnsOutboundUsage)
		case limitNameStreams:
			ri.LimitName = limitNameStreams
			ri.LimitValue = rcmgr.LimitVal64(stats.Streams)
			ri.CurrentUsage = int64(stats.StreamsUsage)
		case limitNameStreamsInbound:
			ri.LimitName = limitNameStreamsInbound
			ri.LimitValue = rcmgr.LimitVal64(stats.StreamsInbound)
			ri.CurrentUsage = int64(stats.StreamsInboundUsage)
		case limitNameStreamsOutbound:
			ri.LimitName = limitNameStreamsOutbound
			ri.LimitValue = rcmgr.LimitVal64(stats.StreamsOutbound)
			ri.CurrentUsage = int64(stats.StreamsOutboundUsage)
		}

		if ri.LimitValue == rcmgr.Unlimited64 || ri.LimitValue == rcmgr.DefaultLimit64 {
			// ignore unlimited and unset limits to remove noise from output.
			continue
		}

		result = append(result, ri)
	}

	return result
}
