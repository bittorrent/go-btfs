package bindata

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ip2location/ip2location-go/v9"
	"github.com/multiformats/go-multiaddr"
)

var DB *ip2location.DB

func Init() error {
	filename := "data/IP2LOCATION-LITE-DB1.IPV6.BIN"
	var err error
	DB, err = InitDB(filename)
	if err != nil {
		return err
	}
	return nil
}

// InitDB init ip2locaiton db with bin data file filename
func InitDB(filename string) (*ip2location.DB, error) {
	data, err := Asset(filename)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	return ip2location.OpenDBWithReader(NopCloser(reader))
}

type nopCloser struct {
	*bytes.Reader
}

func NopCloser(r *bytes.Reader) nopCloser {
	return nopCloser{r}
}

func (nc nopCloser) Close() error {
	return nil
}

// CountryShortCodeByIP4 return country short by ipv4 address
func CountryShortCodeByIP(addr string) (string, error) {
	if DB == nil {
		return "", errors.New("ip2location db not initialized")
	}
	location, err := DB.Get_country_short(addr)
	if err != nil {
		return "", err
	}

	return location.Country_short, nil
}

func CountryShortCode(addr multiaddr.Multiaddr) (string, error) {
	ipv4, err := addr.ValueForProtocol(multiaddr.P_IP4)
	if err != nil {
		fmt.Printf("get ipv4 err:%+v", err)
	} else {
		return CountryShortCodeByIP(ipv4)
	}

	ipv6, err := addr.ValueForProtocol(multiaddr.P_IP6)
	if err != nil {
		fmt.Printf("get ipv6 err:%+v", err)
		return "", err
	}
	return CountryShortCodeByIP(ipv6)
}
