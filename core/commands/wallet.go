package commands

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bittorrent/go-btfs/cmd/btfs/util"
	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/path"
	"github.com/bittorrent/go-btfs/core/wallet"
	walletpb "github.com/bittorrent/go-btfs/protos/wallet"

	"github.com/TRON-US/go-btfs-config"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs-cmds/http"
	"github.com/tron-us/go-btfs-common/crypto"

	"github.com/libp2p/go-libp2p-core/peer"
)

const (
	Domain = "wallet"
)

func init() {
	http.AddDomainWhiteList(Domain)
	for _, c := range []*cmds.Command{
		walletInitCmd,
		walletDepositCmd,
		walletWithdrawCmd,
		walletPasswordCmd,
		walletKeysCmd,
		walletImportCmd,
		walletTransferCmd,
		walletBalanceCmd,
		walletDiscoveryCmd,
		walletGenerateKeyCmd,
	} {
		c.Extra = c.Extra.SetValue(http.DomainKey, Domain)
	}
}

var WalletCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "BTFS wallet",
		ShortDescription: `'btfs wallet' is a set of commands to interact with block chain and ledger.`,
		LongDescription: `'btfs wallet' is a set of commands interact with block chain and ledger to deposit,
withdraw and query balance of token used in BTFS.`,
	},

	Subcommands: map[string]*cmds.Command{
		"init":              walletInitCmd,
		"deposit":           walletDepositCmd,
		"withdraw":          walletWithdrawCmd,
		"balance":           walletBalanceCmd,
		"password":          walletPasswordCmd,
		"keys":              walletKeysCmd,
		"transactions":      walletTransactionsCmd,
		"import":            walletImportCmd,
		"transfer":          walletTransferCmd,
		"discovery":         walletDiscoveryCmd,
		"validate_password": walletCheckPasswordCmd,
		"generate_key":      walletGenerateKeyCmd,
	},
}

var walletInitCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "initialize BTFS wallet",
		ShortDescription: "initialize BTFS wallet",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("private_key", true, false, "private key"),
		cmds.StringArg("encrypted_private_key", true, false, "encrypted private key"),
		cmds.StringArg("encrypted_mnemonic", true, false, "encrypted mnemonic"),
	},
	Options: []cmds.Option{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}
		cfg.Identity.PrivKey = req.Arguments[0]
		cfg.Identity.Mnemonic = ""
		cfg.Identity.EncryptedPrivKey = req.Arguments[1]
		cfg.Identity.EncryptedMnemonic = req.Arguments[2]
		ks, err := crypto.ToPrivKey(req.Arguments[0])
		if err != nil {
			return err
		}
		id, err := peer.IDFromPrivateKey(ks)
		if err != nil {
			return err
		}
		cfg.Identity.PeerID = id.String()
		cfg.UI.Wallet.Initialized = true
		if err = n.Repo.SetConfig(cfg); err != nil {
			return err
		}
		go path.DoRestart(false)
		return nil
	},
	NoRemote: true,
}

var walletGenerateKeyCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Generate new private_key and Mnemonic",
		ShortDescription: "Generate new private_key and Mnemonic",
	},

	Arguments: []cmds.Argument{},
	Options:   []cmds.Option{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		k, m, err := util.GenerateKey("", "BIP39", "")
		if err != nil {
			return err
		}
		ks, err := crypto.FromPrivateKey(k)
		if err != nil {
			return err
		}
		k64, err := crypto.Hex64ToBase64(ks.HexPrivateKey)
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, Keys{
			PrivateKey: k64,
			Mnemonic:   m,
			SkInBase64: k64,
			SkInHex:    ks.HexPrivateKey,
		})
	},
	Type:     Keys{},
	NoRemote: true,
}

const (
	asyncOptionName    = "async"
	passwordOptionName = "password"
)

var walletDepositCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "BTFS wallet deposit",
		ShortDescription: "BTFS wallet deposit from block chain to ledger. Use '-p=<password>' to specific password.",
		Options:          "unit is µBTT (=0.000001BTT)",
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("amount", true, false, "amount to deposit."),
	},
	Options: []cmds.Option{
		cmds.BoolOption(asyncOptionName, "a", "Deposit asynchronously."),
		cmds.StringOption(passwordOptionName, "p", "password"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}
		if err := checkWhetherPasswordSet(cfg); err != nil {
			return err
		}
		password, err := getPassword(req)
		if err != nil {
			return err
		}
		if err := validatePassword(cfg, password); err != nil {
			return err
		}
		async, _ := req.Options[asyncOptionName].(bool)

		amount, err := strconv.ParseInt(req.Arguments[0], 10, 64)
		if err != nil {
			return err
		}

		runDaemon := false

		currentNode, err := cmdenv.GetNode(env)
		if err != nil {
			log.Error("Wrong while get current Node information", err)
			return err
		}
		runDaemon = currentNode.IsDaemon

		err = wallet.WalletDeposit(req.Context, cfg, n, amount, runDaemon, async)
		if err != nil {
			if strings.Contains(err.Error(), "Please deposit at least") {
				err = errors.New("Please deposit at least 10,000,000µBTT(=10BTT)")
			}
			return err
		}
		s := fmt.Sprintf("BTFS wallet deposit submitted. Please wait one minute for the transaction to confirm.")
		if !runDaemon {
			s = fmt.Sprintf("BTFS wallet deposit Done.")
		}
		return cmds.EmitOnce(res, &MessageOutput{s})
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *MessageOutput) error {
			fmt.Fprint(w, out.Message)
			return nil
		}),
	},
	Type:     MessageOutput{},
	NoRemote: true,
}

var walletWithdrawCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "BTFS wallet withdraw",
		ShortDescription: "BTFS wallet withdraw from ledger to block chain. Use '-p=<password>' to specific password.",
		Options:          "unit is µBTT (=0.000001BTT)",
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("amount", true, false, "amount to deposit."),
	},
	Options: []cmds.Option{
		cmds.StringOption(passwordOptionName, "p", "password"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}
		if err := checkWhetherPasswordSet(cfg); err != nil {
			return err
		}
		password, err := getPassword(req)
		if err != nil {
			return err
		}
		if err := validatePassword(cfg, password); err != nil {
			return err
		}
		amount, err := strconv.ParseInt(req.Arguments[0], 10, 64)
		if err != nil {
			return err
		}

		err = wallet.WalletWithdraw(req.Context, cfg, n, amount)
		if err != nil {
			if strings.Contains(err.Error(), "Please withdraw at least") {
				err = errors.New("Please withdraw at least 1,000,000,000µBTT(=1000BTT)")
			}
			return err
		}

		s := fmt.Sprintf("BTFS wallet withdraw submitted. Please wait one minute for the transaction to confirm.")
		return cmds.EmitOnce(res, &MessageOutput{s})
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *MessageOutput) error {
			fmt.Fprint(w, out.Message)
			return nil
		}),
	},
	Type:     MessageOutput{},
	NoRemote: true,
}

var walletBalanceCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "BTFS wallet balance",
		ShortDescription: "Query BTFS wallet balance in ledger and block chain.",
		Options:          "unit is µBTT (=0.000001BTT)",
	},

	Arguments: []cmds.Argument{},
	Options:   []cmds.Option{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}

		tronBalance, ledgerBalance, err := wallet.GetBalance(req.Context, cfg)
		if err != nil {
			log.Error("wallet get balance failed, ERR: ", err)
			return err
		}
		s := fmt.Sprintf("BTFS wallet tron balance '%d', ledger balance '%d'\n", tronBalance, ledgerBalance)
		log.Info(s)

		return cmds.EmitOnce(res, &BalanceResponse{
			BtfsWalletBalance: uint64(ledgerBalance),
			BttWalletBalance:  uint64(tronBalance),
		})
	},
	Type:     BalanceResponse{},
	NoRemote: true,
}

type BalanceResponse struct {
	BtfsWalletBalance uint64
	BttWalletBalance  uint64
}

var walletPasswordCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "BTFS wallet password",
		ShortDescription: "set password for BTFS wallet",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("password", true, false, "new password of BTFS wallet."),
	},
	Options: []cmds.Option{
		cmds.StringOption(passwordOptionName, "p", "old password").WithDefault(""),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}
		if err := checkWhetherPasswordSet(cfg); err == nil {
			pw, err := getPassword(req)
			if err != nil {
				return err
			}
			if err := validatePassword(cfg, pw); err != nil {
				return err
			}
		}
		cipherMnemonic, err := wallet.EncryptWithAES(req.Arguments[0], cfg.Identity.Mnemonic)
		if err != nil {
			return err
		}
		cipherPrivKey, err := wallet.EncryptWithAES(req.Arguments[0], cfg.Identity.PrivKey)
		if err != nil {
			return err
		}
		cfg.Identity.EncryptedMnemonic = cipherMnemonic
		cfg.Identity.EncryptedPrivKey = cipherPrivKey
		err = n.Repo.SetConfig(cfg)
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, &MessageOutput{"Password set."})
	},
	Type:     MessageOutput{},
	NoRemote: true,
}

var walletCheckPasswordCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "check password",
		ShortDescription: "check password",
	},
	Arguments: []cmds.Argument{},
	Options: []cmds.Option{
		cmds.StringOption(passwordOptionName, "p", "password"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}
		if err := checkWhetherPasswordSet(cfg); err != nil {
			return err
		}
		password, err := getPassword(req)
		if err != nil {
			return err
		}
		if err := validatePassword(cfg, password); err != nil {
			return err
		}
		return cmds.EmitOnce(res, &MessageOutput{"Password is correct."})
	},
	Type:     MessageOutput{},
	NoRemote: true,
}

var walletKeysCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "BTFS wallet keys",
		ShortDescription: "get keys of BTFS wallet",
	},
	Arguments: []cmds.Argument{},
	Options:   []cmds.Option{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}
		var keys *Keys
		if !cfg.UI.Wallet.Initialized {
			keys = &Keys{
				PrivateKey: cfg.Identity.PrivKey,
				Mnemonic:   cfg.Identity.Mnemonic,
			}
		} else {
			keys = &Keys{
				PrivateKey: cfg.Identity.EncryptedPrivKey,
				Mnemonic:   cfg.Identity.EncryptedMnemonic,
			}
		}
		return cmds.EmitOnce(res, keys)
	},
	Type:     Keys{},
	NoRemote: true,
}

type Keys struct {
	PrivateKey string
	Mnemonic   string
	SkInBase64 string
	SkInHex    string
}

var walletTransactionsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "BTFS wallet transactions",
		ShortDescription: "get transactions of BTFS wallet",
	},
	Arguments: []cmds.Argument{},
	Options:   []cmds.Option{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		txs, err := wallet.GetTransactions(n.Repo.Datastore(), n.Identity.Pretty())
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, txs)
	},
	Type: []*walletpb.TransactionV1{},
}

var walletTransferCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Send to another BTT wallet.",
		ShortDescription: "Send to another BTT wallet from current BTT wallet. Use '-p=<password>' to specific password.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("to", true, false, "address of another BTFS wallet to transfer to"),
		cmds.StringArg("amount", true, false, "amount of µBTT (=0.000001BTT) to transfer"),
		cmds.StringArg("memo", false, false, "attached memo"),
	},
	Options: []cmds.Option{
		cmds.StringOption(passwordOptionName, "p", "password"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}
		if err := checkWhetherPasswordSet(cfg); err != nil {
			return err
		}
		password, err := getPassword(req)
		if err != nil {
			return err
		}
		if err := validatePassword(cfg, password); err != nil {
			return err
		}
		amount, err := strconv.ParseInt(req.Arguments[1], 10, 64)
		if err != nil {
			return err
		}
		memo := ""
		if len(req.Arguments) == 3 {
			memo = req.Arguments[2]
		}
		ret, err := wallet.TransferBTTWithMemo(req.Context, n, cfg, nil, "", req.Arguments[0], amount, memo)
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("transaction %v sent", ret.TxId)
		return cmds.EmitOnce(res, &TransferResult{
			Result:  ret.Result,
			Message: msg,
		})
	},
	Type:     &TransferResult{},
	NoRemote: true,
}

func checkWhetherPasswordSet(cfg *config.Config) error {
	if cfg.Identity.EncryptedPrivKey == "" {
		return errors.New(`Password required, but none set.
Read 'btfs wallet password --help' and assign a password.`)
	}
	return nil
}

func getPassword(req *cmds.Request) (string, error) {
	password, _ := req.Options[passwordOptionName].(string)
	if password == "" {
		return "", errors.New("Password required, please use '-p <password>' to specify the password.")
	}
	return password, nil
}

func validatePassword(cfg *config.Config, password string) error {
	privK, err := wallet.DecryptWithAES(password, cfg.Identity.EncryptedPrivKey)
	if err != nil || cfg.Identity.PrivKey != privK {
		return errors.New("incorrect password")
	}
	return nil
}

type TransferResult struct {
	Result  bool
	Message string
}

const privateKeyOptionName = "privateKey"
const mnemonicOptionName = "mnemonic"

var walletImportCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "BTFS wallet import",
		ShortDescription: "import BTFS wallet",
	},
	Arguments: []cmds.Argument{},
	Options: []cmds.Option{
		cmds.StringOption(privateKeyOptionName, "p", "Private Key to import."),
		cmds.StringOption(mnemonicOptionName, "m", "Mnemonic to import."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		privKey, _ := req.Options[privateKeyOptionName].(string)
		mnemonic, _ := req.Options[mnemonicOptionName].(string)
		if privKey == "" && mnemonic == "" {
			return errors.New("required private key or mnemonic")
		}
		if err = doSetKeys(n, privKey, mnemonic); err != nil {
			return err
		}
		go path.DoRestart(false)
		return nil
	},
	NoRemote: true,
}

func doSetKeys(n *core.IpfsNode, privKey string, mnemonic string) error {
	return wallet.SetKeys(n, privKey, mnemonic)
}

var walletDiscoveryCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Speed wallet discovery",
		ShortDescription: "Speed wallet discovery",
	},
	Arguments: []cmds.Argument{},
	Options: []cmds.Option{
		cmds.StringOption(passwordOptionName, "p", "password"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}
		if cfg.UI.Wallet.Initialized {
			return errors.New("Already init, cannot discovery.")
		}
		key, err := wallet.DiscoverySpeedKey(req.Options[passwordOptionName].(string))
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, DiscoveryResult{Key: key})
	},
	NoRemote: true,
}

type DiscoveryResult struct {
	Key string
}
