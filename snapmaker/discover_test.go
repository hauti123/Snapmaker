package snapmaker

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBroadcastIpDiscovery(t *testing.T) {

	networkInterfaceAddresses, err := net.InterfaceAddrs()
	require.NoError(t, err)
	assert.NoError(t, err)

	for _, address := range networkInterfaceAddresses {

		var ipAddress net.IP
		switch typedAddress := address.(type) {
		case *net.TCPAddr:
			ipAddress = typedAddress.IP.To4()
		case *net.UDPAddr:
			ipAddress = typedAddress.IP.To4()
		case *net.IPNet:
			ipAddress = typedAddress.IP.To4()

		default:
			assert.Fail(t, "unknown address type %v", typedAddress)
		}

		ipAddress.IsLinkLocalUnicast()
		if ipAddress != nil && len(ipAddress) == net.IPv4len && ipAddress.IsGlobalUnicast() {
			ipAddress[3] = 255
			fmt.Printf("%s\n", ipAddress.String())

		}
	}
}
