package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/assets"
	"github.com/bittorrent/go-btfs/chain"
	chaincfg "github.com/bittorrent/go-btfs/chain/config"
	"github.com/bittorrent/go-btfs/cmd/btfs/util"
	oldcmds "github.com/bittorrent/go-btfs/commands"
	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands"
	"github.com/bittorrent/go-btfs/namesys"
	fsrepo "github.com/bittorrent/go-btfs/repo/fsrepo"

	cmds "github.com/bittorrent/go-btfs-cmds"
	config "github.com/bittorrent/go-btfs-config"
	files "github.com/bittorrent/go-btfs-files"
)

const (
	bitsOptionName      = "bits"
	emptyRepoOptionName = "empty-repo"
	profileOptionName   = "profile"
	keyTypeDefault      = "BIP39"
	keyTypeOptionName   = "key"
	importKeyOptionName = "import"
	rmOnUnpinOptionName = "rm-on-unpin"
	seedOptionName      = "seed"
	simpleMode          = "simple-mode"
	recoveryOptionName  = "recovery"
	/*
		passWordOptionName     = "password"
		passwordFileoptionName = "password-file"
	*/
)

var initCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Initializes btfs config file.",
		ShortDescription: `
Initializes btfs configuration files and generates a new keypair.

If you are going to run BTFS in server environment, you may want to
initialize it using 'server' profile.

For the list of available profiles see 'btfs config profile --help'

btfs uses a repository in the local file system. By default, the repo is
located at ~/.btfs. To change the repo location, set the $BTFS_PATH
environment variable:

    export BTFS_PATH=/path/to/btfsrepo
`,
	},
	Arguments: []cmds.Argument{
		cmds.FileArg("default-config", false, false, "Initialize with the given configuration.").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.IntOption(bitsOptionName, "b", "Number of bits to use in the generated RSA private key.").WithDefault(util.NBitsForKeypairDefault),
		cmds.BoolOption(emptyRepoOptionName, "e", "Don't add and pin help files to the local storage."),
		cmds.StringOption(profileOptionName, "p", "Apply profile settings to config. Multiple profiles can be separated by ','"),
		cmds.StringOption(keyTypeOptionName, "k", "Key generation algorithm, e.g. RSA, Ed25519, Secp256k1, ECDSA, BIP39. By default is BIP39"),
		cmds.StringOption(importKeyOptionName, "i", "Import TRON private key to generate btfs PeerID."),
		cmds.BoolOption(rmOnUnpinOptionName, "r", "Remove unpinned files.").WithDefault(false),
		cmds.StringOption(seedOptionName, "s", "Import seed phrase"),
		cmds.BoolOption(simpleMode, "sm", "init with simple mode or not."),
		cmds.StringOption(recoveryOptionName, "Recovery data from a backup"),
		/*
			cmds.StringOption(passWordOptionName, "", "password for decrypting keys."),
			cmds.StringOption(passwordFileoptionName, "", "path to a file that contains password for decrypting keys"),
		*/

		// TODO need to decide whether to expose the override as a file or a
		// directory. That is: should we allow the user to also specify the
		// name of the file?
		// TODO cmds.StringOption("event-logs", "l", "Location for machine-readable event logs."),
	},
	NoRemote: true,
	Extra:    commands.CreateCmdExtras(commands.SetDoesNotUseRepo(true), commands.SetDoesNotUseConfigAsInput(true)),
	PreRun: func(req *cmds.Request, env cmds.Environment) error {
		cctx := env.(*oldcmds.Context)
		daemonLocked, err := fsrepo.LockedByOtherProcess(cctx.ConfigRoot)
		if err != nil {
			fmt.Println(`What causes this error: there is already one daemon process running in background
			Solution: kill it first and run btfs daemon again.
			If the user has the need to start multiple nodes on the same machine, the configuration needs to be modified.`)
			return err
		}

		log.Info("checking if daemon is running...")
		if daemonLocked {
			log.Debug("btfs daemon is running")
			e := "btfs daemon is running. please stop it to run this command"
			return cmds.ClientError(e)
		}

		return nil
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		cctx := env.(*oldcmds.Context)
		empty, _ := req.Options[emptyRepoOptionName].(bool)
		nBitsForKeypair, _ := req.Options[bitsOptionName].(int)
		rmOnUnpin, _ := req.Options[rmOnUnpinOptionName].(bool)

		var conf *config.Config

		f := req.Files
		if f != nil {
			it := req.Files.Entries()
			if !it.Next() {
				if it.Err() != nil {
					return it.Err()
				}
				return fmt.Errorf("file argument was nil")
			}
			file := files.FileFromEntry(it)
			if file == nil {
				return fmt.Errorf("expected a regular file")
			}

			conf = &config.Config{}
			if err := json.NewDecoder(file).Decode(conf); err != nil {
				return err
			}
		}

		profile, _ := req.Options[profileOptionName].(string)
		importKey, _ := req.Options[importKeyOptionName].(string)
		keyType, _ := req.Options[keyTypeOptionName].(string)
		seedPhrase, _ := req.Options[seedOptionName].(string)
		simpleModeIn, _ := req.Options[simpleMode].(bool)
		/*
			password, _ := req.Options[passWordOptionName].(string)
			passwordFile, _ := req.Options[passwordFileoptionName].(string)
		*/
		backupPath, ok := req.Options[recoveryOptionName].(string)
		if ok {
			btfsPath := env.(*oldcmds.Context).ConfigRoot
			dstPath := filepath.Dir(btfsPath)
			if fsrepo.IsInitialized(btfsPath) {
				newPath := filepath.Join(dstPath, fmt.Sprintf(".btfs_backup_%d", time.Now().Unix()))
				// newPath := filepath.Join(filepath.Dir(btfsPath), backup)
				err := os.Rename(btfsPath, newPath)
				if err != nil {
					return err
				}
				fmt.Println("btfs configuration file already exists!")
				fmt.Println("We have renamed it to ", newPath)
			}

			if err := commands.UnTar(backupPath, dstPath); err != nil {
				err = commands.UnZip(backupPath, dstPath)
				if err != nil {
					return errors.New("your file format is not tar.gz or zip, please check again")
				}
			}
			fmt.Println("Recovery successful!")
			return nil
		}
		return doInit(os.Stdout, cctx.ConfigRoot, empty, nBitsForKeypair, profile, conf, keyType, importKey, seedPhrase, rmOnUnpin, simpleModeIn)
	},
}

var errRepoExists = errors.New(`btfs configuration file already exists!
Reinitializing would overwrite your keys.
`)

func doInit(out io.Writer, repoRoot string, empty bool, nBitsForKeypair int, confProfiles string, conf *config.Config,
	keyType string, importKey string, mnemonic string, rmOnUnpin bool, simpleModeIn bool) error {

	importKey, mnemonic, err := util.GenerateKey(importKey, keyType, mnemonic)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "initializing BTFS node at %s\n", repoRoot); err != nil {
		return err
	}

	if err := checkWritable(repoRoot); err != nil {
		return err
	}

	if fsrepo.IsInitialized(repoRoot) {
		return errRepoExists
	}

	if conf == nil {
		var err error
		conf, err = config.Init(out, nBitsForKeypair, keyType, importKey, mnemonic, rmOnUnpin)
		if err != nil {
			return err
		}
		if rmOnUnpin {
			raw := json.RawMessage(`{"rmOnUnpin":"` + strconv.FormatBool(rmOnUnpin) + `"}`)
			conf.Datastore.Params = &raw
		}
	}

	if err := applyProfiles(conf, confProfiles); err != nil {
		return err
	}

	if err := addChainInfo(conf); err != nil {
		return err
	}

	if err := addIdentityInfo(conf, importKey); err != nil {
		return err
	}

	if err := storeChainId(conf, repoRoot); err != nil {
		return err
	}

	conf.SimpleMode = simpleModeIn

	if err := fsrepo.Init(repoRoot, conf); err != nil {
		return err
	}

	if !empty {
		if err := addDefaultAssets(out, repoRoot); err != nil {
			return err
		}
	}

	return initializeIpnsKeyspace(repoRoot)
}

// add chain id into leveldb
// btfs init cmd, not node process
func storeChainId(conf *config.Config, repoRoot string) error {
	statestore, err := chain.InitStateStore(repoRoot)
	if err != nil {
		fmt.Println("init statestore err: ", err)
		return err
	}

	defer statestore.Close()

	err = chain.StoreChainIdToDisk(conf.ChainInfo.ChainId, statestore)
	if err != nil {
		fmt.Println("init StoreChainId err: ", err)
		return err
	}

	return nil
}

// add chain info
func addChainInfo(conf *config.Config) error {
	chainId := conf.ChainInfo.ChainId
	chainCfg, found := chaincfg.GetChainConfig(chainId)
	if !found {
		return errors.New(fmt.Sprintf("chainid=%d is not found.", chainId))
	}

	conf.ChainInfo.CurrentFactory = chainCfg.CurrentFactory.Hex()
	// conf.ChainInfo.PriceOracleAddress = chainCfg.PriceOracleAddress.Hex()
	conf.ChainInfo.Endpoint = chainCfg.Endpoint
	return nil
}

// add Identity info
func addIdentityInfo(conf *config.Config, importKey string) error {
	conf.Identity.HexPrivKey = importKey

	bttcAddr, err := chain.GetBttcByKey(conf.Identity.PrivKey)
	if err != nil {
		return err
	}
	conf.Identity.BttcAddr = bttcAddr
	return nil
}

func applyProfiles(conf *config.Config, profiles string) error {
	if profiles == "" {
		return nil
	}
	for _, profile := range strings.Split(profiles, ",") {
		transformer, ok := config.Profiles[profile]
		if !ok {
			return fmt.Errorf("invalid configuration profile: %s", profile)
		}
		if err := transformer.Transform(conf); err != nil {
			return err
		}
	}
	return nil
}

func checkWritable(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		// dir exists, make sure we can write to it
		testfile := filepath.Join(dir, "test")
		fi, err := os.Create(testfile)
		if err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("%s is not writeable by the current user", dir)
			}
			return fmt.Errorf("unexpected error while checking writeablility of repo root: %s", err)
		}
		fi.Close()
		return os.Remove(testfile)
	}

	if os.IsNotExist(err) {
		// dir doesn't exist, check that we can create it
		return os.MkdirAll(dir, 0775)
	}

	if os.IsPermission(err) {
		return fmt.Errorf("cannot write to %s, incorrect permissions", err)
	}

	return err
}

func addDefaultAssets(out io.Writer, repoRoot string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := fsrepo.Open(repoRoot)
	if err != nil { // NB: repo is owned by the node
		return err
	}

	buildCfg := &core.BuildCfg{Repo: r}
	nd, err := core.NewNode(ctx, buildCfg)
	if err != nil {
		return err
	}
	defer nd.Close()

	dkey, err := assets.SeedInitDocs(nd)
	if err != nil {
		return fmt.Errorf("init: seeding init docs failed: %s", err)
	}
	log.Debugf("init: seeded init docs %s", dkey)

	if _, err = fmt.Fprintf(out, "to get started, enter:\n"); err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "\n\tbtfs cat /btfs/%s/readme\n\n", dkey)
	return err
}

func initializeIpnsKeyspace(repoRoot string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := fsrepo.Open(repoRoot)
	if err != nil { // NB: repo is owned by the node
		return err
	}

	nd, err := core.NewNode(ctx, &core.BuildCfg{Repo: r})
	if err != nil {
		return err
	}
	defer nd.Close()

	return namesys.InitializeKeyspace(ctx, nd.Namesys, nd.Pinning, nd.PrivateKey)
}
