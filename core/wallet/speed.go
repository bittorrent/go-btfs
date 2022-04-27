package wallet

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tron-us/go-btfs-common/crypto"
)

const (
	portFileName = "port"
	api          = "http://127.0.0.1:%d/api"
	keyUrl       = api + "/private_key?pw=%s&t=%s"
	tokenUrl     = api + "/token"
)

var (
	portFile = filepath.Join(portPath, portFileName)
)

// return speed key in base64
func DiscoverySpeedKey(password string) (string, error) {
	password = url.QueryEscape(password)
	if err := validateOs(); err != nil {
		return "", err
	}
	pf, err := os.Open(portFile)
	if err != nil {
		return "", err
	}
	port, err := readPort(pf)
	if err != nil {
		return "", err
	}
	token, err := get(fmt.Sprintf(tokenUrl, port))
	if err != nil {
		return "", err
	}
	key, err := get(fmt.Sprintf(keyUrl, port, password, token))
	if err != nil {
		return "", err
	}
	if key == "" {
		return "", errors.New("invalid private key")
	}
	base64, err := crypto.Hex64ToBase64(key)
	if err != nil {
		return "", err
	}
	return base64, nil
}

func readPort(r io.Reader) (int64, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return -1, err
	}
	return strconv.ParseInt(strings.TrimSpace(string(bytes)), 10, 32)
}

func get(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
