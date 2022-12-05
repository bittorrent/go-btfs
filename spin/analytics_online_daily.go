package spin

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	config "github.com/TRON-US/go-btfs-config"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core"
	onlinePb "github.com/tron-us/go-btfs-common/protos/online"
	cgrpc "github.com/tron-us/go-btfs-common/utils/grpc"

	"github.com/cenkalti/backoff/v4"
)

func (dc *dcWrap) doSendOnlineDaily(ctx context.Context, config *config.Config, sm *onlinePb.ReqSignMetrics) error {
	onlineService := config.Services.OnlineServerDomain
	if len(onlineService) <= 0 {
		onlineService = chain.GetOnlineServer(config.ChainInfo.ChainId)
	}
	cb := cgrpc.OnlineClient(onlineService)
	return cb.WithContext(ctx, func(ctx context.Context, client onlinePb.OnlineServiceClient) error {
		resp, err := client.DoDailyStatusReport(ctx, sm)
		if err != nil {
			chain.CodeStatus = chain.ConstCodeError
			chain.ErrStatus = err
			return err
		} else {
			chain.CodeStatus = chain.ConstCodeSuccess
			chain.ErrStatus = nil
		}

		//fmt.Printf("--- online, resp, SignedInfo = %+v, signature = %+v \n", resp.SignedInfo, resp.Signature)
		if resp.Code != onlinePb.ResponseCode_SUCCESS {
			return errors.New("DoDailySignReportHandler err: " + resp.Message)
		}

		lastOnline, err := chain.GetLastOnline()
		if err != nil {
			return err
		}

		if lastOnline != nil {
			_, err = chain.SetReportOnlineLastTimeDailyOK()
			if err != nil {
				return err
			}

			_, err = chain.SetReportOnlineListDailyOK(lastOnline)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (dc *dcWrap) SendOnlineDaily(node *core.IpfsNode, config *config.Config) {
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
		fmt.Printf("--- online 1, doSendDataOnline \n")
		err := dc.doSendOnlineDaily(node.Context(), config, sm)
		if err != nil {
			fmt.Printf("--- online 2, doSendDataOnline error = %+v \n", err)
			log.Infof("failed： doSendDataOnline to online server: %+v ", err)
		} else {
			log.Debug("sent OK, doSendDataOnline to online server")
		}
		return err
	}, bo)
}

func (dc *dcWrap) collectionAgentOnlineDaily(node *core.IpfsNode) {
	//tick := time.NewTicker(dailyReportOnline)
	tick := time.NewTicker(10 * time.Second)
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

			dc.SendOnlineDaily(node, cfg)
		}
	}
}

func GetLastOnlineInfoWhenNodeMigration(ctx context.Context, config *config.Config) error {
	onlineService := config.Services.OnlineServerDomain
	if len(onlineService) <= 0 {
		onlineService = chain.GetOnlineServer(config.ChainInfo.ChainId)
	}
	cb := cgrpc.OnlineClient(onlineService)
	return cb.WithContext(ctx, func(ctx context.Context, client onlinePb.OnlineServiceClient) error {
		req := onlinePb.ReqLastDailySignedInfo{
			PeerId: config.Identity.PeerID,
		}

		resp, err := client.GetLastDailySignedInfo(ctx, &req)
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
			err := chain.StoreOnline(&onlineInfo)
			if err != nil {
				return err
			}
		}

		return err
	})
}
