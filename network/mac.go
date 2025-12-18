package network

import (
	"fmt"
	"net"
)

// GetMacAddresses returns a list of MAC addresses for all active network interfaces
func GetMacAddresses() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var macs []string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		if iface.HardwareAddr == nil {
			continue
		}

		macs = append(macs, iface.HardwareAddr.String())
	}

	if len(macs) == 0 {
		return nil, fmt.Errorf("no mac address found")
	}

	return macs, nil
}
