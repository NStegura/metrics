package ip

import (
	"fmt"
	"net"
)

// GetIP get first not local ip.
func GetIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip, fmt.Errorf("failed to get ip's: %w", err)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}
	return "", nil
}
