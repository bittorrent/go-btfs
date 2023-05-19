package resolve

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bittorrent/go-btfs/namesys"

	logging "github.com/ipfs/go-log"
	"github.com/ipfs/go-path"
)

var log = logging.Logger("nsresolv")

// ErrNoNamesys is an explicit error for when an BTFS node doesn't
// (yet) have a name system
var ErrNoNamesys = errors.New(
	"core/resolve: no Namesys on BtfsNode - can't resolve btns entry")

// ResolveIPNS resolves /btns paths
func ResolveIPNS(ctx context.Context, nsys namesys.NameSystem, p path.Path) (path.Path, error) {
	if strings.HasPrefix(p.String(), "/btns/") {
		evt := log.EventBegin(ctx, "resolveBtnsPath")
		defer evt.Done()
		// resolve btns paths

		// TODO(cryptix): we should be able to query the local cache for the path
		if nsys == nil {
			evt.Append(logging.LoggableMap{"error": ErrNoNamesys.Error()})
			return "", ErrNoNamesys
		}

		seg := p.Segments()

		if len(seg) < 2 || seg[1] == "" { // just "/<protocol/>" without further segments
			err := fmt.Errorf("invalid path %q: btns path missing BTNS ID", p)
			evt.Append(logging.LoggableMap{"error": err})
			return "", err
		}

		extensions := seg[2:]
		resolvable, err := path.FromSegments("/", seg[0], seg[1])
		if err != nil {
			evt.Append(logging.LoggableMap{"error": err.Error()})
			return "", err
		}

		respath, err := nsys.Resolve(ctx, resolvable.String())
		if err != nil {
			evt.Append(logging.LoggableMap{"error": err.Error()})
			return "", err
		}

		segments := append(respath.Segments(), extensions...)
		p, err = path.FromSegments("/", segments...)
		if err != nil {
			evt.Append(logging.LoggableMap{"error": err.Error()})
			return "", err
		}
	}
	return p, nil
}

// Resolve resolves the given path by parsing out protocol-specific
// entries (e.g. /btns/<node-key>) and then going through the /btfs/
// entries and returning the final node.
/* func Resolve(ctx context.Context, nsys namesys.NameSystem, r *resolver.Resolver, p path.Path) (format.Node, error) {
	p, err := ResolveIPNS(ctx, nsys, p)
	if err != nil {
		return nil, err
	}

	// ok, we have an BTFS path now (or what we'll treat as one)
	return r.ResolvePath(ctx, p)
}
*/
