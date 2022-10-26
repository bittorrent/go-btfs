module github.com/bittorrent/go-btfs

go 1.18

require (
	bazil.org/fuse v0.0.0-20200117225306-7b5117fecadc
	github.com/TRON-US/go-btfs-api v0.3.0
	github.com/TRON-US/go-btfs-chunker v0.3.0
	github.com/TRON-US/go-btfs-config v0.11.11
	github.com/TRON-US/go-btfs-files v0.2.0
	github.com/TRON-US/go-btfs-pinner v0.1.1
	github.com/TRON-US/go-btns v0.1.1
	github.com/TRON-US/go-eccrypto v0.0.1
	github.com/TRON-US/go-mfs v0.3.1
	github.com/TRON-US/go-unixfs v0.6.1
	github.com/TRON-US/interface-go-btfs-core v0.7.0
	github.com/Workiva/go-datastructures v1.0.52
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d
	github.com/bittorrent/go-btfs-cmds v0.2.14
	github.com/blang/semver v3.5.1+incompatible
	github.com/bren2010/proquint v0.0.0-20160323162903-38337c27106d
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/coreos/go-systemd/v22 v22.1.0
	github.com/dustin/go-humanize v1.0.0
	github.com/elgris/jsondiff v0.0.0-20160530203242-765b5c24c302
	github.com/ethereum/go-ethereum v1.10.3
	github.com/ethersphere/go-sw3-abi v0.4.0
	github.com/fsnotify/fsnotify v1.5.4
	github.com/gabriel-vasile/mimetype v1.1.2
	github.com/go-bindata/go-bindata/v3 v3.1.3
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/ip2location/ip2location-go/v9 v9.0.0
	github.com/ipfs/go-bitswap v0.2.20
	github.com/ipfs/go-block-format v0.0.2
	github.com/ipfs/go-blockservice v0.1.3
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-cidutil v0.0.2
	github.com/ipfs/go-datastore v0.4.5
	github.com/ipfs/go-detect-race v0.0.1
	github.com/ipfs/go-ds-badger v0.2.4
	github.com/ipfs/go-ds-flatfs v0.4.5
	github.com/ipfs/go-ds-leveldb v0.4.2
	github.com/ipfs/go-ds-measure v0.1.0
	github.com/ipfs/go-filestore v0.0.3
	github.com/ipfs/go-fs-lock v0.0.6
	github.com/ipfs/go-graphsync v0.2.0
	github.com/ipfs/go-ipfs-blockstore v0.1.4
	github.com/ipfs/go-ipfs-ds-help v0.1.1
	github.com/ipfs/go-ipfs-exchange-interface v0.0.1
	github.com/ipfs/go-ipfs-exchange-offline v0.0.1
	github.com/ipfs/go-ipfs-posinfo v0.0.1
	github.com/ipfs/go-ipfs-provider v0.4.3
	github.com/ipfs/go-ipfs-routing v0.1.0
	github.com/ipfs/go-ipfs-util v0.0.2
	github.com/ipfs/go-ipld-cbor v0.0.4
	github.com/ipfs/go-ipld-format v0.2.0
	github.com/ipfs/go-ipld-git v0.0.3
	github.com/ipfs/go-log v1.0.5
	github.com/ipfs/go-merkledag v0.3.2
	github.com/ipfs/go-metrics-interface v0.0.1
	github.com/ipfs/go-metrics-prometheus v0.0.2
	github.com/ipfs/go-path v0.0.8
	github.com/ipfs/go-verifcid v0.0.1
	github.com/ipld/go-car v0.1.1-0.20200429200904-c222d793c339
	github.com/jbenet/go-is-domain v1.0.5
	github.com/jbenet/go-random v0.0.0-20190219211222-123a90aedc0c
	github.com/jbenet/go-temp-err-catcher v0.1.0
	github.com/jbenet/goprocess v0.1.4
	github.com/klauspost/reedsolomon v1.9.9
	github.com/libp2p/go-libp2p v0.11.0
	github.com/libp2p/go-libp2p-circuit v0.4.0
	github.com/libp2p/go-libp2p-connmgr v0.2.4
	github.com/libp2p/go-libp2p-core v0.14.0
	github.com/libp2p/go-libp2p-crypto v0.1.0
	github.com/libp2p/go-libp2p-discovery v0.5.1
	github.com/libp2p/go-libp2p-http v0.1.5
	github.com/libp2p/go-libp2p-kad-dht v0.9.0
	github.com/libp2p/go-libp2p-kbucket v0.4.7
	github.com/libp2p/go-libp2p-loggables v0.1.0
	github.com/libp2p/go-libp2p-mplex v0.4.1
	github.com/libp2p/go-libp2p-noise v0.2.0
	github.com/libp2p/go-libp2p-peerstore v0.2.7
	github.com/libp2p/go-libp2p-pubsub v0.3.5
	github.com/libp2p/go-libp2p-pubsub-router v0.3.2
	github.com/libp2p/go-libp2p-record v0.1.3
	github.com/libp2p/go-libp2p-routing-helpers v0.2.3
	github.com/libp2p/go-libp2p-secio v0.2.2
	github.com/libp2p/go-libp2p-swarm v0.5.0
	github.com/libp2p/go-libp2p-testing v0.9.0
	github.com/libp2p/go-libp2p-tls v0.1.3
	github.com/libp2p/go-libp2p-yamux v0.5.4
	github.com/libp2p/go-socket-activation v0.0.2
	github.com/libp2p/go-tcp-transport v0.2.4
	github.com/libp2p/go-testutil v0.1.0
	github.com/libp2p/go-ws-transport v0.4.0
	github.com/looplab/fsm v0.1.0
	github.com/markbates/pkger v0.17.0
	github.com/mholt/archiver/v3 v3.3.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mr-tron/base58 v1.2.0
	github.com/multiformats/go-multiaddr v0.4.1
	github.com/multiformats/go-multiaddr-dns v0.2.0
	github.com/multiformats/go-multibase v0.0.3
	github.com/multiformats/go-multihash v0.0.15
	github.com/opentracing/opentracing-go v1.2.0
	github.com/orcaman/concurrent-map v0.0.0-20190826125027-8c72a8bb44f6
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.2
	github.com/shirou/gopsutil/v3 v3.20.12
	github.com/status-im/keycard-go v0.0.0-20200402102358-957c09536969
	github.com/stretchr/testify v1.8.0
	github.com/syndtr/goleveldb v1.0.1-0.20210305035536-64b5b1c73954
	github.com/thedevsaddam/gojsonq/v2 v2.5.2
	github.com/tron-us/go-btfs-common v0.8.10
	github.com/tron-us/go-common/v2 v2.3.0
	github.com/tron-us/protobuf v1.3.7
	github.com/tyler-smith/go-bip32 v0.0.0-20170922074101-2c9cfd177564
	github.com/tyler-smith/go-bip39 v1.1.0
	github.com/whyrusleeping/base32 v0.0.0-20170828182744-c30ac30633cc
	github.com/whyrusleeping/go-sysinfo v0.0.0-20190219211824-4a357d4b90b1
	github.com/whyrusleeping/multiaddr-filter v0.0.0-20160516205228-e903e4adabd7
	github.com/whyrusleeping/tar-utils v0.0.0-20180509141711-8c6c8ba81d5c
	go.uber.org/fx v1.13.1
	go.uber.org/zap v1.19.1
	go4.org v0.0.0-20200411211856-f5505b9728dd
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	golang.org/x/net v0.0.0-20220630215102-69896b714898
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f
	golang.org/x/sys v0.0.0-20220702020025-31831981b65f
	gopkg.in/cheggaaa/pb.v1 v1.0.28
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools v2.2.0+incompatible
)

require (
	github.com/AndreasBriese/bbloom v0.0.0-20190306092124-e2d15f34fcf9 // indirect
	github.com/FactomProject/basen v0.0.0-20150613233007-fe3947df716e // indirect
	github.com/FactomProject/btcutilecc v0.0.0-20130527213604-d3a63a5752ec // indirect
	github.com/Kubuxu/go-os-helper v0.0.1 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/Stebalien/go-bitfield v0.0.1 // indirect
	github.com/VictoriaMetrics/fastcache v1.5.7 // indirect
	github.com/alexbrainman/goissue34681 v0.0.0-20191006012335-3fc7a47baff5 // indirect
	github.com/andybalholm/brotli v0.0.0-20190621154722-5f990b63d2d6 // indirect
	github.com/benbjohnson/clock v1.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cmars/basen v0.0.0-20150613233007-fe3947df716e // indirect
	github.com/codemodus/kace v0.5.1 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/crackcomm/go-gitignore v0.0.0-20170627025303-887ab5e44cc3 // indirect
	github.com/cskr/pubsub v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/davidlazar/go-crypto v0.0.0-20170701192655-dcfb0a7ac018 // indirect
	github.com/deckarep/golang-set v0.0.0-20180603214616-504e848d77ea // indirect
	github.com/dgraph-io/badger v1.6.1 // indirect
	github.com/dgraph-io/ristretto v0.0.2 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/edsrzf/mmap-go v1.1.0 // indirect
	github.com/facebookgo/atomicfile v0.0.0-20151019160806-2de1f203e7d5 // indirect
	github.com/flynn/noise v1.0.0 // indirect
	github.com/fomichev/secp256k1 v0.0.0-20180413221153-00116ff8c62f // indirect
	github.com/gballet/go-libpcsclite v0.0.0-20190607065134-2772fd86a8ff // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-pg/migrations/v7 v7.1.11 // indirect
	github.com/go-pg/pg/v9 v9.2.0 // indirect
	github.com/go-pg/zerochecker v0.2.0 // indirect
	github.com/go-redis/redis/v7 v7.4.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gobuffalo/here v0.6.0 // indirect
	github.com/golang/gddo v0.0.0-20190419222130-af0f2af80721 // indirect
	github.com/golang/snappy v0.0.3-0.20201103224600-674baa8c7fc3 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/gopacket v1.1.19 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190910122728-9d188e94fb99 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/hannahhoward/go-pubsub v0.0.0-20200423002714-8d62886cc36e // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/holiman/uint256 v1.1.1 // indirect
	github.com/huin/goupnp v1.0.1-0.20210310174557-0ca763054c88 // indirect
	github.com/hypnoglow/go-pg-monitor v0.1.0 // indirect
	github.com/hypnoglow/go-pg-monitor/gopgv9 v0.1.0 // indirect
	github.com/ipfs/bbloom v0.0.4 // indirect
	github.com/ipfs/go-ipfs-delay v0.0.1 // indirect
	github.com/ipfs/go-ipfs-pq v0.0.2 // indirect
	github.com/ipfs/go-log/v2 v2.5.0 // indirect
	github.com/ipfs/go-peertaskqueue v0.2.0 // indirect
	github.com/ipld/go-ipld-prime v0.5.1-0.20201021195245-109253e8a018 // indirect
	github.com/ipld/go-ipld-prime-proto v0.1.0 // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/karalabe/usb v0.0.0-20190919080040-51dc0efba356 // indirect
	github.com/kisielk/errcheck v1.5.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/klauspost/cpuid v1.2.4 // indirect
	github.com/klauspost/cpuid/v2 v2.0.4 // indirect
	github.com/klauspost/pgzip v1.2.1 // indirect
	github.com/koron/go-ssdp v0.0.0-20191105050749-2e1c40ed0b5d // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/libp2p/go-addr-util v0.0.2 // indirect
	github.com/libp2p/go-buffer-pool v0.0.2 // indirect
	github.com/libp2p/go-cidranger v1.1.0 // indirect
	github.com/libp2p/go-conn-security-multistream v0.2.0 // indirect
	github.com/libp2p/go-eventbus v0.2.1 // indirect
	github.com/libp2p/go-flow-metrics v0.0.3 // indirect
	github.com/libp2p/go-libp2p-asn-util v0.0.0-20200825225859-85005c6cf052 // indirect
	github.com/libp2p/go-libp2p-autonat v0.3.2 // indirect
	github.com/libp2p/go-libp2p-blankhost v0.2.0 // indirect
	github.com/libp2p/go-libp2p-gostream v0.2.1 // indirect
	github.com/libp2p/go-libp2p-metrics v0.1.0 // indirect
	github.com/libp2p/go-libp2p-nat v0.0.6 // indirect
	github.com/libp2p/go-libp2p-netutil v0.1.0 // indirect
	github.com/libp2p/go-libp2p-peer v0.2.0 // indirect
	github.com/libp2p/go-libp2p-pnet v0.2.0 // indirect
	github.com/libp2p/go-libp2p-transport-upgrader v0.3.0 // indirect
	github.com/libp2p/go-mplex v0.1.2 // indirect
	github.com/libp2p/go-msgio v0.0.6 // indirect
	github.com/libp2p/go-nat v0.0.5 // indirect
	github.com/libp2p/go-netroute v0.2.0 // indirect
	github.com/libp2p/go-openssl v0.0.7 // indirect
	github.com/libp2p/go-reuseport v0.0.2 // indirect
	github.com/libp2p/go-reuseport-transport v0.0.4 // indirect
	github.com/libp2p/go-stream-muxer-multistream v0.3.0 // indirect
	github.com/libp2p/go-yamux v1.3.7 // indirect
	github.com/marten-seemann/tcp v0.0.0-20210406111302-dfbc87cc63fd // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/miekg/dns v1.1.31 // indirect
	github.com/mikioh/tcpinfo v0.0.0-20190314235526-30a79bb1804b // indirect
	github.com/mikioh/tcpopt v0.0.0-20190314235656-172688c1accc // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mmcloughlin/avo v0.0.0-20200523190732-4439b6b2c061 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/multiformats/go-base32 v0.0.3 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-multiaddr-fmt v0.1.0 // indirect
	github.com/multiformats/go-multiaddr-net v0.2.0 // indirect
	github.com/multiformats/go-multistream v0.2.0 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/nwaples/rardecode v1.0.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.17.0 // indirect
	github.com/peterh/liner v1.1.1-0.20190123174540-a2c9a5303de7 // indirect
	github.com/pierrec/lz4 v2.0.5+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/polydawn/refmt v0.0.0-20190809202753-05966cbd336a // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.35.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/prometheus/tsdb v0.7.1 // indirect
	github.com/rjeczalik/notify v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/segmentio/encoding v0.1.15 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/spacemonkeygo/spacelog v0.0.0-20180420211403-2296661a0572 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/texttheater/golang-levenshtein v0.0.0-20180516184445-d188e65d659e // indirect
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/tklauser/numcpus v0.2.2 // indirect
	github.com/ulikunitz/xz v0.5.6 // indirect
	github.com/vmihailenco/bufpool v0.1.11 // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.12 // indirect
	github.com/vmihailenco/tagparser v0.1.1 // indirect
	github.com/whyrusleeping/cbor-gen v0.0.0-20200710004633-5379fc63235d // indirect
	github.com/whyrusleeping/chunker v0.0.0-20181014151217-fe64bd25879f // indirect
	github.com/whyrusleeping/go-keyspace v0.0.0-20160322163242-5b898ac5add1 // indirect
	github.com/whyrusleeping/mdns v0.0.0-20190826153040-b9b60ed33aa9 // indirect
	github.com/whyrusleeping/timecache v0.0.0-20160911033111-cfcb2f1abfee // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	go.opencensus.io v0.22.4 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.10.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220609170525-579cf78fd858 // indirect
	golang.org/x/tools v0.1.10 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20211118181313-81c1377c94b1 // indirect
	google.golang.org/grpc v1.46.2 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	launchpad.net/gocheck v0.0.0-20140225173054-000000000087 // indirect
	mellium.im/sasl v0.2.1 // indirect
)

replace github.com/ipfs/go-ipld-format => github.com/TRON-US/go-ipld-format v0.2.0

replace github.com/ipfs/go-cid => github.com/TRON-US/go-cid v0.3.0

replace github.com/libp2p/go-libp2p-core => github.com/TRON-US/go-libp2p-core v0.7.1

replace github.com/libp2p/go-libp2p-kad-dht => github.com/TRON-US/go-libp2p-kad-dht v0.10.1

replace github.com/multiformats/go-multiaddr => github.com/TRON-US/go-multiaddr v0.4.0

replace github.com/ipfs/go-path => github.com/TRON-US/go-path v0.2.0

replace github.com/ipfs/go-graphsync => github.com/TRON-US/go-graphsync v0.2.1

replace github.com/ipld/go-car => github.com/TRON-US/go-car v0.3.0

replace github.com/ipld/go-ipld-prime-proto => github.com/TRON-US/go-ipld-prime-proto v0.1.0

replace github.com/libp2p/go-libp2p-yamux => github.com/libp2p/go-libp2p-yamux v0.2.8

replace github.com/libp2p/go-libp2p-swarm => github.com/libp2p/go-libp2p-swarm v0.2.8

replace github.com/libp2p/go-libp2p-mplex => github.com/libp2p/go-libp2p-mplex v0.2.4

replace github.com/libp2p/go-libp2p => github.com/libp2p/go-libp2p v0.11.0

replace github.com/libp2p/go-libp2p-circuit => github.com/libp2p/go-libp2p-circuit v0.3.1

exclude github.com/anacrolix/dht/v2 v2.15.2-0.20220123034220-0538803801cb
