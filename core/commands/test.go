package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"io"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	hubpb "github.com/bittorrent/go-btfs-common/protos/hub"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/corehttp/remote"
	"github.com/bittorrent/go-btfs/core/hub"
	"github.com/bittorrent/go-btfs/settlement/swap/swapprotocol"
	"github.com/bittorrent/go-btfs/settlement/swap/swapprotocol/pb"

	"github.com/ethereum/go-ethereum/common"
	peerInfo "github.com/libp2p/go-libp2p/core/peer"
)

type TestOutput struct {
	Status string
}

var testOptionDesc = "test it, get hosts or send cheque."

var TestCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "test it.",
		ShortDescription: `test it. get hosts or send cheque.`,
	},

	Run:      testCmd.Run,
	Encoders: testCmd.Encoders,
	Type:     testCmd.Type,

	Subcommands: map[string]*cmds.Command{
		"hosts":        testHostsCmd,
		"cheque":       testChequeCmd,
		"p2phandshake": testP2pShakeCmd,
	},
}

var testCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Show peers in the bootstrap list.",
		ShortDescription: "Peers are output in the format '<multiaddr>/<peerID>'.",
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		return cmds.EmitOnce(res, &TestOutput{"test ok"})
	},
	Type: TestOutput{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *TestOutput) error {
			return testWrite(w, "", out.Status)
		}),
	},
}

var testHostsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Show peers in the bootstrap list.",
		ShortDescription: "Peers are output in the format '<multiaddr>/<peerID>'.",
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nodes, err := getHosts(req, env)
		if err != nil {
			return err
		}

		fmt.Println("len(nodes) = ", len(nodes))
		for i := range nodes {
			fmt.Printf("nodes[%d] = %+v", i, nodes[i])
			if i >= 5 {
				break
			}
		}

		return cmds.EmitOnce(res, &TestOutput{"get hosts ok"})
	},
	Type: TestOutput{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *TestOutput) error {
			return testWrite(w, "", out.Status)
		}),
	},
}

func getHosts(req *cmds.Request, env cmds.Environment) ([]*hubpb.Host, error) {
	cfg, err := cmdenv.GetConfig(env)
	if err != nil {
		return nil, err
	}
	if !cfg.Experimental.StorageClientEnabled {
		return nil, fmt.Errorf("storage client api not enabled")
	}

	mode, ok := req.Options["host-sync-mode"].(string)
	if !ok {
		mode = cfg.Experimental.HostsSyncMode
	}

	node, err := cmdenv.GetNode(env)
	if err != nil {
		return nil, err
	}

	nodes, err := hub.QueryHosts(req.Context, node, mode)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

var testChequeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Show peers in the bootstrap list.",
		ShortDescription: "Peers are output in the format '<multiaddr>/<peerID>'.",
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nodes, err := getHosts(req, env)
		if err != nil {
			return err
		}
		//fmt.Println("get hosts: ", nodes)

		if len(nodes) <= 0 {
			return errors.New("get hosts, it is none")
		}

		toPeer := "16Uiu2HAmQ4VAZdAmXJKapZYAib2BZpsSkVQfTwowCYvSqbEyuTnQ"
		//toPeer := nodes[0].NodeId
		chain.SettleObject.SwapService.Settle(toPeer, big.NewInt(10), "", tokencfg.GetWbttToken())

		return cmds.EmitOnce(res, &TestOutput{"send cheque ok"})
	},
	Type: TestOutput{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *TestOutput) error {
			return testWrite(w, "", out.Status)
		}),
	},
}

var testP2pShakeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "p2p handshake.",
		ShortDescription: "p2p handshake.",
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nodes, err := getHosts(req, env)
		if err != nil {
			return err
		}
		//fmt.Println("get hosts: ", nodes)

		if len(nodes) <= 0 {
			return errors.New("get hosts, it is none")
		}

		peer := "16Uiu2HAm5TZ8547Xzbqoynt4cXGg6mjkCGoLBWnNxSXcoW7eTBdW"
		//peer := nodes[0].NodeId

		peerhostPid, err := peerInfo.Decode(peer)
		if err != nil {
			log.Infof("peer.IDB58Decode(peer:%s) error: %s", peer, err)
			return err
		}

		ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
		ctxParams, err := helper.ExtractContextParams(req, env)
		if err != nil {
			return err
		}

		//get beneficiary
		output, err := remote.P2PCall(ctx, ctxParams.N, ctxParams.Api, peerhostPid, "/p2p/handshake",
			swapprotocol.SwapProtocol.GetChainID(),
			ctxParams.N.Identity,
		)

		beneficiary := &pb.Handshake{}
		err = json.Unmarshal(output, beneficiary)
		if err != nil {
			return err
		}

		fmt.Println("remote.P2PCall, beneficiary = ", common.BytesToAddress(beneficiary.Beneficiary))

		return cmds.EmitOnce(res, &TestOutput{"p2p handshake ok"})
	},
	Type: TestOutput{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *TestOutput) error {
			return testWrite(w, "", out.Status)
		}),
	},
}

func testWrite(w io.Writer, prefix string, status string) error {
	_, err := w.Write([]byte(prefix + status + "\n"))
	if err != nil {
		return err
	}
	return nil
}
