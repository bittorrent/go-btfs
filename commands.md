# BTFS Command Reference

## **btfs accesskey**
Manage S3-Compatible API access keys.
1. **btfs accesskey delete \<key\>** → Generate a new access key record.
2. **btfs accesskey disable \<key\>** → Disable the specified access key.
3. **btfs accesskey enable \<key\>** → Enable the access key.
4. **btfs accesskey generate** → Generate a new access key record.
5. **btfs accesskey get \<key\>** → Get detailed information for an access key.
6. **btfs accesskey list** → List all access keys.
7. **btfs accesskey reset** → Reset the secret of the specified access key.

## **btfs add**
1. **btfs add \<fileName\>** → Add a file to BTFS.
2. **btfs add -r \<directoryName\>** → Add a directory to BTFS.

## **btfs backup \<file\>**
Back up BTFS data.

## **btfs bitswap**
Interact with the bitswap agent.
1. **btfs bitswap stat** → Show diagnostic information about the bitswap agent (must run in online mode).
2. **btfs bitswap wantlist** → Show blocks currently on the wantlist.
3. **btfs bitswap ledger \<peer\>** → Show the current ledger for a peer.
4. **btfs bitswap reprovide** → Trigger reprovider.

## **btfs bittorrent**
Integrate with the BitTorrent network (supports BitTorrent seed or a magnet URI scheme).
1. **btfs bittorrent bencode \<path\>** → Print the bencoded info from a BitTorrent seed file.
2. **btfs bittorrent download \<uri\>** → Download a BitTorrent file from a seed or a magnet URL.
3. **btfs bittorrent metainfo** → Print the metainfo of a BitTorrent file.
4. **btfs bittorrent scrape** → Fetch swarm metrics for info-hashes from the tracker.
5. **btfs bittorrent serve \<path\>** → Serve as a BitTorrent client with the specified files.

## **btfs block**
Interact with raw BTFS blocks.
1. **btfs block stat** → Print information about a raw BTFS block.
2. **btfs block get \<hash\>** → Get a raw BTFS block.
3. **btfs block put \<data\>** → Store input as a BTFS block.
4. **btfs block rm \<hash\>** → Remove a BTFS block.

## **btfs bootstrap**
Show or edit the list of bootstrap peers.
1. **btfs bootstrap add \<peer\>** → Add a peer to the bootstrap list.
2. **btfs bootstrap add default** → Add default bootstrap nodes.
3. **btfs bootstrap list** → Show peers in the bootstrap list.
4. **btfs bootstrap rm \<peer\>** → Remove a peer from the bootstrap list.
5. **btfs bootstrap rm all** → Remove all bootstrap peers.

## **btfs bttc**
Interact with BTT and WBTT services.
1. **btfs bttc btt2wbtt \<amount\>** → Swap BTT to WBTT at your BTT address.
2. **btfs bttc send-btt-to \<addr\> \<amount\>** → Transfer your BTT to another BTT address.
3. **btfs bttc send-token-to --token-type=\<token-type\> \<addr\> \<amount\>** → Transfer your WBTT to another BTT address.
4. **btfs bttc send-wbtt-to \<addr\> \<amount\>** → Transfer your WBTT to another BTT address.
5. **btfs bttc wbtt2btt \<amount\>** → Swap WBTT to BTT at your BTT address.

## **btfs cat \<hash\>**
Show BTFS object data.

## **btfs cheque**
Interact with vault services on BTFS.
1. **btfs cheque all_token_balance \<addr\>** → Get all token balances by address.
2. **btfs cheque batch-cash \<peer-ids\>...** → Batch cash the cheques by peer IDs.
3. **btfs cheque bttbalance \<addr\>** → Get BTT balance by address.
4. **btfs cheque cash \<peer-id\>** → Cash a cheque by peer ID.
5. **btfs cheque cashlist \<from\> \<limit\>** → Get cash status by peer ID.
6. **btfs cheque cashstatus \<peer-id\>** → Get cash status by peer ID.
7. **btfs cheque chaininfo** → Show current chain info.
8. **btfs cheque fix_cheque_cashout** → List cheques received from peers.
9. **btfs cheque price** → Get BTFS token price.
10. **btfs cheque price-all** → Get all BTFS token prices.
11. **btfs cheque receive \<peer-id\>** → List cheques received from peers.
12. **btfs cheque receive-history-list \<from\> \<limit\>** → Display the received cheques from a peer.
13. **btfs cheque receive-history-peer \<peer-id\>** → Display received cheques from a specific peer.
14. **btfs cheque receive-history-stats** → Display received cheques from a peer.
15. **btfs cheque receive-history-stats-all** → Display received cheques from all peers.
16. **btfs cheque receive-total-count** → Get total count of received cheques.
17. **btfs cheque receivelist \<offset\> \<limit\>** → List received cheques from peers.
18. **btfs cheque receivelistall \<offset\> \<limit\>** → List all received cheques from peers.
19. **btfs cheque send \<peer-id\>** → List cheques sent to peers.
20. **btfs cheque send-history-list \<from\> \<limit\>** → Display sent cheques from a peer.
21. **btfs cheque send-history-peer \<peer-id\>** → Display sent cheques from a specific peer.
22. **btfs cheque send-history-stats** → Display sent cheques from a peer.
23. **btfs cheque send-history-stats-all** → Display sent cheques from all peers.
24. **btfs cheque send-total-count** → Get total count of sent cheques.
25. **btfs cheque sendlist** → List cheques sent to peers.
26. **btfs cheque sendlistall** → List all cheques sent to peers.
27. **btfs cheque stats** → List cheques received from peers.
28. **btfs cheque stats-all** → List cheques received from all peers.
29. **btfs cheque token_balance \<addr\>** → Get token balance by address.

## **btfs cid**
Convert and discover properties of CIDs.
1. **btfs cid base32 \<cid\>...** → Convert CIDs to Base32 CID version 1.
2. **btfs cid bases** → List available multibase encodings.
3. **btfs cid codecs** → List available CID codecs.
4. **btfs cid format \<cid\>...** → Format and convert a CID in various ways.
5. **btfs cid hashes** → List available multihashes.

## **btfs commands**
List all commands.

## **btfs config \<key\> \<value\>**
Get and set BTFS config values.
1. **btfs config edit** → Open the config file for editing in $EDITOR.
2. **btfs config optin** → Opt-in to enable analytic data collection (default).
3. **btfs config optout** → Opt-out of data collection (enabled by default).
4. **btfs config replace \<file\>** → Replace the config with the specified file.
5. **btfs config reset** → Reset config file contents.
6. **btfs config show** → Output config file contents.
7. **btfs config storage-host-enable \<enable\>** → Enable or disable storage hosting.
8. **btfs config sync-chain-info** → Sync chain info.
9. **btfs config sync-simple-mode \<value\>** → Set simple mode to true or false.

## **btfs daemon**
Run a network-connected BTFS node.

## **btfs dag**
Interact with IPLD DAG objects.
1. **btfs dag export \<root\>** → Stream the selected DAG as a .car stream on stdout.
2. **btfs dag get \<ref\>** → Get a DAG node from BTFS.
3. **btfs dag import \<path\>...** → Import the contents of .car files.
4. **btfs dag put \<object data\>...** → Add a DAG node to BTFS.
5. **btfs dag resolve \<ref\>** → Resolve IPLD block.
6. **btfs dag stat \<root\>** → Get stats for a DAG.

## **btfs decrypt \<cid\>**
Decrypt the content of a CID with the private key of this peer.

## **btfs dht**
Issue commands directly through the DHT.
1. **btfs dht findpeer \<peerID\>...** → Find the multiaddresses associated with a Peer ID.
2. **btfs dht findprovs \<key\>...** → Find peers that can provide a specific value, given a key.
3. **btfs dht get \<key\>...** → Query the routing system for the best value for a key.
4. **btfs dht provide \<key\>...** → Announce to the network that you are providing given values.
5. **btfs dht put \<key\> \<value-file\>** → Write a key/value pair to the routing system.
6. **btfs dht query \<peerID\>...** → Find the closest Peer IDs to a given Peer ID by querying the DHT.

## **btfs diag**
Generate diagnostic reports.
1. **btfs diag cmds** → List commands run on this BTFS node.
2. **btfs diag cmds clear** → Clear inactive requests from the log.
3. **btfs diag cmds set-time \<time\>** → Set how long to keep inactive requests in the log.
4. **btfs diag sys** → Print system diagnostic information.

## **btfs dns**
Resolve DNS links.

## **btfs encrypt \<path\>...**
Encrypt a file with the public key of the peer.

## **btfs file**
Interact with BTFS objects representing Unix filesystems.
1. **btfs file ls \<btfs-path\>...** → List directory contents for Unix filesystem objects.

## **btfs files**
Interact with UnixFS files.
1. **btfs files chcid [\<path\>]** → Change the CID version or hash function of the root node of a given path.
2. **btfs files cp \<source\> \<dest\>** → Copy BTFS files and directories into MFS (or copy within MFS).
3. **btfs files flush [\<path\>]** → Flush a given path's data to disk.
4. **btfs files ls [\<path\>]** → List directories in the local mutable namespace.
5. **btfs files mkdir \<path\>** → Create directories.
6. **btfs files mv \<source\> \<dest\>** → Move files.
7. **btfs files read \<path\>** → Read a file in a given MFS.
8. **btfs files rm \<path\>...** → Remove a file.
9. **btfs files stat \<path\>** → Display file status.
10. **btfs files write \<path\> \<data\>** → Write to a mutable file in a given filesystem.

## **btfs filestore**
Interact with filestore objects.
1. **btfs filestore dups** → List blocks in both the filestore and standard block storage.
2. **btfs filestore ls [\<obj\>]...** → List objects in the filestore.
3. **btfs filestore verify [\<obj\>]...** → Verify objects in the filestore.

## **btfs get**
Download BTFS objects.

## **btfs guard \<btfs-path\>**
Interact with guard services from the BTFS client.
1. **btfs guard test** → Send tests to guard service endpoints from the BTFS client.
2. **btfs guard test send-challenges** → Send shard challenge questions from the BTFS client.

## **btfs id \<peer-id\>**
Show BTFS node ID info.

## **btfs init**
Initialize BTFS local configuration.

## **btfs key**
Create and list BTNS name keypairs.
1. **btfs key gen \<name\>** → Create a new keypair.
2. **btfs key list** → List all local keypairs.
3. **btfs key rename \<name\> \<newName\>** → Rename a keypair.
4. **btfs key rm \<name\>...** → Remove a keypair.

## **btfs log**
Interact with the daemon log output.
1. **btfs log level \<subsystem\> \<level\>** → Change the logging level.
2. **btfs log ls** → List logging subsystems.
3. **btfs log tail** → Read the event log.

## **btfs ls**
List directory contents for Unix filesystem objects.

## **btfs metadata**
Interact with metadata for BTFS files.
1. **btfs metadata add \<file-hash\> \<metadata\>** → Add token metadata to a BTFS file.
2. **btfs metadata rm \<file-hash\> \<metadata\>** → Remove token metadata from a BTFS file.

## **btfs mount**
Mount BTFS to the filesystem (read-only).

## **btfs multibase**
Encode and decode files or stdin with multibase format.
1. **btfs multibase decode \<encoded_file\>** → Decode a multibase string.
2. **btfs multibase encode \<file\>** → Encode data into a multibase string.
3. **btfs multibase list** → List available multibase encodings.
4. **btfs multibase transcode \<encoded_file\>** → Transcode a multibase string between bases.

## **btfs name**
Publish and resolve BTFS names.
1. **btfs name publish \<btfs-path\>** → Publish BTNS names.
2. **btfs name pubsub** → Manage BTNS pubsub.
3. **btfs name resolve [\<name\>]** → Resolve BTNS names.

## **btfs network**
Get BTFS network information.

## **btfs object**
Interact with BTFS objects.
1. **btfs object data \<key\>** → Output the raw bytes of a BTFS object.
2. **btfs object diff \<obj_a\> \<obj_b\>** → Display the diff between two BTFS objects.
3. **btfs object get \<key\>** → Get and serialize the DAG node named by \<key\>.
4. **btfs object links \<key\>** → Output the links pointed to by the specified object.
5. **btfs object new [\<template\>]** → Create a new object from a BTFS template.
6. **btfs object patch** → Create a new MerkleDAG object based on an existing one.
7. **btfs object patch append-data** → Append data to the data segment of a DAG node.
8. **btfs object patch add-link** → Add a link to a given object.
9. **btfs object patch rm-link** → Remove a link from a given object.
10. **btfs object patch set-data** → Set the data field of a BTFS object.
11. **btfs object put \<data\>** → Store input as a DAG object and print its key.
12. **btfs object stat \<key\>** → Get stats for the DAG node named by \<key\>.

## **btfs p2p**
LibP2P stream mounting.
1. **btfs p2p close** → Stop listening for new connections to forward.
2. **btfs p2p forward \<protocol\> \<listen-address\> \<target-address\>** → Forward connections to a LibP2P service.
3. **btfs p2p handshake \<chain-id\> \<peer-id\>** → P2P handshake.
4. **btfs p2p listen \<protocol\> \<target-address\>** → Create LibP2P service.
5. **btfs p2p ls** → List active P2P listeners.
6. **btfs p2p stream** → Manage P2P streams.
7. **btfs p2p stream close** → Close active P2P stream.
8. **btfs p2p stream ls** → List active P2P streams.

## **btfs pin**
Pin (and unpin) objects to local storage.
1. **btfs pin add \<btfs-path\>...** → Pin objects to local storage.
2. **btfs pin ls [\<btfs-path\>]...** → List objects pinned to local storage.
3. **btfs pin rm \<btfs-path\>...** → Remove pinned objects from local storage.
4. **btfs pin update \<from-path\> \<to-path\>** → Update a recursive pin.
5. **btfs pin verify** → Verify that recursive pins are complete.

## **btfs ping**
Send echo request packets to BTFS hosts.

## **btfs pubsub**
An experimental publish-subscribe system on BTFS.
1. **btfs pubsub ls** → List subscribed topics by name.
2. **btfs pubsub peers [\<topic\>]** → List peers we are currently pubsubbing with.
3. **btfs pubsub pub \<topic\> \<data\>...** → Publish a message to a given pubsub topic.
4. **btfs pubsub sub \<topic\>** → Subscribe to messages on a given topic.

## **btfs recovery**
Recover BTFS data from an archived backup file.

## **btfs refs**
List links (references) from an object.
1. **btfs refs local** → List all local references.

## **btfs repo**
Manipulate the BTFS repository.
1. **btfs repo fsck** → Remove repo lockfiles.
2. **btfs repo gc** → Perform a garbage collection sweep on the repo.
3. **btfs repo stat** → Get stats for the currently used repo.
4. **btfs repo verify** → Verify all blocks in the repo are not corrupted.
5. **btfs repo version** → Show the repo version.

## **btfs resolve**
Resolve the value of names to BTFS.

## **btfs restart**
Restart the daemon.

## **btfs rm**
Remove files or directories from a local BTFS node.

## **btfs settlement**
Interact with chequebook services on BTFS.
1. **btfs settlement list** → List all settlements.
2. **btfs settlement peer \<peer-id\>** → Get chequebook balance.

## **btfs shutdown**
Shut down the BTFS daemon.

## **btfs stats**
Query BTFS statistics.
1. **btfs stats bitswap** → Show diagnostic information on the bitswap agent.
2. **btfs stats bw** → Print BTFS bandwidth information.
3. **btfs stats dht [\<dht\>]...** → Return statistics about the node's DHT(s).
4. **btfs stats repo** → Get stats for the currently used repo.

## **btfs statuscontract**
Report status-contract commands.
1. **btfs statuscontract config** → Get reporting status-contract config.
2. **btfs statuscontract daily_last_report_time** → Report total status-contract info (total count, total gas spent, and contract address).
3. **btfs statuscontract daily_report_list \<from\> \<limit\>** → Report daily list with pagination.
4. **btfs statuscontract daily_report_online_server** → Daily report online server.
5. **btfs statuscontract daily_total** → Report total status-contract info.
6. **btfs statuscontract lastinfo** → Get the last reporting status-contract info.
7. **btfs statuscontract report_online_server** → Report online server.
8. **btfs statuscontract reportlist \<from\> \<limit\>** → Report status-contract list with pagination.

## **btfs storage**
Interact with storage services on BTFS.
1. **btfs storage announce** → Update and announce storage host information.
2. **btfs storage challenge** → Interact with storage challenge requests and responses.
3. **btfs storage challenge request** → Challenge storage hosts with Proof-of-Storage requests.
4. **btfs storage challenge response** → Storage host responds to Proof-of-Storage requests.
5. **btfs storage contracts** → Get node storage contracts info.
6. **btfs storage contracts sync** → Synchronize contract stats based on role.
7. **btfs storage contract list** → Get contracts list based on role.
8. **btfs storage contract stat** → Get contracts stats based on role.
9. **btfs storage dcrepair** → Interact with host repair requests and responses for decentralized shards repair.
10. **btfs storage dcrepair request** → Negotiate with hosts for repair jobs.
11. **btfs storage dcrepair response** → Host responds to repair jobs.
12. **btfs storage hosts** → Interact with information on hosts.
13. **btfs storage hosts info** → Display saved host information.
14. **btfs storage hosts sync** → Synchronize host information from BTFS hub.
15. **btfs storage info [\<peer-id\>]** → Show storage host information.
16. **btfs storage path \<path-name\> \<storage-size\>** → Modify the host storage folder path for the BTFS client.
17. **btfs storage path capacity** → Get free space of the specified path.
18. **btfs storage path list** → List directories.
19. **btfs storage path migrate** → Migrate path (e.g., `btfs storage path migrate /Users/tron/.btfs.new`).
20. **btfs storage path mkdir** → Create a directory.
21. **btfs storage path status** → Get status of resetting the path.
22. **btfs storage path volumes** → List disk volumes.
23. **btfs storage stats** → Get node storage stats.
24. **btfs storage stats info** → Get node stats.
25. **btfs storage stats list** → List node stats.
26. **btfs storage stats sync** → Synchronize node stats.
27. **btfs storage upload \<file-hash\> [\<upload-peer-id\>] [\<upload-nonce-ts\>] [\<upload-signature\>]** → Store files on BTFS network nodes through BTT payment.
28. **btfs storage upload cheque \<encoded-cheque\> \<amount\> [\<contract-id\>] [\<token\>]** → Receive upload cheque and process it.
29. **btfs storage upload getcontractbatch \<session-id\> \<peer-id\> \<nonce-timestamp\> \<upload-session-signature\> \<contracts-type\>** → Get all contracts from the upload session.
30. **btfs storage upload getunsigned \<session-id\> \<peer-id\> \<nonce-timestamp\> \<upload-session-signature\> \<session-status\>** → Get the input data for upload signing.
31. **btfs storage upload init \<session-id\> \<file-hash\> \<shard-hash\> \<price\> \<escrow-contract\> \<guard-contract-meta\> \<storage-length\> \<shard-size\> \<shard-index\> [\<upload-peer-id\>]** → Initialize storage handshake with the client.
32. **btfs storage upload recvcontract \<session-id\> \<shard-hash\> \<shard-index\> \<escrow-contract\> \<guard-contract\>** → For renter client to receive half-signed contracts.
33. **btfs storage upload repair \<file-hash\> \<repair-shards\> \<renter-pid\> \<blacklist\>** → Repair specific shards of a file.
34. **btfs storage upload sign \<session-id\> \<peer-id\> \<nonce-timestamp\> \<upload-session-signature\> \<session-status\> \<signed\>** → Return the signed data to the upload session.
35. **btfs storage upload signcontractbatch \<session-id\> \<peer-id\> \<nonce-timestamp\> \<upload-session-signature\> \<contracts-type\> \<signed-data-items\>** → Get the unsigned contracts from the upload session.
36. **btfs storage upload status \<session-id\>** → Check storage upload and payment status (from client's perspective).
37. **btfs storage upload supporttokens** → Support cheques and return tokens.

## **btfs swarm**
Interact with the swarm.
1. **btfs swarm addrs** → List known addresses (useful for debugging).
2. **btfs swarm addrs listen** → List interface listening addresses.
3. **btfs swarm addrs local** → List local addresses.
4. **btfs swarm connect \<address\>...** → Open a connection to a given address.
5. **btfs swarm disconnect \<address\>...** → Close connection to a given address.
6. **btfs swarm filters** → Manipulate address filters.
7. **btfs swarm filters add \<address\>...** → Add an address filter.
8. **btfs swarm filters rm \<address\>...** → Remove an address filter.
9. **btfs swarm peers** → List peers with open connections.

## **btfs tar**
Utility functions for tar files in BTFS.
1. **btfs tar add \<file\>** → Import a tar file into BTFS.
2. **btfs tar cat \<path\>** → Export a tar file from BTFS.

## **btfs test**
1. **btfs test cheque** → Show peers in the bootstrap list.
2. **btfs test hosts** → Show peers in the bootstrap list.
3. **btfs test p2phandshake** → P2P handshake.

## **btfs urlstore**
Interact with urlstore.
1. **btfs urlstore add \<url\>** → Add a URL via urlstore.

## **btfs vault**
Interact with vault services on BTFS.
1. **btfs vault address** → Get vault address.
2. **btfs vault balance** → Get vault balance.
3. **btfs vault balance_all** → Get vault balance.
4. **btfs vault deposit \<amount\>** → Deposit from beneficiary to vault contract account.
5. **btfs vault upgrade** → Upgrade vault contract to the latest version.
6. **btfs vault wbttbalance \<addr\>** → Get WBTT balance by address.
7. **btfs vault withdraw \<amount\>** → Withdraw from vault contract account to beneficiary.

## **btfs version**
Show BTFS version information.
1. **btfs version deps** → Show information about build dependencies.
