package resolve_test

import (
	"testing"

	coremock "github.com/bittorrent/go-btfs/core/mock"
	"github.com/bittorrent/go-btfs/namesys/resolve"

	path "github.com/ipfs/go-path"
)

func TestResolveNoComponents(t *testing.T) {
	n, err := coremock.NewMockNode()
	if n == nil || err != nil {
		t.Fatal("Should have constructed a mock node", err)
	}

	_, err = resolve.Resolve(n.Context(), n.Namesys, n.Resolver, path.Path("/btns/"))
	if err.Error() != "invalid path \"/btns/\": btns path missing BTNS ID" {
		t.Error("Should error with no components (/ipns/).", err)
	}

	_, err = resolve.Resolve(n.Context(), n.Namesys, n.Resolver, path.Path("/btfs/"))
	if err.Error() != "invalid path \"/btfs/\": not enough path components" {
		t.Error("Should error with no components (/ipfs/).", err)
	}

	_, err = resolve.Resolve(n.Context(), n.Namesys, n.Resolver, path.Path("/../.."))
	if err.Error() != "invalid path \"/../..\": unknown namespace \"..\"" {
		t.Error("Should error with invalid path.", err)
	}
}
