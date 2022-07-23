package snapmaker

import (
	"fmt"
	"net"
	"testing"
)

func TestBroadcastIpDiscovery(t *testing.T) {

	networkInterfaces, err := net.Interfaces()
	assert.NoError(t, err)
	fmt.Printf("Interfaces:\n%v\n", networkInterfaces)

	for _, networkInterface := range networkInterfaces {
		fmt.Printf("%v\n", networkInterface)
	}
}
