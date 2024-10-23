- **btfs accesskey**                                        → Manage S3-Compatible-API access-keys
    1. **btfs accesskey delete  <key>**      →  generate a new access-key record
    2. **btfs accesskey disable  <key>**      → disable the specified access-key
    3. **btfs accesskey enable  <key>**       → enable the access-key
    4. **btfs accesskey generate**                → generate a new access-key record
    5. **btfs accesskey get   <key>**            → get an access-key detail info
    6. **btfs accesskey list**                          → list all access-keys
    7. **btfs accesskey reset**                       → reset secret of the specified access-key  

- **btfs add**
    1. **btfs add <fileName>**                                  → add a file to btfs
    2. **btfs add -r <directoryName>**                   → add a directory to btfs

- **btfs backup <file>**                                → back up BTFS’s data
- **btfs bitswap**                                           → interact with bitswap agent
    1. **btfs bitswap stat**                                     → show some diagnostic information on the             bitswap agent. must run in online mode
    2. **btfs bitswap wantlist**                              → show blocks that currently on the wantlist
    3. **btfs bitswap ledger  <peer>**                                                → show the current ledger for a peer
    4. **btfs bitswap reprovide**                             → trigger reprovider

- **btfs bittorrent**                                        → A tool command to integrate with bittorrent net(support bittorrent seed or a magnet URI scheme)
    1. **btfs bittorrent bencode  <path>**          → Print the bencoded info person-friendly of a bittorrent file from a bittorrent seed file
    2. **btfs bittorrent download  <uri>**              → Download a bittorrent file from the bittorrent seed or a magnet URL
    3. **btfs bittorrent metainfo**                             → print the metainfo of a bittorrent file from a seed file
    4. **btfs bittorrent scrape**                                → Fetch swarm metrics for info-hashes from tracker
    5. **btfs bittorrent serve   <path>**                  → serve as a bittorrent client with the specified files
    
- **btfs block**                                                   → interact with raw BTFS blocks
    1. **btfs block stat**                              → print information of a raw BTFS block
    2. **btfs block get**                               → get a raw BTFS block
    3. **btfs block put**                                → store input as an BTFS block
    4. **btfs block rm** **<hash>**                                 → remove BTFS blocks
    
- **btfs bootstrap**                                           → ****Show or edit the list of bootstrap peers
    1. **btfs bootstrap add  <peer>**              → Add peers to the bootstrap list
    2. **btfs bootstrap add default**              → Add default bootstrap nodes
    3. **btfs bootstrap list**                             → Show peers in the bootstrap list
    4. **btfs bootstrap rm  <peer>**                → Remove peers from the bootstrap list
    5. **btfs bootstrap rm all**                         → Remove all bootstrap peers
    
- **btfs bttc**                                                    → Interact with bttc related services
    1. **btfs bttc btt2wbtt  <amount>**                            → Swap BTT to WBTT at your bttc address
    2. **btfs bttc send-btt-to  <addr> <amount>**          → Transfer your BTT to other bttc address 
    3. **btfs bttc send-token-to  --token-type=<token-type>  <addr>  <amount>**      → Transfer your WBTT to other bttc address 
    4. **btfs bttc send-wbtt-to    <addr>  <amount>**       →  Transfer your WBTT to other bttc address           
    5. **btfs bttc wbtt2btt  <amount>**                            → swap WBTT to BTT at your bttc address
    
- **btfs cat   <hash>**                                      → show BTFS object data
- **btfs cheque**                                                 → Interact with vault services on BTFS
    1. **btfs cheque all_token_balance <addr>**            → Get all token balance by addr.
    2. **btfs cheque batch-cash <peer-ids>...**              → Batch cash the cheques by peerIDs.
    3. **btfs cheque bttbalance <addr>**                          → Get btt balance by addr.
    4. **btfs cheque cash <peer-id>**                                → Cash a cheque by peerID.
    5. **btfs cheque cashlist <from> <limit>**                  → Get cash status by peerID.
    6. **btfs cheque cashstatus <peer-id>**                     → Get cash status by peerID.
    7. **btfs cheque chaininfo**                                           → Show current chain info.
    8. **btfs cheque fix_cheque_cashout**                       → List cheque(s) received from peers.
    9. **btfs cheque price**                                                   → Get btfs token price.
    10. **btfs cheque price-all**                                             → Get btfs all price.
    11. **btfs cheque receive <peer-id>**                            → List cheque(s) received from peers.
    12. **btfs cheque receive-history-list <from> <limit>** → Display the received cheques from peer.
    13. **btfs cheque receive-history-peer <peer-id>**      → Display the received cheques from peer.
    14. **btfs cheque receive-history-stats**                         → Display the received cheques from peer.
    15. **btfs cheque receive-history-stats-all**                   → Display the received cheques from peer, of all tokens.
    16. **btfs cheque receive-total-count**                            → send cheque(s) count
    17. **btfs cheque receivelist <offset> <limit>**               → List cheque(s) received from peers.
    18. **btfs cheque receivelistall <offset> <limit>**           → List cheque(s) received from peers.
    19. **btfs cheque send <peer-id>**                                    → List cheque send to peers.
    20. **btfs cheque send-history-list <from> <limit>**     → Display the send cheques from peer.
    21. **btfs cheque send-history-peer <peer-id>**           → Display the send cheques from peer.
    22. **btfs cheque send-history-stats**                             → Display the received cheques from peer.
    23. **btfs cheque send-history-stats-all**                       → Display the received cheques from peer, of all tokens
    24. **btfs cheque send-total-count**                                → send cheque(s) count
    25. **btfs cheque sendlist**                                                 → List cheque(s) send to peers.
    26. **btfs cheque sendlistall**                                             → List cheque(s) send to peers.
    27. **btfs cheque stats**                                                      → List cheque(s) received from peers.
    28. **btfs cheque stats-all**                                                → List cheque(s) received from peers, of all tokens
    29. **btfs cheque token_balance <addr>**                      → Get one token balance by addr.
    
- **btfs cid**                                                    → convert and discover properties of CIDs
    1. **btfs cid base32 <cid>...**             → Convert CIDs to Base32 CID version 1.
    2. **btfs cid bases**                              → List available multibase encodings.
    3. **btfs cid codecs**                            → List available CID codecs.
    4. **btfs cid format <cid>...**              → Format and convert a CID in various useful ways.
    5. **btfs cid hashes**                            → List available multihashes.
    
- **btfs commands**                                        → list all commands
    
    
- **btfs config  <key>  <value>**                  → get and set btfs config values
    1. **btfs config edit**                                          → Open the config file for editing in $EDITOR.
    2. **btfs config optin**                                        → Opt-in enables analytic data collection (default).
    3. **btfs config optout**                                     → Opt-out disables collection of the analytics data(enabled by default).
    4. **btfs config replace <file>**                         → Replace the config with <file>.
    5. **btfs config reset**                                        → Reset config file contents.
    6. **btfs config show**                                       → Output config file contents.
    7. **btfs config storage-host-enable <enable>** → host is or not.
    8. **btfs config sync-chain-info**                          → sync chain info.
    9. **btfs config sync-simple-mode <value>**     → simple mode is true or not.

- **btfs daemon**                                             →  Run a network-connected BTFS node
    
    
- **btfs dag**                                                     → Interact with ipld dag objects
    1. **btfs dag export <root> **                     → Streams the selected DAG as a .car stream on stdout.
    2. **btfs dag get <ref>**                              → Get a dag node from btfs.
    3. **btfs dag import <path> **                 → Import the contents of .car files
    4. **btfs dag put <object data>**           → Add a dag node to btfs.
    5. **btfs dag resolve <ref>**                       → Resolve ipld block
    6. **btfs dag stat <root>**                           → Gets stats for a DAG

- **btfs decrypt**  **<cid>**                                            → decrypt the content of a CID with the private key of this peer
- **btfs dht**                                                     → Issue commands directly through the DHT.
    1. **btfs dht findpeer <peerID>...**          → Find the multiaddresses associated with a Peer ID.
    2. **btfs dht findprovs <key>...**              → Find peers that can provide a specific value, given a key.
    3. **btfs dht get <key>...**                         → Given a key, query the routing system for its best value.
    4. **btfs dht provide <key>...**                 → Announce to the network that you are providing given values.
    5. **btfs dht put <key> <value-file>**      → Write a key/value pair to the routing system.
    6. **btfs dht query <peerID>...**              → Find the closest Peer IDs to a given Peer ID by querying the DHT.

- **btfs diag**                                                    → Generate diagnostic reports
    1. **btfs diag cmds**                                  → List commands run on this BTFS node.
    2. **btfs diag cmds clear**                         → Clear inactive requests from the log
    3. **btfs diag cmds set-time  <time>**    → Set how long to keep inactive requests in the log
    4. **btfs diag sys**                                      → Print system diagnostic information.