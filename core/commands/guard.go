package commands

import (
	"fmt"
	"github.com/bittorrent/go-btfs/utils"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/guard"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/core/hub"

	cmds "github.com/bittorrent/go-btfs-cmds"
	cconfig "github.com/bittorrent/go-btfs-common/config"

	"github.com/ipfs/go-cid"
)

const (
	guardUrlOptionName                   = "url"
	guardQuestionCountPerShardOptionName = "questions-per-shard"
	guardHostsOptionName                 = "hosts"
)

var GuardCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Interact with guard services from BTFS client.",
		ShortDescription: `
Connect with guard functions directly through this command.
The subcommands here are mostly for debugging and testing purposes.`,
	},
	Subcommands: map[string]*cmds.Command{
		"test": guardTestCmd,
	},
}

var guardTestCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Send tests to guard service endpoints from BTFS client.",
		ShortDescription: `
This command contains subcommands that are typically for development purposes
by letting the BTFS client test individual guard endpoints.`,
	},
	Subcommands: map[string]*cmds.Command{
		"send-challenges": guardSendChallengesCmd,
	},
}

var guardSendChallengesCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Send shard challenge questions from BTFS client.",
		ShortDescription: `
Sends all shard challenge questions under a reed-solomon encoded file
to the guard service.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("file-hash", true, false, "File hash to generate the questions from."),
		cmds.StringArg("session-id", true, false, "Session id to generate the questions from."),
	},
	Options: []cmds.Option{
		cmds.StringOption(guardUrlOptionName, "u", "Guard service url including protocol and port. Default: reads from BTFS config."),
		cmds.IntOption(guardQuestionCountPerShardOptionName, "q", "Number of challenge questions per shard to generate").WithDefault(cconfig.ConstMinQuestionsCountPerChallenge),
		cmds.StringOption(guardHostsOptionName, "sh", "List of hosts for each shard, ordered sequentially and separated by ','. Default: reads from BTFS datastore."),
	},
	RunTimeout: 30 * time.Second,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}
		fmt.Println("CheckSimpleMode ", err)

		// get config settings
		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}
		if !cfg.Experimental.StorageClientEnabled {
			return fmt.Errorf("storage client api not enabled")
		}
		// get node
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		// get core api
		api, err := cmdenv.GetApi(env, req)
		if err != nil {
			return err
		}
		rootHash, err := cid.Parse(req.Arguments[0])
		if err != nil {
			return err
		}
		shardHashes, _, err := helper.CheckAndGetReedSolomonShardHashes(req.Context, n, api, rootHash)
		if err != nil {
			return err
		}
		var hostIDs []string
		if hl, found := req.Options[guardHostsOptionName].(string); found {
			hostIDs = strings.Split(hl, ",")
		} else {
			hosts, err := helper.GetHostsFromDatastore(req.Context, n, hub.HubModeAll, len(shardHashes))
			if err != nil {
				return err
			}
			for _, ni := range hosts {
				hostIDs = append(hostIDs, ni.NodeId)
			}
		}

		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}
		qCount, _ := req.Options[guardQuestionCountPerShardOptionName].(int)
		hashStrs := make([]string, 0)
		for _, c := range shardHashes {
			hashStrs = append(hashStrs, c.String())
		}
		rss, err := sessions.GetRenterSession(ctxParams, req.Arguments[1], req.Arguments[0], hashStrs)
		if err != nil {
			return err
		}
		questions, err := guard.PrepCustomFileChallengeQuestions(rss, rootHash, shardHashes, hostIDs, qCount, false,
			n.Identity.Pretty())
		if err != nil {
			return err
		}
		// check if we need to update config for a different guard url
		if gu, found := req.Options[guardUrlOptionName].(string); found {
			cfg, err = cfg.Clone()
			if err != nil {
				return err
			}
			cfg.Services.GuardDomain = gu
		}
		// send to guard
		return guard.SendChallengeQuestions(req.Context, cfg, rootHash, questions)
	},
}
