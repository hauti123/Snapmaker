package snapmaker

import (
	"fmt"
	"net"
	"time"
)

func DiscoverSnapmaker(timeout time.Duration) (string, error) {
	err := sendDiscoveryPacket()
	if err != nil {
		return "", err
	}
	return waitForSnapmakerResponse(timeout)
}

func sendDiscoveryPacket() error {
	fmt.Println("starting discovery")
	pc, err := net.ListenPacket("udp4", ":50000")
	if err != nil {
		return err
	}
	defer pc.Close()

	addr, err := net.ResolveUDPAddr("udp4", "192.168.188.255:20054")
	if err != nil {
		return err
	}

	fmt.Println("sending discovery packet")
	_, err = pc.WriteTo([]byte("discover"), addr)
	if err != nil {
		return err
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

/*
func getBroadcastIp() (string, error) {
	networkInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, networkInterface := range networkInterfaces {
		addresses, err := networkInterface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// process IP address
		}
	}
}
*/
