package system

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// Define reliable providers for fetching public IPs.
var ipProviders = []string{
	"https://api.ipify.org?format=text",
	"https://ifconfig.me/ip",
	"http://checkip.amazonaws.com",
	"https://icanhazip.com",
}

/*
GetPublicIP attempts to fetch the public IP address of the server.
It sequentially requests external APIs and falls back to the local outbound IP if all fail.
*/
func GetPublicIP() (string, error) {
	// 1. Attempt to fetch real public IP via external APIs.
	// Set a 3-second timeout to prevent blocking.
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	for _, url := range ipProviders {
		ip, err := fetchIP(client, url)
		if err == nil && isValidIP(ip) {
			fmt.Printf("✅ Get the server IP: %s\n", ip)
			return ip, nil
		}
	}

	// 2. If all external APIs fail, fallback to local outbound IP (LAN IP).
	fmt.Println("Warning: Could not fetch public IP from external APIs, falling back to local IP.")
	return GetLocalOutboundIP()
}

// fetchIP performs an HTTP request to get the IP.
func fetchIP(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}

/*
GetLocalOutboundIP obtains the preferred local outbound IP address.

Using the UDP Dial technique (without actually sending packets), this is the most accurate way to get the local LAN IP of the machine.
*/
func GetLocalOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP.String()
	fmt.Printf("✅ Get the local outbound IP (may not work for WireGuard): %s\n", ip)
	return ip, nil
}

// isValidIP validate the form of the ip string.
func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}
