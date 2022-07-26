package snapmaker

import (
	"fmt"
	"net"
	"time"
)

func DiscoverSnapmaker(timeout time.Duration, apiToken string) (Snapmaker, error) {

	startTime := time.Now()
	for time.Now().Sub(startTime) < timeout {
		err := sendDiscoveryPacket()
		if err != nil {
			continue
		}
		ip, err := waitForSnapmakerResponse(timeout)

		if err != nil {
			return Snapmaker{}, err
		}

		return Snapmaker{
			ipAdress: ip,
			port:     snapmakerApiPort,
			token:    apiToken,
		}, nil
	}
	return Snapmaker{}, fmt.Errorf("discovery timeout reached")
}

func sendDiscoveryPacket() error {
	fmt.Println("starting discovery")
	pc, err := net.ListenPacket("udp4", ":50000")
	if err != nil {
		return err
	}
	defer pc.Close()

	broadcastAddresses, err := getBroadcastIp()
	if err != nil {
		return err
	}

	for _, broadcastAddress := range broadcastAddresses {
		addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", broadcastAddress, snapmakerDiscoveryPort))
		if err != nil {
			return err
		}

		fmt.Printf("sending discovery packet to %s\n", addr.String())
		_, err = pc.WriteTo([]byte(snapmakerDiscoveryPayload), addr)
		if err != nil {
			return err
		}
	}
	return nil
}

func waitForSnapmakerResponse(timeout time.Duration) (string, error) {
	fmt.Println("waiting for response")
	pc, err := net.ListenPacket("udp4", ":50000")
	if err != nil {
		return "", err
	}
	defer pc.Close()

	buf := make([]byte, 1024)

	pc.SetReadDeadline(time.Now().Add(timeout))
	_, addr, err := pc.ReadFrom(buf)
	if err != nil {
		return "", err
	}

	host, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		return "", err
	}

	return host, nil
}

func getBroadcastIp() ([]string, error) {

	networkInterfaceAddresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	broadcastAddresses := make([]string, 0)
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
		}

		if ipAddress != nil && len(ipAddress) == net.IPv4len && ipAddress.IsGlobalUnicast() {
			ipAddress[3] = 255
			broadcastAddresses = append(broadcastAddresses, ipAddress.String())
		}
	}
	return broadcastAddresses, nil
}
