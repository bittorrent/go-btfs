Tutorial for Upgrading to V2.1.2

This tutorial shows how to upgrade your node to v2.1.2. The main change in this upgradation is that we imporoved the security and maintainability both of the `VaultFactory` and `Vault` contract, and you can find the detail [here](https://github.com/bittorrent/btfs-vault/pull/8).

New version `VaultFactory` address is [`0x763d7858287B9a33F4bE5bb3df0241dACc59BCc7`](https://bttcscan.com/address/0x763d7858287B9a33F4bE5bb3df0241dACc59BCc7), and new version `Vault` logic contract address is [`0x11a91B7270ea000768F7A2C543547e832b5cb031`](https://bttcscan.com/address/0x11a91b7270ea000768f7a2c543547e832b5cb031), welcome for review.

1. Cashing Your Cheques before Upgrading

You had better cashing all your cheques you received before upgrading, and withdraw all your WBTT from vault to your BTTC address as well.

We provide a bash script to facilitate these burdensome operations, and you can run it as following. By the way, you can cash cheques manually via commands or dashboard. 

Note that the script will skip cashing cheques whose cashable amount is less than estimated gas fee.



```shell
# When API endpoint is 127.0.0.1:5001, which is the default one
$ curl -s https://raw.githubusercontent.com/bittorrent/go-btfs/master/scripts/batch_cash.sh | bash -s 127.0.0.1:5001

# or, if your API endpoint differs, please specify it
$ curl -s https://raw.githubusercontent.com/bittorrent/go-btfs/master/scripts/batch_cash.sh | bash -s <your-api-host>
```

Before running this script, please confirm that 'curl' and 'bc' tools are installed on your system. 
If not, you can install them through the system's corresponding package management tools, for example:

```shell
# on centos
$ sudo yum install -y curl
$ sudo yum install -y bc

# on ubuntu
$ sudo apt install curl
$ sudo apt install bc

# on macos
$ brew install curl
$ brew install bc
```

If your system is windows, you may need to install 'git-bash' and the corresponding 'bc' tool to run this script. 
If you have difficulty installing these tools, you can cash cheques and withdraw balance manually via commands or dashboard, just a little more tedious.

The scripts will prompts "Success, all tasks completed!" if everything gose well.
You can run the script again if it failed due to occasional reasons, weak network for example.

2. Upgrade your Vault to New Version

Upgradation will re-deploy a vault for you, so make sure there are at least 30000 BTT in your BTTC address. Note that your original private key and BTTC address won't be changed in this upgradation. Before upgrading, the program will back your `statestore` and `config` up in your `$BTFS_PATH`.

Upgrade steps:

1. Stop your BTFS daemon that are currently running;
2. Download BTFS v2.1.2 at https://github.com/bittorrent/go-btfs/releases;
3. Unzip downloaded file and rename it to `btfs`, and replace old `btfs` with it;
4. Run btfs daemon command with corrent `--chain-id` parameter, e.g.:

```shell
# If your node is a mainnet node, run:
$ btfs daemon --chain-id 199

# If your node is a testnet node, run:
$ btfs daemon --chain-id 1029

# and it will show logs like...

Initializing daemon...
go-btfs version: 2.1.2
Repo version: 10
System version: amd64/darwin
Golang version: go1.16.15
Repo location: ~/.btfs
Peer identity: 16Uiu2HAmQa9wH6Zx98rDEJvN6poFDDpBVjfHK6nKqPvEo2BUCydg
the address of Bttc format is:  0x69bBf2d2779F05dAdba6947aCD94c3e151bb93dF
the address of Tron format is:  TKcH76P3fcUWqGXBDEk7XCMBjLGsun9YvP

# upgrading related logs...

prepare upgrading your vault contract
backup statestore folder successfully to ~/.btfs/statestore.backup26
backup config file successfully to ~/.btfs/config.backup26
your old vault address is 0xaAFA1cEEF792aF65d68Edb545a542137A69A87Ba
will re-deploy a vault contract for you

cannot continue until there is sufficient (100 Suggested) BTT (for Gas) available on 0x69bbf2d2779f05dadba6947acd94c3e151bb93df 
cannot continue until there is sufficient (100 Suggested) BTT (for Gas) available on 0x69bbf2d2779f05dadba6947acd94c3e151bb93df 
self vault: 0x2d34448026415a7cc65a69f795528915ef685a5f
...
```
