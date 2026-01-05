package system

import (
	"errors"
	"fast-wireguard/internal/tracker"
	"fmt"
	"net"
	"strings"
)

/*
GetPhysicalInterface attempts to find the default physical network interface.

It uses a "dial" trick to find the outbound interface, then filters it against
loopback and managed interfaces.
*/
func GetPhysicalInterface() (string, error) {
	// Strategy 1: Try to set up one UDP connection (won't send data) to search for the route list
	// It is more accuary than purely enumerating the interface list.
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err == nil {
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)

		// Search the name of the interface based on the IP adress
		iface, err := getInterfaceByIP(localAddr.IP)
		if err == nil {
			// Validate the name (low probability, for safety)
			isManaged, _ := tracker.IsManagedByUs(iface.Name)
			if !isManaged {
				fmt.Printf("✅ Get the physical interface: %s\n", iface.Name)
				return iface.Name, nil
			}
		}
	}

	// Strategy 2: If there is no network or the method above does not work, fall back to search all the interfaces
	// Search for the first interface Up, not Loopback and not managed by us
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to list interfaces: %w", err)
	}

	for _, i := range ifaces {
		// 1. Ignore the interfaces that do not work
		if i.Flags&net.FlagUp == 0 {
			continue
		}
		// 2. Ignore Loopback interfaces
		if i.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 3. Ignore common prefixes of virtual interfaces (optional, improves robustness)
		if strings.HasPrefix(i.Name, "docker") || strings.HasPrefix(i.Name, "veth") {
			continue
		}

		// 4. Core logic: Check whether this interface is managed by us
		isManaged, err := tracker.IsManagedByUs(i.Name)
		if err != nil {
			// If there is an error reading the log, we can choose to skip it or just print a warning
			// Here we choose a conservative approach: if there is an error, assume it is not a physical network card and skip it
			fmt.Printf("Warning: checking managed status for %s failed: %v\n", i.Name, err)
			continue
		}
		if isManaged {
			continue
		}

		// Found a interface that meets the requirements
		fmt.Printf("✅ Get the physical interface: %s\n", i.Name)
		return i.Name, nil
	}

	return "", errors.New("no suitable physical interface found")
}

/*
getInterfaceByIP finds the network interface object corresponding to the IP address.
*/
func getInterfaceByIP(ip net.IP) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var currentIP net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				currentIP = v.IP
			case *net.IPAddr:
				currentIP = v.IP
			}
			if currentIP.Equal(ip) {
				return &i, nil
			}
		}
	}
	return nil, errors.New("interface not found for IP")
}
