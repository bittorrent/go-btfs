package spin

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	onlinePb "github.com/bittorrent/go-btfs-common/protos/online"
	cgrpc "github.com/bittorrent/go-btfs-common/utils/grpc"
	config "github.com/bittorrent/go-btfs-config"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core"

	"github.com/cenkalti/backoff/v4"
)

var (
	startTime = time.Now()
)

func (dc *dcWrap) doSendOnlineDaily(ctx context.Context, config *config.Config, sm *onlinePb.ReqSignMetrics) (msg string, err error) {
	onlineService := config.Services.OnlineServerDomain
	if len(onlineService) <= 0 {
		onlineService = chain.GetOnlineServer(config.ChainInfo.ChainId)
	}
	cb := cgrpc.OnlineClient(onlineService)
	err = cb.WithContext(ctx, func(ctx context.Context, client onlinePb.OnlineServiceClient) error {
		resp, err := client.DoDailyStatusReport(ctx, sm)
		if err != nil {
			fmt.Printf("daily report online, resp = %+v, err = %+v \n", resp, err)
			chain.CodeStatus = chain.ConstCodeError
			chain.ErrStatus = err
			return err
		} else {
			chain.CodeStatus = chain.ConstCodeSuccess
			chain.ErrStatus = nil
		}

		if resp != nil && len(resp.Message) > 0 {
			msg = resp.Message
		}

		//return errors.New("xxx") //test err
		if resp.Code != onlinePb.ResponseCode_SUCCESS {
			if resp.Message == "to many request" {
				return nil
			}
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

	return msg, err
}

func (dc *dcWrap) SendOnlineDaily(node *core.IpfsNode, config *config.Config) (msg string, err error) {
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
		return "", err
	}

	// retry: max 3 times
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 3 * 60 * time.Second
	bo.InitialInterval = 60 * time.Second
	bo.MaxInterval = 60 * time.Second
	backoff.Retry(func() error {
		msg, err = dc.doSendOnlineDaily(node.Context(), config, sm)
		if err != nil {
			log.Infof("failed： doSendDataOnline to online server: %+v ", err)
		} else {
			log.Debug("sent OK, doSendDataOnline to online server")
		}

		return err
	}, bo)

	return msg, err
}

func (dc *dcWrap) collectionAgentOnlineDaily(node *core.IpfsNode) {
	tick := time.NewTicker(interReportOnlineDaily)
	//tick := time.NewTicker(60 * time.Second)
	defer tick.Stop()

	// Force tick on immediate start
	// make the configuration available in the for loop
	for ; true; <-tick.C {
		cfg, err := dc.node.Repo.Config()
		if err != nil {
			continue
		}
		if !isReportOnlineEnabled(cfg) {
			return
		}

		report, err := chain.GetReportOnlineLastTimeDaily()
		//fmt.Printf("... GetReportOnlineLastTimeDaily, report: %+v err:%+v \n", report, err)
		if err != nil {
			log.Errorf("GetReportOnlineLastTimeDaily err:%+v", err)
			continue
		}

		now := time.Now()
		if now.Sub(startTime) < 30*time.Minute {
			continue
		}

		nowMod := now.Unix() % 86400
		// report only 1 hour every, and must after 10 hour.
		if nowMod > report.EveryDaySeconds &&
			nowMod < report.EveryDaySeconds+3600*2 &&
			now.Sub(report.LastReportTime) > 10*time.Hour {

			fmt.Printf("every day, SendOnlineDaily, time:%+v\n", time.Now().String())
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
