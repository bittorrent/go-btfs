package spin

import (
	"context"
	"fmt"
	"strings"
	"time"

	onlinePb "github.com/bittorrent/go-btfs-common/protos/online"
	cgrpc "github.com/bittorrent/go-btfs-common/utils/grpc"
	config "github.com/bittorrent/go-btfs-config"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core"

	"github.com/cenkalti/backoff/v4"
	"github.com/gogo/protobuf/proto"
	ic "github.com/libp2p/go-libp2p/core/crypto"
)

func isReportOnlineEnabled(cfg *config.Config) bool {
	return cfg.Experimental.StorageHostEnabled || cfg.Experimental.ReportOnline
}

func (dc *dcWrap) doSendDataOnline(ctx context.Context, config *config.Config, sm *onlinePb.ReqSignMetrics) error {
	onlineService := config.Services.OnlineServerDomain
	if len(onlineService) <= 0 {
		onlineService = chain.GetOnlineServer(config.ChainInfo.ChainId)
	}
	cb := cgrpc.OnlineClient(onlineService)
	return cb.WithContext(ctx, func(ctx context.Context, client onlinePb.OnlineServiceClient) error {
		resp, err := client.UpdateSignMetrics(ctx, sm)
		if err != nil {
			chain.CodeStatus = chain.ConstCodeError
			chain.ErrStatus = err
			return err
		} else {
			chain.CodeStatus = chain.ConstCodeSuccess
			chain.ErrStatus = nil
		}

		//fmt.Printf("--- online, resp, SignedInfo = %+v, signature = %+v \n", resp.SignedInfo, resp.Signature)
		if (resp.SignedInfo != nil) && len(resp.SignedInfo.Peer) > 0 {
			onlineInfo := chain.LastOnlineInfo{
				LastSignedInfo: onlinePb.SignedInfo{
					Peer:        resp.SignedInfo.Peer,
					CreatedTime: resp.SignedInfo.CreatedTime,
					Version:     resp.SignedInfo.Version,
					Nonce:       resp.SignedInfo.Nonce,
					BttcAddress: resp.SignedInfo.BttcAddress,
					SignedTime:  resp.SignedInfo.SignedTime,
				},
				LastSignature: resp.Signature,
				LastTime:      time.Now(),
			}
			chain.StoreOnline(&onlineInfo)
		}

		return err
	})
}

func (dc *dcWrap) SendDataOnline(node *core.IpfsNode, config *config.Config) {
	sm, errs, err := dc.doPrepDataOnline(node)
	if errs == nil {
		errs = make([]error, 0)
	}
	var sb strings.Builder
	if err != nil {
		errs = append(errs, err)
	}
	for _, err := range errs {
		sb.WriteString(err.Error())
		sb.WriteRune('\n')
	}
	log.Debug(sb.String())
	// If complete prep failure we return
	if err != nil {
		return
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = maxRetryTotal
	backoff.Retry(func() error {
		err := dc.doSendDataOnline(node.Context(), config, sm)
		if err != nil {
			//fmt.Printf("--- online, doSendDataOnline error = %+v \n", err)
			log.Infof("failed： doSendDataOnline to online server: %+v ", err)
		} else {
			log.Debug("sent OK, doSendDataOnline to online server")
		}
		return err
	}, bo)
}

// doPrepData gathers the latest analytics and returns (signed object, list of reporting errors, failure)
func (dc *dcWrap) doPrepDataOnline(btfsNode *core.IpfsNode) (*onlinePb.ReqSignMetrics, []error, error) {
	errs := dc.update(btfsNode)
	payload, err := dc.getPayloadOnline(btfsNode)
	if err != nil {
		return nil, errs, fmt.Errorf("failed to marshal dataCollection object to a byte array: %s", err.Error())
	}
	if dc.node.PrivateKey == nil {
		return nil, errs, fmt.Errorf("node's private key is null")
	}

	signature, err := dc.node.PrivateKey.Sign(payload)
	if err != nil {
		return nil, errs, fmt.Errorf("failed to sign raw data with node private key: %s", err.Error())
	}

	publicKey, err := ic.MarshalPublicKey(dc.node.PrivateKey.GetPublic())
	if err != nil {
		return nil, errs, fmt.Errorf("failed to marshal node public key: %s", err.Error())
	}

	sm := new(onlinePb.ReqSignMetrics)
	sm.Payload = payload
	sm.Signature = signature
	sm.PublicKey = publicKey
	return sm, errs, nil
}

func (dc *dcWrap) getPayloadOnline(btfsNode *core.IpfsNode) ([]byte, error) {
	var lastSignedInfo *onlinePb.SignedInfo
	var lastSignature string
	var lastTime time.Time

	lastOnline, err := chain.GetLastOnline()
	if err != nil {
		return nil, err
	}

	if lastOnline == nil {
		lastSignedInfo = nil
		lastSignature = ""
	} else {
		lastSignedInfo = &onlinePb.SignedInfo{
			Peer:        lastOnline.LastSignedInfo.Peer,
			CreatedTime: lastOnline.LastSignedInfo.CreatedTime,
			Version:     lastOnline.LastSignedInfo.Version,
			Nonce:       lastOnline.LastSignedInfo.Nonce,
			BttcAddress: lastOnline.LastSignedInfo.BttcAddress,
			SignedTime:  lastOnline.LastSignedInfo.SignedTime,
		}
		lastSignature = lastOnline.LastSignature
		lastTime = lastOnline.LastTime
	}

	pn := &onlinePb.PayLoadInfo{
		NodeId:         btfsNode.Identity.Pretty(),
		Node:           dc.pn,
		LastSignedInfo: lastSignedInfo,
		LastSignature:  lastSignature,
		LastTime:       lastTime,
	}
	bytes, err := proto.Marshal(pn)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("--- online, req,  LastSignedInfo:%+v, LastSignature:%+v, pn.LastTime:%+v \n", pn.LastSignedInfo, pn.LastSignature, pn.LastTime)

	return bytes, nil
}

func (dc *dcWrap) collectionAgentOnline(node *core.IpfsNode) {
	tick := time.NewTicker(heartBeatOnline)
	//tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()

	// Force tick on immediate start
	// make the configuration available in the for loop
	for ; true; <-tick.C {
		cfg, err := dc.node.Repo.Config()
		if err != nil {
			continue
		}

		if isReportOnlineEnabled(cfg) {
			//fmt.Println("")
			//fmt.Println("--- online agent ---")

			dc.SendDataOnline(node, cfg)
		}
	}
}
