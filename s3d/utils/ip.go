package utils

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/yann-y/fds/internal/iam/set"
	"net"
	"runtime"
)

var log = logging.Logger("utils")

// MustGetLocalIP4 returns IPv4 addresses of localhost.  It panics on error.
func MustGetLocalIP4() (ipList set.StringSet) {
	ipList = set.NewStringSet()
	ifs, err := net.Interfaces()
	if err != nil {
		log.Errorf("Unable to get IP addresses of this host %v", err)

	}

	for _, interf := range ifs {
		addrs, err := interf.Addrs()
		if err != nil {
			continue
		}
		if runtime.GOOS == "windows" && interf.Flags&net.FlagUp == 0 {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip.To4() != nil {
				ipList.Add(ip.String())
			}
		}
	}

	return ipList
}
