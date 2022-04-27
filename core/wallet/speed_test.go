package wallet

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/tron-us/go-btfs-common/crypto"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestReadPort(t *testing.T) {
	switch runtime.GOOS {
	case "darwin":
		e, err := homedir.Expand("~/Library/Application Support/uTorrent Web/BitTorrentHelper/")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, e, portPath)
	case "windows":
		assert.Equal(t, "%AppData%/../Local/BitTorrentHelper/", portPath)
	default:
		assert.Equal(t, "", portPath)
	}
	port, err := readPort(strings.NewReader("\n8888\r\n "))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(8888), port)
}

func TestGetPlainKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		switch r.URL.EscapedPath() {
		case "/api/private_key":
			w.Write([]byte("c14f99e28b64abfb743a88002939f776b7f9ebab9aeac5cb7340daf7be81c2a1"))
		case "/api/token":
			w.Write([]byte("token1"))
		}
		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}
		if r.URL.EscapedPath() != "/api/private_key" && r.URL.EscapedPath() != "/api/token" {
			t.Errorf("Expected request to '/api/private_key' or '/api/token', got '%s'", r.URL.EscapedPath())
		}
	}))
	defer ts.Close()
	token, err := get(ts.URL + "/api/token")
	assert.Equal(t, "token1", token)
	key, err := get(ts.URL + "/api/private_key")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "c14f99e28b64abfb743a88002939f776b7f9ebab9aeac5cb7340daf7be81c2a1", key)
}

func TestHexToBase64(t *testing.T) {
	base64, err := crypto.Hex64ToBase64("c14f99e28b64abfb743a88002939f776b7f9ebab9aeac5cb7340daf7be81c2a1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "CAISIMFPmeKLZKv7dDqIACk593a3+eurmurFy3NA2ve+gcKh", base64)
}
