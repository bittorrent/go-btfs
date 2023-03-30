package commands

import (
	"errors"

	"github.com/TRON-US/interface-go-btfs-core/options"
	"github.com/TRON-US/interface-go-btfs-core/path"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
)

type MetaResult struct {
	Hash string
}

const (
	metaOverwriteOptionName = "overwrite"
	metaPinOptionName       = "pin"
)

// MetadataCmd is the 'btfs metadata' command
var MetadataCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Interact with metadata for BTFS files.",
		ShortDescription: `
'btfs metadata' is a command to manipulate token metadata for BTFS files
 that are stored through BTT payment.`,
	},

	Subcommands: map[string]*cmds.Command{
		"add": metadataAddCmd,
		"rm":  metadataRemoveCmd,
	},
}

var metadataAddCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Add token metadata to a BTFS file.",
		ShortDescription: `
'btfs metadata add' adds token metadata item(s) to a BTFS file that is
        stored on the BTFS network through BTT payment. 
        We specify the target BTFS file hash and metadata items key-value pair in JSON string format.

Example:

To add metadata for a file, specify the file hash and metadata key-value pair:

        $btfs metadata add <file-hash> '{"price":11.2}'
        
This command returns a new file-hash for the file.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("file-hash", true, false, "BTFS target file hash."),
		cmds.StringArg("metadata", true, false, "Token metadata to append in JSON string."),
	},
	Options: []cmds.Option{
		cmds.BoolOption(metaPinOptionName, "Pin this object when adding.").WithDefault(true),
		cmds.BoolOption(metaOverwriteOptionName, "Overwrite metadata when there are existing key-value pairs.").WithDefault(false),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		api, err := cmdenv.GetApi(env, req)
		if err != nil {
			return err
		}

		pin, _ := req.Options[metaPinOptionName].(bool)

		enc, err := cmdenv.GetCidEncoder(req)
		if err != nil {
			return err
		}
		opts := []options.UnixfsAddMetaOption{
			options.Unixfs.OverwriteToAdd(pin),
			options.Unixfs.PinToAdd(pin),
		}
		fileHash := req.Arguments[0]
		tokenMetadata := req.Arguments[1]
		// TODO: use for loop or batch for token metadata items.

		p, err := api.Unixfs().AddMetadata(req.Context, path.New(fileHash), tokenMetadata, opts...)
		if err != nil {
			return err
		}
		h := ""
		if p != nil {
			h = enc.Encode(p.Cid())
		} else {
			return errors.New("got nil path")
		}

		err = res.Emit(&MetaResult{
			Hash: h,
		})
		if err != nil {
			return err
		}

		return nil
	},
	Type: MetaResult{},
}

var metadataRemoveCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Remove token metadata on a BTFS file.",
		ShortDescription: `
'btfs metadata rm' removes specified token metadata on a BTFS file that is
        stored on the BTFS network through BTT payment. 
        We specify the target BTFS file hash and the metadata item keys.

Example:
        
To remove the metadata for file, specify the file hash and metadata key:
	
       	$btfs metadata rm <file-hash> 'price'
        
The output returns a new file-hash for the file.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("file-hash", true, false, "BTFS target file hash."),
		cmds.StringArg("metadata", true, false, "Token metadata keys to remove."),
	},
	Options: []cmds.Option{
		cmds.BoolOption(metaPinOptionName, "Pin this object when removing.").WithDefault(true),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		api, err := cmdenv.GetApi(env, req)
		if err != nil {
			return err
		}

		pin, _ := req.Options[metaPinOptionName].(bool)

		enc, err := cmdenv.GetCidEncoder(req)
		if err != nil {
			return err
		}
		opts := []options.UnixfsRemoveMetaOption{
			options.Unixfs.PinToRemove(pin),
		}
		fileHash := req.Arguments[0]
		tokenMetadata := req.Arguments[1]
		// TODO: use for loop or batch for token metadata items.

		p, err := api.Unixfs().RemoveMetadata(req.Context, path.New(fileHash), tokenMetadata, opts...)
		if err != nil {
			return err
		}
		h := ""
		if p != nil {
			h = enc.Encode(p.Cid())
		} else {
			return errors.New("got nil path")
		}

		err = res.Emit(&MetaResult{
			Hash: h,
		})
		if err != nil {
			return err
		}

		return nil
	},
	Type: MetaResult{},
}
