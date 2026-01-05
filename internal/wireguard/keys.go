package wireguard

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	configDir = "/etc/wireguard"
)

/*
GenerateWGKeys generates a WireGuard private and public key pair and saves them to files.

Returns the private key, public key, and an error if any operation fails.
*/
func GenerateWGKeys(interfaceName string) (string, string, error) {
	priKeyPath := filepath.Join(configDir, fmt.Sprintf("%s.key", interfaceName))
	pubKeyPath := filepath.Join(configDir, fmt.Sprintf("%s.pub", interfaceName))

	cmdGenKey := exec.Command("wg", "genkey")
	priKeyBytes, err := cmdGenKey.Output()
	if err != nil {
		return "", "", fmt.Errorf("Failed to generate private key for wireguard: %w", err)
	}
	// Write the bytes into file
	if err := os.WriteFile(priKeyPath, priKeyBytes, 0600); err != nil {
		return "", "", fmt.Errorf("Failed to write private key into %s: %w", priKeyPath, err)
	}

	cmdPubKey := exec.Command("wg", "pubkey")
	cmdPubKey.Stdin = bytes.NewReader(priKeyBytes)
	pubKeyBytes, err := cmdPubKey.Output()
	if err != nil {
		return "", "", fmt.Errorf("Failed to generate public key for wireguard: %w", err)
	}
	// Write the public key into file
	if err := os.WriteFile(pubKeyPath, pubKeyBytes, 0600); err != nil {
		return "", "", fmt.Errorf("Failed to write public key into %s: %w", pubKeyPath, err)
	}

	// Return the keys as strings
	privateKey := strings.TrimSpace(string(priKeyBytes))
	publicKey := strings.TrimSpace(string(pubKeyBytes))
	return privateKey, publicKey, nil
}
