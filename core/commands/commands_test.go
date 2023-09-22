package commands

import (
	"strings"
	"testing"

	cmds "github.com/bittorrent/go-btfs-cmds"
)

func collectPaths(prefix string, cmd *cmds.Command, out map[string]struct{}) {
	for name, sub := range cmd.Subcommands {
		path := prefix + "/" + name
		out[path] = struct{}{}
		collectPaths(path, sub, out)
	}
}

func TestROCommands(t *testing.T) {
	list := []string{
		"/block",
		"/block/get",
		"/block/stat",
		"/cat",
		"/commands",
		"/dag",
		"/dag/get",
		"/dag/resolve",
		"/dag/stat",
		"/dns",
		"/get",
		"/ls",
		"/name",
		"/name/resolve",
		"/object",
		"/object/data",
		"/object/get",
		"/object/links",
		"/object/stat",
		"/refs",
		"/resolve",
		"/version",
	}

	cmdSet := make(map[string]struct{})
	collectPaths("", RootRO, cmdSet)

	for _, path := range list {
		if _, ok := cmdSet[path]; !ok {
			t.Errorf("%q not in result", path)
		} else {
			delete(cmdSet, path)
		}
	}

	for path := range cmdSet {
		t.Errorf("%q in result but shouldn't be", path)
	}

	for _, path := range list {
		path = path[1:] // remove leading slash
		split := strings.Split(path, "/")
		sub, err := RootRO.Get(split)
		if err != nil {
			t.Errorf("error getting subcommand %q: %v", path, err)
		} else if sub == nil {
			t.Errorf("subcommand %q is nil even though there was no error", path)
		}
	}
}
func TestCommands(t *testing.T) {
	list := []string{
		"/add",
		"/bitswap",
		"/bitswap/ledger",
		"/bitswap/reprovide",
		"/bitswap/stat",
		"/bitswap/wantlist",
		"/block",
		"/block/get",
		"/block/put",
		"/block/rm",
		"/block/stat",
		"/bootstrap",
		"/bootstrap/add",
		"/bootstrap/add/default",
		"/bootstrap/list",
		"/bootstrap/rm",
		"/bootstrap/rm/all",
		"/cat",
		"/commands",
		"/config",
		"/config/edit",
		"/config/replace",
		"/config/reset",
		"/config/show",
		//"/config/profile",
		//"/config/profile/apply",
		"/config/storage-host-enable",
		"/config/sync-chain-info",
		"/config/sync-simple-mode",
		"/config/optin",
		"/config/optout",
		"/dag",
		"/dag/get",
		"/dag/export",
		"/dag/put",
		"/dag/import",
		"/dag/resolve",
		"/dag/stat",
		"/dht",
		"/dht/findpeer",
		"/dht/findprovs",
		"/dht/get",
		"/dht/provide",
		"/dht/put",
		"/dht/query",
		"/diag",
		"/diag/cmds",
		"/diag/cmds/clear",
		"/diag/cmds/set-time",
		"/diag/sys",
		"/dns",
		"/file",
		"/file/ls",
		"/files",
		"/files/chcid",
		"/files/cp",
		"/files/flush",
		"/files/ls",
		"/files/mkdir",
		"/files/mv",
		"/files/read",
		"/files/rm",
		"/files/stat",
		"/filestore",
		"/filestore/dups",
		"/filestore/ls",
		"/filestore/verify",
		"/files/write",
		"/get",
		"/id",
		"/key",
		"/key/gen",
		"/key/list",
		"/key/rename",
		"/key/rm",
		"/log",
		"/log/level",
		"/log/ls",
		"/log/tail",
		"/ls",
		"/mount",
		"/name",
		"/name/publish",
		"/name/pubsub",
		"/name/pubsub/state",
		"/name/pubsub/subs",
		"/name/pubsub/cancel",
		"/name/resolve",
		"/object",
		"/object/data",
		"/object/diff",
		"/object/get",
		"/object/links",
		"/object/new",
		"/object/patch",
		"/object/patch/add-link",
		"/object/patch/append-data",
		"/object/patch/rm-link",
		"/object/patch/set-data",
		"/object/put",
		"/object/stat",
		"/p2p",
		"/p2p/close",
		"/p2p/forward",
		"/p2p/listen",
		"/p2p/ls",
		"/p2p/stream",
		"/p2p/stream/close",
		"/p2p/stream/ls",
		"/pin",
		"/pin/add",
		"/ping",
		"/pin/ls",
		"/pin/rm",
		"/pin/update",
		"/pin/verify",
		"/pubsub",
		"/pubsub/ls",
		"/pubsub/peers",
		"/pubsub/pub",
		"/pubsub/sub",
		"/refs",
		"/refs/local",
		"/repo",
		"/repo/fsck",
		"/repo/gc",
		"/repo/stat",
		"/repo/verify",
		"/repo/version",
		"/resolve",
		"/rm",
		"/shutdown",
		"/restart",
		"/stats",
		"/stats/bitswap",
		"/stats/bw",
		"/stats/dht",
		"/stats/repo",
		"/swarm",
		"/swarm/addrs",
		"/swarm/addrs/listen",
		"/swarm/addrs/local",
		"/swarm/connect",
		"/swarm/disconnect",
		"/swarm/filters",
		"/swarm/filters/add",
		"/swarm/filters/rm",
		"/swarm/peers",
		"/tar",
		"/tar/add",
		"/tar/cat",
		"/urlstore",
		"/urlstore/add",
		"/version",
		"/version/deps",
		"/cid",
		"/cid/format",
		"/cid/base32",
		"/cid/codecs",
		"/cid/bases",
		"/cid/hashes",
		"/storage",
		"/storage/path",
		"/storage/path/capacity",
		"/storage/path/status",
		"/storage/path/migrate",
		"/storage/path/list",
		"/storage/path/mkdir",
		"/storage/path/volumes",
		"/storage/upload",
		"/storage/upload/init",
		"/storage/upload/recvcontract",
		"/storage/upload/status",
		"/storage/upload/repair",
		"/storage/upload/getcontractbatch",
		"/storage/upload/signcontractbatch",
		"/storage/upload/getunsigned",
		"/storage/upload/sign",
		"/storage/announce",
		"/storage/info",
		"/storage/hosts",
		"/storage/hosts/sync",
		"/storage/hosts/info",
		"/storage/challenge",
		"/storage/challenge/request",
		"/storage/challenge/response",
		"/storage/dcrepair",
		"/storage/dcrepair/request",
		"/storage/dcrepair/response",
		"/storage/stats",
		"/storage/stats/info",
		"/storage/stats/sync",
		"/storage/stats/list",
		"/storage/contracts",
		"/storage/contracts/list",
		"/storage/contracts/stat",
		"/storage/contracts/sync",
		"/metadata",
		"/metadata/add",
		"/metadata/rm",
		"/guard",
		"/guard/test",
		"/guard/test/send-challenges",
		"/cheque",
		"/cheque/stats",
		"/cheque/stats-all",
		"/cheque/send-history-stats",
		"/cheque/send-history-stats-all",
		"/cheque/cashlist",
		"/cheque/receive-history-stats",
		"/cheque/receive-history-stats-all",
		"/cheque/bttbalance",
		"/cheque/token_balance",
		"/cheque/all_token_balance",
		"/cheque/cash",
		"/cheque/cashstatus",
		"/cheque/chaininfo",
		"/cheque/price",
		"/cheque/price-all",
		"/cheque/receive",
		"/cheque/receive-history-list",
		"/cheque/receive-history-peer",
		"/cheque/receive-total-count",
		"/cheque/receivelist",
		"/cheque/receivelistall",
		"/cheque/send",
		"/cheque/send-history-list",
		"/cheque/send-history-peer",
		"/cheque/send-total-count",
		"/cheque/sendlist",
		"/cheque/sendlistall",
		"/p2p/handshake",
		"/settlement",
		"/settlement/list",
		"/settlement/peer",
		"/storage/upload/cheque",
		"/storage/upload/supporttokens",
		"/test",
		"/test/cheque",
		"/test/hosts",
		"/test/p2phandshake",
		"/vault",
		"/vault/address",
		"/vault/balance",
		"/vault/balance_all",
		"/vault/deposit",
		"/vault/wbttbalance",
		"/vault/withdraw",
		"/vault/upgrade",
		"/network",
		"/bttc",
		"/bttc/btt2wbtt",
		"/bttc/wbtt2btt",
		"/bttc/send-btt-to",
		"/bttc/send-wbtt-to",
		"/bttc/send-token-to",
		"/statuscontract",
		"/statuscontract/total",
		"/statuscontract/reportlist",
		"/statuscontract/lastinfo",
		"/statuscontract/config",
		"/statuscontract/report_online_server",
		//"/statuscontract/report_status_contract",
		"/statuscontract/daily_report_online_server",
		"/statuscontract/daily_report_list",
		"/statuscontract/daily_total",
		"/statuscontract/daily_last_report_time",
		"/bittorrent",
		"/bittorrent/download",
		"/bittorrent/serve",
		"/bittorrent/scrape",
		"/bittorrent/metainfo",
		"/bittorrent/bencode",
		"/accesskey",
		"/accesskey/generate",
		"/accesskey/enable",
		"/accesskey/disable",
		"/accesskey/reset",
		"/accesskey/delete",
		"/accesskey/get",
		"/accesskey/list",
	}

	cmdSet := make(map[string]struct{})
	collectPaths("", Root, cmdSet)

	for _, path := range list {
		if _, ok := cmdSet[path]; !ok {
			t.Errorf("%q not in result", path)
		} else {
			delete(cmdSet, path)
		}
	}

	for path := range cmdSet {
		t.Errorf("%q in result but shouldn't be", path)
	}

	for _, path := range list {
		path = path[1:] // remove leading slash
		split := strings.Split(path, "/")
		sub, err := Root.Get(split)
		if err != nil {
			t.Errorf("error getting subcommand %q: %v", path, err)
		} else if sub == nil {
			t.Errorf("subcommand %q is nil even though there was no error", path)
		}
	}
}
