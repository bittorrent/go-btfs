package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/bittorrent/go-btfs/core"

	iface "github.com/TRON-US/interface-go-btfs-core"

	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/peer"
)

type P2PRemoteCall struct {
	Node *core.IpfsNode
	ID   peer.ID
}

type ErrorMessage struct {
	Message string
	Code    int
	Type    string
}

const P2PRemoteCallProto = "/rapi"

// P2PCall is a wrapper for creating a client and calling a get
// If passed a nil context, a new one will be created
func P2PCall(ctx context.Context, n *core.IpfsNode, coreApi iface.CoreAPI, pid peer.ID, api string, args ...interface{}) ([]byte, error) {
	// new context if not caller-passed down
	if ctx == nil {
		ctx = context.Background()
	}
	err := coreApi.Swarm().Connect(ctx, peer.AddrInfo{
		ID: pid,
	})
	if err != nil {
		return nil, err
	}
	remoteCall := &P2PRemoteCall{
		Node: n,
		ID:   pid,
	}
	return remoteCall.CallGet(ctx, api, args)
}

// P2PCallStrings is a helper to pass string arguments to P2PCall
func P2PCallStrings(ctx context.Context, n *core.IpfsNode, coreApi iface.CoreAPI, pid peer.ID, api string, strs ...string) ([]byte, error) {
	var args []interface{}
	for _, str := range strs {
		args = append(args, str)
	}
	return P2PCall(ctx, n, coreApi, pid, api, args...)
}

func (r *P2PRemoteCall) CallGet(ctx context.Context, api string, args []interface{}) ([]byte, error) {
	var sb strings.Builder
	for i, arg := range args {
		if i == 0 {
			sb.WriteString("?")
		} else {
			sb.WriteString("&")
		}
		sb.WriteString("arg=")
		switch arg.(type) {
		case []byte:
			s := url.QueryEscape(string(arg.([]byte)))
			sb.WriteString(s)
		case string:
			sb.WriteString(arg.(string))
		default:
			sb.WriteString(fmt.Sprintf("%v", arg))
		}
	}
	// setup url
	reqUrl := fmt.Sprintf("libp2p://%s%s%s%s", r.ID.Pretty(), apiPrefix, api, sb.String())
	// perform context setup
	req, err := http.NewRequestWithContext(ctx, "POST", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	// libp2p protocol register
	tr := &http.Transport{}
	tr.RegisterProtocol("libp2p",
		p2phttp.NewTransport(r.Node.PeerHost, p2phttp.ProtocolOption(P2PRemoteCallProto)))
	client := &http.Client{Transport: tr}
	// call
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		var tmp IoError = fmt.Errorf("fail to read response body: %s", err)
		return nil, tmp
	}
	if resp.StatusCode != http.StatusOK {
		e := &ErrorMessage{}
		if err = json.Unmarshal(body, e); err != nil {
			return nil, err
		}
		var tmp BusinessError = fmt.Errorf(e.Message)
		return nil, tmp
	}
	return body, nil
}

type IoError error
type BusinessError error

func (r *P2PRemoteCall) CallPost() {
}
