# go-btfs

## What is BTFS 2.0?

BitTorrent File System (BTFS) is a next-generation file sharing protocol in the BitTorrent ecosystem. Current mainstream public blockchains mostly focus on computational tasks but lack cost-effective, scalable, and high-performing file storage and sharing solutions.

These are exactly what BTFS aims to clear up. Besides, underpinned by BTTC, BTFS enables cross-chain connectivity and multi-channel payments, making itself a more convenient choice. The intgration of BTFS, BitTorrent, and the BTTC network will boost DApp developers' efficiency in serving a wider market.

* The [documentation](https://docs.btfs.io/v2.0) walks developers through BTFS 2.0 setup, usage, and API references.
* Please join the BTFS community at [discord](https://discord.gg/PQWfzWS).

## BTFS 2.0 Architecture Diagram

![Architecture Diagram](https://files.readme.io/a21e9fb--min.png)

## Table of Contents
- [go-btfs](#go-btfs)
  - [What is BTFS 2.0?](#what-is-btfs-20)
  - [BTFS 2.0 Architecture Diagram](#btfs-20-architecture-diagram)
  - [Table of Contents](#table-of-contents)
  - [Faucet](#faucet)
  - [Install BTFS](#install-btfs)
    - [System Requirements](#system-requirements)
    - [Install Pre-Built Packages](#install-pre-built-packages)
      - [Initialize a BTFS Daemon](#initialize-a-btfs-daemon)
      - [Start the Daemon](#start-the-daemon)
    - [Build from Source](#build-from-source)
      - [Requires](#requires)
      - [Install Go](#install-go)
    - [Docker](#docker)
    - [Notices](#notices)
  - [Getting Started](#getting-started)
    - [Some things to try](#some-things-to-try)
    - [Usage](#usage)
  - [Development](#development)
    - [Development Dependencies](#development-dependencies)
    - [BTFS Gateway](#btfs-gateway)
  - [License](#license)

## Faucet

In order to ensure the normal use of btfs 2.0 testnet, you need to apply for BTT at BTTC testnet, which is obtained [**here**](https://testfaucet.bt.io/#/).

## Install BTFS

The download and install instructions for BTFS are over at: https://docs.btfs.io/v2.0/docs/install-run-btfs20-node.

### System Requirements

BTFS can run on most Linux, macOS, and Windows systems. We recommend
running it on a machine with at least 2 GB of RAM (it’ll do fine with
only one CPU core), but it should run fine with as little as 1 GB of
RAM. On systems with less memory, it may not be completely stable.
Only support compiling from source for mac and unix-based system.

### Install Pre-Built Packages

We host pre-built binaries at https://github.com/bittorrent/go-btfs/releases/latest

#### Initialize a BTFS Daemon
**On testnet**
```
$ btfs init -p storage-host-testnet
Generating TRON key with BIP39 seed phrase...
Master public key:  xpub661MyMwAqRbcGgHpeMqFkS5hnwoGeAcHG5KkDQwke7wFxtKqfsXTCTjWsoU2dYVXVGvV7EuGcviEzEJ143TezxxXvs2zZ9FYTtCei8iRQ66
initializing BTFS node at /Users/btfs/.btfs
generating btfs node keypair with TRON key...done
peer identity: 16Uiu2HAmKFQPM72SssFRrqcH1qwUsPwcp7vXSg3SEzfdYua1J5qc
to get started, enter:

    btfs cat /btfs/QmZjrLVdUpqVU6Pnc8pBnyQxVdpn9J8tfcsycP84W6N93C/readme
```

#### Start the Daemon

Start the BTFS Daemon
```
$ btfs daemon --chain-id <chainid>
```

Specify the chain for btfs to run by `--chain-id`, the chainid of the test network is `1029`, and the start command becomes: `btfs daemon --chain-id 1029`
```
$ btfs daemon --chain-id 1029
Initializing daemon...
go-btfs version: 2.0
Repo version: 10
System version: amd64/darwin
Golang version: go1.16.5
Repo location: /Users/btfs/.btfs
Peer identity: 16Uiu2HAmKFQPM72SssFRrqcH1qwUsPwcp7vXSg3SEzfdYua1J5qc
the address of Bttc format is:  0x7Cf4B71017F0312037D53fe966CE625BF98FFff6
the address of Tron format is:  TMMuwwxsuQGrDrN3aanc5y5r4FbibgLYDa
cannot continue until there is sufficient (30000 Suggested) BTT (for Gas) available on 0x7cf4b71017f0312037d53fe966ce625bf98ffff6
```

**Run the Daemon**

When starting the BTFS daemon for the first time, the system will create a node account and at the same time print a string of messages: cannot continue until there is sufficient (30000 Suggested) BTT (for Gas) available on After seeing such a message, it is necessary to recharge the node account with BTT through an external account, and the system suggests a minimum of 30000 BTT, which is used as gas to deploy a node vault contract by the node account.
After the recharge, the BTFS node will create the vault contract.

Get BTT on BTTC testnet reference [Faucet](#Faucet)
```
cannot continue until there is sufficient (30000 Suggested) BTT (for Gas) available on 0x7cf4b71017f0312037d53fe966ce625bf98ffff6 
self vault: 0x1f8b3e7e691d733f5eb17e5570c49de3e5aecef9 
Swarm listening on /ip4/127.0.0.1/tcp/4001
Swarm listening on /ip4/192.168.21.149/tcp/4001
Swarm listening on /ip6/::1/tcp/4001
Swarm listening on /p2p-circuit
Swarm announcing /ip4/127.0.0.1/tcp/4001
Swarm announcing /ip4/192.168.21.149/tcp/4001
Swarm announcing /ip6/::1/tcp/4001
API server listening on /ip4/127.0.0.1/tcp/5001
Dashboard: http://127.0.0.1:5001/dashboard
Gateway (readonly) server listening on /ip4/127.0.0.1/tcp/8080
Remote API server listening on /ip4/127.0.0.1/tcp/5101
Daemon is ready
```
At this point, the BTFS node is up and running


### Build from Source

#### Requires
* GO
* GNU make
* Git
* GCC (or some other go compatible C Compiler) (optional)

#### Install Go
If you need to update: [Download latest version of Go](https://golang.org/dl/).

You'll need to add Go's bin directories to your `$PATH` environment variable e.g., by adding these lines to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile`:

```
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$GOPATH/bin
```

(If you run into trouble, see the [Go install instructions](https://golang.org/doc/install)).

Clone the go-btfs repository
```
$ git clone https://github.com/bittorrent/go-btfs
```

Navigate to the go-btfs directory and run `make install`.
```
$ cd go-btfs
$ make install
```

A successful make install outputs something like:
```
$ make install
go: downloading github.com/tron-us/go-btfs-common v0.2.28
go: extracting github.com/tron-us/go-btfs-common v0.2.28
go: finding github.com/tron-us/go-btfs-common v0.2.28
go version go1.14.1 darwin/amd64
bin/check_go_version 1.14
go install  "-asmflags=all='-trimpath='" "-gcflags=all='-trimpath='" -ldflags="-X "github.com/bittorrent/go-btfs".CurrentCommit=e4848946d" ./cmd/btfs
```
Afterwards, run `btfs init` and `btfs daemon` to initialize and start the daemon. To re-initialize a new pair of keys, you can shut down the daemon first via `btfs shutdown`. Then run `rm -r .btfs` and `btfs init` again.

### Docker

Developers also have the option to build a BTFS daemon within a Docker container. After cloning the go-btfs repository, navigate into the go-btfs directory. This is where the Dockerfile is located. Build the docker image:
```
$ cd go-btfs
$ docker image build -t btfs_docker .   // Builds the docker image and tags "btfs_docker" as the name 
```

A successful build should have an output like:
```
Sending build context to Docker daemon  2.789MB
Step 1/37 : FROM golang:1.15
 ---> 4fe257ac564c
Step 2/37 : MAINTAINER bittorrent <support@tron.network>
 ---> Using cache
 ---> 02409001f528

...

Step 37/37 : CMD ["daemon", "--migrate=true"]
 ---> Running in 3660f91dce94
Removing intermediate container 3660f91dce94
 ---> b4e1523cf264
Successfully built b4e1523cf264
Successfully tagged btfs_docker:latest
```

Start the container based on the new image. Starting the container also initializes and starts the BTFS daemon.
```
$ docker container run --publish 5001:5001 --detach --name btfs1 btfs_docker
```

The CLI flags are as such:

* `--publish` asks Docker to forward traffic incoming on the host’s port 8080, to the container’s port 5001.
* `--detach` asks Docker to run this container in the background.
* `--name` specifies a name with which you can refer to your container in subsequent commands, in this case btfs1.

Configure cross-origin(CORS)
You need to configure cross-origin (CORS) to access the container from the host.
```
(host) docker exec -it btfs1 /bin/sh // Enter the container's shell
```

Then configure cross-origin(CORS) with btfs

```
(container) btfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["http://$IP:$PORT"]'
(container) btfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "GET", "POST"]'
```

E.g:
```
(container) btfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["http://localhost:5001"]'
(container) btfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "GET", "POST"]'
```

Exit the container and restart the container
```
(container) exit
(host) docker restart btfs1
```

You can access the container from the host with http://localhost:5001/webui.

Execute commands within the docker container:
```
docker exec CONTAINER btfs add FILE
```

### Notices
After upgrade to go-btfs v2.3.1, if you find that your node connectivity has deteriorated, you can execute the following command:

```
btfs config Swarm.ResourceMgr.Limits.System --json '{"ConnsInbound":0}'
```
Then restart your btfs node.
This command will help you unblock the inbound connections to improve connectivity and thus increase the chances of getting a contract.
## Getting Started

### Some things to try

Basic proof of 'btfs working' locally:

    echo "hello world" > hello
    btfs add hello
    # This should output a hash string that looks something like:
    # QmaN4MmXMduZe7Y7XoMKFPuDFunvEZU6DWtBPg3L8kkAuS
    btfs cat <that hash>

### Usage

```
  btfs  - Global p2p merkle-dag filesystem.

  btfs [--config=<config> | -c] [--debug | -D] [--help] [-h] [--api=<api>] [--offline] [--cid-base=<base>] [--upgrade-cidv0-in-output] [--encoding=<encoding> | --enc] [--timeout=<timeout>] <command> ...

SUBCOMMANDS
  BASIC COMMANDS
    init          Initialize btfs local configuration
    add <path>    Add a file to BTFS
    cat <ref>     Show BTFS object data
    get <ref>     Download BTFS objects
    ls <ref>      List links from an object
    refs <ref>    List hashes of links from an object

  BTFS COMMANDS
    storage       Manage client and host storage features
    rm            Clean up locally stored files and objects

  DATA STRUCTURE COMMANDS
    block         Interact with raw blocks in the datastore
    object        Interact with raw dag nodes
    files         Interact with objects as if they were a unix filesystem
    dag           Interact with IPLD documents (experimental)
    metadata      Interact with metadata for BTFS files

  ADVANCED COMMANDS
    daemon        Start a long-running daemon process
    mount         Mount an BTFS read-only mount point
    resolve       Resolve any type of name
    name          Publish and resolve BTNS names
    key           Create and list BTNS name keypairs
    dns           Resolve DNS links
    pin           Pin objects to local storage
    repo          Manipulate the BTFS repository
    stats         Various operational stats
    p2p           Libp2p stream mounting
    filestore     Manage the filestore (experimental)

  NETWORK COMMANDS
    id            Show info about BTFS peers
    bootstrap     Add or remove bootstrap peers
    swarm         Manage connections to the p2p network
    dht           Query the DHT for values or peers
    ping          Measure the latency of a connection
    diag          Print diagnostics

  TOOL COMMANDS
    config        Manage configuration
    version       Show btfs version information
    commands      List all available commands
    cid           Convert and discover properties of CIDs
    log           Manage and show logs of running daemon

  Use 'btfs <command> --help' to learn more about each command.

  btfs uses a repository in the local file system. By default, the repo is
  located at ~/.btfs. To change the repo location, set the $BTFS_PATH
  environment variable:

    export BTFS_PATH=/path/to/btfsrepo
```

## Development

Some places to get you started on the codebase:

- Main file: [./cmd/btfs/main.go](https://github.com/bittorrent/go-btfs/blob/master/cmd/btfs/main.go)
- CLI Commands: [./core/commands/](https://github.com/bittorrent/go-btfs/tree/master/core/commands)
- libp2p
  - libp2p: [libp2p](https://github.com/libp2p/go-libp2p)
  - DHT: [DHT](https://github.com/libp2p/go-libp2p-kad-dht)
  - PubSub: [PubSub](https://github.com/libp2p/go-libp2p-pubsub)

### Development Dependencies

If you make changes to the protocol buffers, you will need to install the [protoc compiler](https://github.com/google/protobuf).

### BTFS Gateway

BTFS Gateway is a free service that allows you to retrieve files from the BTFS network in your browser directly.

[How to use BTFS Gateway](https://docs.btfs.io/v2.0/docs/btfs-gateway-user-guide-1)

## License

[MIT](./LICENSE)

