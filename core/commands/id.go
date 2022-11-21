package commands

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	version "github.com/bittorrent/go-btfs"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	ke "github.com/bittorrent/go-btfs/core/commands/keyencode"

	kb "github.com/libp2p/go-libp2p-kbucket"
	ic "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	pstore "github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/tron-us/go-btfs-common/crypto"
)

const offlineIdErrorMessage = `'btfs id' currently cannot query information on remote
peers without a running daemon; we are working to fix this.
In the meantime, if you want to query remote peers using 'btfs id',
please run the daemon:

    btfs daemon &
    btfs id QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ
`

type IdOutput struct {
	ID              string
	PublicKey       string
	Addresses       []string
	AgentVersion    string
	ProtocolVersion string
	Protocols       []string
	DaemonProcessID int
	TronAddress     string
	BttcAddress     string
	VaultAddress    string
	ChainID         int64
}

const (
	formatOptionName   = "format"
	idFormatOptionName = "peerid-base"
)

var IDCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Show btfs node id info.",
		ShortDescription: `
Prints out information about the specified peer.
If no peer is specified, prints out information for local peers.

'btfs id' supports the format option for output with the following keys:
<id> : The peers id.
<aver>: Agent version.
<pver>: Protocol version.
<pubkey>: Public key.
<addrs>: Addresses (newline delimited).

EXAMPLE:

    btfs id Qmece2RkXhsKe5CRooNisBTh4SK119KrXXGmoK6V3kb8aH -f="<addrs>\n"
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("peerid", false, false, "Peer.ID of node to look up."),
	},
	Options: []cmds.Option{
		cmds.StringOption(formatOptionName, "f", "Optional output format."),
		cmds.StringOption(idFormatOptionName, "Encoding used for peer IDs: Can either be a multibase encoded CID or a base58btc encoded multihash. Takes {b58mh|base36|k|base32|b...}.").WithDefault("b58mh"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		keyEnc, err := ke.KeyEncoderFromString(req.Options[idFormatOptionName].(string))
		if err != nil {
			return err
		}

		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		var id peer.ID
		if len(req.Arguments) > 0 {
			var err error
			id, err = peer.Decode(req.Arguments[0])
			if err != nil {
				return fmt.Errorf("invalid peer id")
			}
		} else {
			id = n.Identity
		}

		if id == n.Identity {
			output, err := printSelf(keyEnc, n, env)
			if err != nil {
				return err
			}
			return cmds.EmitOnce(res, output)
		}

		// TODO handle offline mode with polymorphism instead of conditionals
		if !n.IsOnline {
			return errors.New(offlineIdErrorMessage)
		}

		// We need to actually connect to run identify.
		err = n.PeerHost.Connect(req.Context, peer.AddrInfo{ID: id})
		switch err {
		case nil:
		case kb.ErrLookupFailure:
			return errors.New(offlineIdErrorMessage)
		default:
			return err
		}

		output, err := printPeer(keyEnc, n.Peerstore, id, n, env)
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, output)
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *IdOutput) error {
			format, found := req.Options[formatOptionName].(string)
			if found {
				output := format
				output = strings.Replace(output, "<id>", out.ID, -1)
				output = strings.Replace(output, "<aver>", out.AgentVersion, -1)
				output = strings.Replace(output, "<pver>", out.ProtocolVersion, -1)
				output = strings.Replace(output, "<pubkey>", out.PublicKey, -1)
				output = strings.Replace(output, "<addrs>", strings.Join(out.Addresses, "\n"), -1)
				output = strings.Replace(output, "<protocols>", strings.Join(out.Protocols, "\n"), -1)
				output = strings.Replace(output, "\\n", "\n", -1)
				output = strings.Replace(output, "\\t", "\t", -1)
				fmt.Fprint(w, output)
			} else {
				marshaled, err := json.MarshalIndent(out, "", "\t")
				if err != nil {
					return err
				}
				marshaled = append(marshaled, byte('\n'))
				fmt.Fprintln(w, string(marshaled))
			}
			return nil
		}),
	},
	Type: IdOutput{},
}

func printPeer(keyEnc ke.KeyEncoder, ps pstore.Peerstore, p peer.ID, node *core.IpfsNode, env cmds.Environment) (interface{}, error) {
	if p == "" {
		return nil, errors.New("attempted to print nil peer")
	}

	info := new(IdOutput)
	info.ID = keyEnc.FormatID(p)

	if pk := ps.PubKey(p); pk != nil {
		pkb, err := ic.MarshalPublicKey(pk)
		if err != nil {
			return nil, err
		}
		info.PublicKey = base64.StdEncoding.EncodeToString(pkb)
	}

	addrInfo := ps.PeerInfo(p)
	addrs, err := peer.AddrInfoToP2pAddrs(&addrInfo)
	if err != nil {
		return nil, err
	}

	for _, a := range addrs {
		info.Addresses = append(info.Addresses, a.String())
	}
	sort.Strings(info.Addresses)

	protocols, _ := ps.GetProtocols(p) // don't care about errors here.
	for _, p := range protocols {
		info.Protocols = append(info.Protocols, string(p))
	}
	sort.Strings(info.Protocols)

	if v, err := ps.Get(p, "ProtocolVersion"); err == nil {
		if vs, ok := v.(string); ok {
			info.ProtocolVersion = vs
		}
	}
	if v, err := ps.Get(p, "AgentVersion"); err == nil {
		if vs, ok := v.(string); ok {
			info.AgentVersion = vs
		}
	}

	if node.IsDaemon {
		info.DaemonProcessID = os.Getpid()

		info.BttcAddress = chain.ChainObject.OverlayAddress.Hex()
		info.VaultAddress = chain.SettleObject.VaultService.Address().Hex()
	} else {
		info.DaemonProcessID = -1

		bttcAddr, err := chain.GetBttcNonDaemon(env)
		if err != nil {
			return nil, err
		}
		info.BttcAddress = bttcAddr

		valutAddr, err := chain.GetVaultNonDaemon(env)
		if err != nil {
			return nil, err
		}
		info.VaultAddress = valutAddr
	}

	return info, nil
}

// printing self is special cased as we get values differently.
func printSelf(keyEnc ke.KeyEncoder, node *core.IpfsNode, env cmds.Environment) (interface{}, error) {
	info := new(IdOutput)
	info.ID = keyEnc.FormatID(node.Identity)

	pk := node.PrivateKey.GetPublic()
	pkb, err := ic.MarshalPublicKey(pk)
	if err != nil {
		return nil, err
	}
	info.PublicKey = base64.StdEncoding.EncodeToString(pkb)

	if node.PeerHost != nil {
		addrs, err := peer.AddrInfoToP2pAddrs(host.InfoFromHost(node.PeerHost))
		if err != nil {
			return nil, err
		}
		for _, a := range addrs {
			info.Addresses = append(info.Addresses, a.String())
		}
		sort.Strings(info.Addresses)
		info.Protocols = node.PeerHost.Mux().Protocols()
		sort.Strings(info.Protocols)
	}
	info.ProtocolVersion = "btfs/0.1.0" //identify.LibP2PVersion
	info.AgentVersion = version.UserAgent
	raw, err := node.PrivateKey.Raw()
	if err != nil {
		return nil, err
	}
	keys, err := crypto.FromPrivateKey(hex.EncodeToString(raw))
	if err != nil {
		return nil, err
	}
	info.TronAddress = keys.Base58Address

	if node.IsDaemon {
		info.DaemonProcessID = os.Getpid()

		info.BttcAddress = chain.ChainObject.OverlayAddress.Hex()
		info.VaultAddress = chain.SettleObject.VaultService.Address().Hex()

		// show chain id only local peer and in daemon mode
		info.ChainID = chain.ChainObject.ChainID
	} else {
		info.DaemonProcessID = -1

		bttcAddr, err := chain.GetBttcNonDaemon(env)
		if err != nil {
			return nil, err
		}
		info.BttcAddress = bttcAddr

		valutAddr, err := chain.GetVaultNonDaemon(env)
		if err != nil {
			return nil, err
		}
		info.VaultAddress = valutAddr
	}

	return info, nil
}
