package addresses

import (
	"net"
	"strconv"

	"github.com/filecoin-project/go-state-types/abi"
	ma "github.com/multiformats/go-multiaddr"
)

func MultiAddrs(addr []abi.Multiaddrs) []ma.Multiaddr {
	var m []ma.Multiaddr
	for _, v := range addr {
		if a, err := ma.NewMultiaddrBytes(v); err == nil {
			m = append(m, a)
		}
	}

	return m
}

func IPAddress(a []ma.Multiaddr) (ips []string, tcpPort int) {
	for _, v := range a {
		if port, err := v.ValueForProtocol(ma.P_TCP); err == nil {
			tcpPort, _ = strconv.Atoi(port)
		}
		if ip, err := v.ValueForProtocol(ma.P_IP4); err == nil {
			ips = append(ips, ip)
		}
		if ip, err := v.ValueForProtocol(ma.P_IP6); err == nil {
			ips = append(ips, ip)
		}
	}

	return ips, tcpPort
}

func GetIPVersion(ipAdd net.IP) int {
	af := 4
	if ipAdd.To4() == nil {
		af = 6
	}

	return af
}
