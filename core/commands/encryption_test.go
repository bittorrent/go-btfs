package commands

import (
	"bytes"
	"context"
	"errors"
	"io"
	"math"
	"testing"
	"time"

	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/bootstrap"
	"github.com/bittorrent/go-btfs/core/coreapi"
	coremock "github.com/bittorrent/go-btfs/core/mock"

	files "github.com/bittorrent/go-btfs-files"
	"github.com/bittorrent/interface-go-btfs-core/options"
	"github.com/libp2p/go-libp2p/core/peer"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/libp2p/go-testutil"
)

func TestEncryption(t *testing.T) {
	addOpts := []options.UnixfsAddOption{
		options.Unixfs.Encrypt(true),
		options.Unixfs.PeerId("16Uiu2HAmRih4otzcxyZ428QoQim8SptHZq5sjBxiUuiKr7ctPmMG"),
	}
	getOpts := []options.UnixfsGetOption{
		options.Unixfs.Decrypt(true),
		options.Unixfs.PrivateKey("CAISIBtp+e228gq2SBTP/bfXbnUx+OQWZNlDuEntq7eOPlBB"),
	}

	msg := "btt to da moon"
	conf := testutil.LatencyConfig{NetworkLatency: 400 * time.Millisecond}
	if err := directAddCat([]byte(msg), conf, addOpts, getOpts); err != nil {
		t.Fatal(err)
	}
}

func directAddCat(data []byte, conf testutil.LatencyConfig, addOpts []options.UnixfsAddOption,
	getOpts []options.UnixfsGetOption) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create network
	mn := mocknet.New()
	mn.SetLinkDefaults(mocknet.LinkOptions{
		Latency: conf.NetworkLatency,
		// TODO add to conf. This is tricky because we want 0 values to be functional.
		Bandwidth: math.MaxInt32,
	})

	adder, err := core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Host:   coremock.MockHostOption(mn),
	})
	if err != nil {
		return err
	}
	defer adder.Close()

	catter, err := core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Host:   coremock.MockHostOption(mn),
	})
	if err != nil {
		return err
	}
	defer catter.Close()

	adderApi, err := coreapi.NewCoreAPI(adder)
	if err != nil {
		return err
	}

	catterApi, err := coreapi.NewCoreAPI(catter)
	if err != nil {
		return err
	}

	err = mn.LinkAll()
	if err != nil {
		return err
	}

	bs1 := []peer.AddrInfo{adder.Peerstore.PeerInfo(adder.Identity)}
	bs2 := []peer.AddrInfo{catter.Peerstore.PeerInfo(catter.Identity)}

	if err := catter.Bootstrap(bootstrap.BootstrapConfigWithPeers(bs1)); err != nil {
		return err
	}
	if err := adder.Bootstrap(bootstrap.BootstrapConfigWithPeers(bs2)); err != nil {
		return err
	}

	added, err := adderApi.Unixfs().Add(ctx, files.NewBytesFile(data), addOpts...)
	if err != nil {
		return err
	}

	readerCatted, err := catterApi.Unixfs().Get(ctx, added, getOpts...)
	if err != nil {
		return err
	}

	// verify
	var bufout bytes.Buffer
	_, err = io.Copy(&bufout, readerCatted.(io.Reader))
	if err != nil {
		return err
	}
	if !bytes.Equal(bufout.Bytes(), data) {
		return errors.New("catted data does not match added data")
	}

	return nil
}
