package addresses

import (
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

func IpAddress(a []ma.Multiaddr) []string {
	var ips []string
	for _, v := range a {
		if ip, err := v.ValueForProtocol(ma.P_IP4); err == nil {
			ips = append(ips, ip)
		} else if ip, err := v.ValueForProtocol(ma.P_IP6); err == nil {
			ips = append(ips, ip)
		}
	}
	return ips
}
