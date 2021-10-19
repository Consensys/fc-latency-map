package addresses

import (
	"fmt"
	"net"
	"testing"

	"github.com/filecoin-project/go-state-types/abi"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
)

var dummyIpAddressV4 = "127.0.0.1"
var dummyIpAddressV6 = "::1"
var dummyPort = 1234
var dummyMultiAddressV4, _ = ma.NewMultiaddr("/ip4/" + dummyIpAddressV4 + "/tcp/" + fmt.Sprintf("%d", dummyPort))
var dummyMultiAddressV6, _ = ma.NewMultiaddr("/ip6/" + dummyIpAddressV6 + "/tcp/" + fmt.Sprintf("%d", dummyPort))

func Test_MultiAddrs_Nil(t *testing.T) {
	addrs := MultiAddrs(nil)
	assert.Nil(t, addrs)
}

func Test_MultiAddrs_Empty(t *testing.T) {
	addrs := MultiAddrs([]abi.Multiaddrs{})
	assert.Nil(t, addrs)
}

func Test_MultiAddrs_OK(t *testing.T) {
	// Arrange
	dummyMultiAddresses := [][]byte{dummyMultiAddressV4.Bytes(), dummyMultiAddressV6.Bytes()}

	// Act
	addrs := MultiAddrs(dummyMultiAddresses)

	// Assert
	assert.NotNil(t, addrs)
	assert.NotEmpty(t, addrs)
}

func Test_IPAddress_OK(t *testing.T) {
	// Arrange
	dummyMultiAddresses := []ma.Multiaddr{dummyMultiAddressV4, dummyMultiAddressV6}

	// Act
	ips, port := IPAddress(dummyMultiAddresses)

	// Assert
	assert.NotNil(t, ips)
	assert.NotEmpty(t, ips)
	assert.Equal(t, dummyIpAddressV4, ips[0])
	assert.Equal(t, dummyIpAddressV6, ips[1])

	assert.NotNil(t, port)
	assert.NotEmpty(t, port)
	assert.Equal(t, dummyPort, port)
}

func Test_GetIPVersion_OK(t *testing.T) {
	// Arrange
	dummyIpV4 := net.ParseIP(dummyIpAddressV4)
	dummyIpV6 := net.ParseIP(dummyIpAddressV6)

	// Act
	v4 := GetIPVersion(dummyIpV4)
	v6 := GetIPVersion(dummyIpV6)

	// Assert
	assert.NotNil(t, v4)
	assert.NotNil(t, v6)
	assert.NotNil(t, 4, v4)
	assert.NotNil(t, 6, v6)
}
