package wireguard

import (
	"bytes"
	_ "embed"
	"fast-wireguard/internal/templates"
	"fast-wireguard/internal/tracker"
	"fast-wireguard/pkg/utils"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	wgConfigDir = "/etc/wireguard"
)

type WgConfTplData struct {
	InterfaceName     string
	PriKeyServer      string
	ListenPort        int
	Address           string
	MTU               int
	PhysicalInterface string
}

type PeerConfTplData struct {
	PeerName     string
	PubKeyClient string
	AllowedIPs   string
}

type ClientConfTplData struct {
	PriKeyClient string
	AllowedIPs   string
	PubKeyServer string
	Endpoint     string
	MTU          int
}

/*
GenerateWGConfig create the service file for the given parameters
*/
func GenerateWGConfig(
	interfaceName string,
	listenPort int,
	priKeyServer string,
	mtu int,
	IPAdressLocal string,
	physicalInterface string,
	force bool,
) error {
	// 1. Make sure the path of the configuration file
	configPath := filepath.Join(wgConfigDir, fmt.Sprintf("%s.conf", interfaceName))

	// 2. Handle the option "force"
	if !force {
		if _, err := os.Stat(configPath); err == nil {
			confirmed := utils.PromptConfirm(fmt.Sprintf("Config file for %s already exists. Do you want to overwrite it?", interfaceName), false)
			if !confirmed {
				return nil
			}
		}
	}

	// 3. Prepare the data for template rendering
	data := WgConfTplData{
		InterfaceName:     interfaceName,
		PriKeyServer:      priKeyServer,
		ListenPort:        listenPort,
		Address:           IPAdressLocal,
		MTU:               mtu,
		PhysicalInterface: physicalInterface,
	}

	// 4. Parse and render the template
	tmpl, err := template.New("wgConfig").Parse(templates.WgConfTpl)
	if err != nil {
		return fmt.Errorf("failed to parse config template: %w", err)
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return fmt.Errorf("failed to render config template: %w", err)
	}

	// 6. Write the file with privilege 0600
	if err := os.WriteFile(configPath, buffer.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write config file to %s: %w", configPath, err)
	}

	// Add this into tracker
	if err := tracker.AddInterfaceToLog(interfaceName); err != nil {
		return fmt.Errorf("failed to track the interface: %w\n", err)
	}

	fmt.Printf("✅ Configuration file generated at: %s\n", configPath)
	return nil
}

/*
DeleteWGConfig removes the WireGuard configuration file for the given interface.
*/
func DeleteWGConfig(interfaceName string) error {
	configPath := filepath.Join(wgConfigDir, fmt.Sprintf("%s.conf", interfaceName))
	pubKeyPath := filepath.Join(wgConfigDir, fmt.Sprintf("%s.pub", interfaceName))
	PriKeyPath := filepath.Join(wgConfigDir, fmt.Sprintf("%s.key", interfaceName))

	// Stop the service and disable it
	if err := DisableServiceAutoStart(interfaceName, true); err != nil {
		return fmt.Errorf("failed to disable service: %w", err)
	}
	if err := StopService(interfaceName, true); err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}

	// Remove the configuration file and key files
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove configuration file %s: %w", configPath, err)
	}
	if err := os.Remove(pubKeyPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove public key file %s: %w", pubKeyPath, err)
	}
	if err := os.Remove(PriKeyPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove private key file %s: %w", PriKeyPath, err)
	}

	// Remove this from tracker
	if err := tracker.RemoveInterfaceFromLog(interfaceName); err != nil {
		return fmt.Errorf("failed to untrack the interface: %w\n", err)
	}

	fmt.Printf("✅ Configuration file for %s removed successfully.\n", interfaceName)
	return nil
}

/*
DeleteAllWGConfigs removes all WireGuard configuration files managed by fast-wireguard.
*/
func DeleteAllWGConfigs() error {
	// Delete all managed interfaces
	interfaces, err := tracker.GetAllManagedInterfaces()
	if err != nil {
		return err
	}

	for _, iface := range interfaces {
		if err := DeleteWGConfig(iface); err != nil {
			return err
		}
	}

	// Clear the tracker log file
	if err := os.Remove(tracker.InterfaceLogPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clear the tracker log file: %w", err)
	}

	return nil
}

/*
AddWGPeerConfig adds a peer configuration to the given WireGuard interface configuration file.
*/
func AddWGPeerConfig(
	interfaceName string,
	serverPublicIP string,
	listenPort int,
	mtu int,
	pubKeyServer string,
	peerName string,
	AllowedIPs string,
	pubKeyClient string,
	priKeyClient string,
) (string, error) {
	// 1. Make sure the path of the configuration file
	configPath := filepath.Join(wgConfigDir, fmt.Sprintf("%s.conf", interfaceName))

	// 2. Check whether the client configuration already exists
	peers, err := parseWGPeerConfig()
	if err != nil {
		return "", err
	}
	for _, peer := range peers {
		if peer.PubKeyClient == pubKeyClient {
			confirmed := utils.PromptConfirm(fmt.Sprintf("A peer with the same public key already exists in %s. Do you want to overwrite it?", configPath), false)
			if !confirmed {
				// Generate client configuration string
				peerName = peer.PeerName
				AllowedIPs = peer.AllowedIPs
			}
			// Overwrite the existing peer configuration whether confirmed or not, to generate client config string
			if err := DeleteWGPeerConfig(interfaceName, pubKeyClient, true); err != nil {
				return "", fmt.Errorf("failed to delete existing peer configuration: %w", err)
			}
			break
		}
	}

	// 3. Open the existing WireGuard configuration file and check if there is already a peer with the same pubKeyClient
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "", fmt.Errorf("configuration file for interface %s does not exist", interfaceName)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read configuration file: %w", err)
	}

	// 4. Handle the given peerName and AllowdIPs
	if strings.Contains(peerName, "[n]") {
		// Count existing peers
		existingPeers := bytes.Count(content, []byte("[Peer]"))
		peerName = strings.ReplaceAll(peerName, "[n]", fmt.Sprintf("%d", existingPeers+1))
	}

	if strings.Contains(AllowedIPs, "[n+1]") {
		// Count existing peers
		existingPeers := bytes.Count(content, []byte("[Peer]"))
		lastOctet := existingPeers + 2 // n starts from 1
		AllowedIPs = strings.ReplaceAll(AllowedIPs, "[n+1]", fmt.Sprintf("%d", lastOctet))
	}

	// 5. Prepare the data for template rendering
	data := PeerConfTplData{
		PeerName:     peerName,
		PubKeyClient: pubKeyClient,
		AllowedIPs:   AllowedIPs,
	}

	// 6. Parse and render the template
	tmpl, err := template.New("peerConfig").Parse(templates.PeerConfTpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse peer config template: %w", err)
	}
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return "", fmt.Errorf("failed to render peer config template: %w", err)
	}

	// 7. Append the peer configuration to the WireGuard config file
	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return "", fmt.Errorf("failed to open configuration file for appending: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(buffer.Bytes()); err != nil {
		return "", fmt.Errorf("failed to append peer configuration: %w", err)
	}
	fmt.Printf("✅ Peer configuration added to %s\n", configPath)

	// 8. Generate Client Configuration
	return GenerateWGClientConfig(
		serverPublicIP,
		listenPort,
		priKeyClient,
		pubKeyServer,
		AllowedIPs,
		mtu,
	)
}

/*
DeleteWGPeerConfig removes a peer configuration from the given WireGuard interface configuration file.
*/
func DeleteWGPeerConfig(
	interfaceName string,
	pubKeyClient string,
	silent bool,
) error {
	// 1. Make sure the path of the configuration file
	configPath := filepath.Join(wgConfigDir, fmt.Sprintf("%s.conf", interfaceName))

	// 2. Open the existing WireGuard configuration file
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file for interface %s does not exist", interfaceName)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// 3. Remove the peer configuration from the WireGuard config file
	sections := strings.Split(string(content), "[Peer]")
	var updatedSections []string
	updatedSections = append(updatedSections, sections[0]) // Keep the [Interface] section

	for _, section := range sections[1:] {
		if !strings.Contains(section, pubKeyClient) {
			updatedSections = append(updatedSections, "[Peer]"+section)
		}
	}

	updatedContent := strings.Join(updatedSections, "")

	// 4. Write the updated content back to the configuration file
	if err := os.WriteFile(configPath, []byte(updatedContent), 0600); err != nil {
		return fmt.Errorf("failed to write updated configuration file: %w", err)
	}

	if !silent {
		fmt.Printf("✅ Peer with public key %s removed from %s\n", pubKeyClient, configPath)
	}
	return nil
}

/*
parseWGPeerConfig reads the WireGuard configuration file and extracts peer configurations.
*/
func parseWGPeerConfig() ([]PeerConfTplData, error) {
	configPath := filepath.Join(wgConfigDir, fmt.Sprintf("%s.conf", "wg0"))

	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}
	// [Peer]
	// # Peer name {{ .PeerName }}
	// PublicKey = {{ .PubKeyClient }}
	// AllowedIPs = {{ .AllowedIP }}
	sections := strings.Split(string(content), "[Peer]")[1:]
	var peers []PeerConfTplData

	for _, section := range sections {
		lines := strings.Split(section, "\n")
		var peer PeerConfTplData
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if after, ok := strings.CutPrefix(line, "# Peer name"); ok {
				peer.PeerName = strings.TrimSpace(after)
			} else if after, ok := strings.CutPrefix(line, "PublicKey ="); ok {
				peer.PubKeyClient = strings.TrimSpace(after)
			} else if after, ok := strings.CutPrefix(line, "AllowedIPs ="); ok {
				peer.AllowedIPs = strings.TrimSpace(after)
			}
		}
		if peer.PubKeyClient != "" {
			peers = append(peers, peer)
		}
	}

	return peers, nil
}

/*
GenerateWGClientConfig generates the WireGuard client configuration string.
*/
func GenerateWGClientConfig(
	serverPublicIP string,
	listenPort int,
	priKeyClient string,
	pubKeyServer string,
	allowedIPs string,
	mtu int,
) (string, error) {
	host := serverPublicIP
	if ip := net.ParseIP(host); ip != nil && ip.To4() == nil {
		host = fmt.Sprintf("[%s]", host)
	}
	endpoint := fmt.Sprintf("%s:%d", host, listenPort)

	clientData := ClientConfTplData{
		PriKeyClient: priKeyClient,
		AllowedIPs:   allowedIPs,
		PubKeyServer: pubKeyServer,
		Endpoint:     endpoint,
		MTU:          mtu,
	}

	tmplClient, err := template.New("clientConfig").Parse(templates.ClientConfTpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse client config template: %w", err)
	}

	var clientBuffer bytes.Buffer
	if err := tmplClient.Execute(&clientBuffer, clientData); err != nil {
		return "", fmt.Errorf("failed to render client config template: %w", err)
	}

	return clientBuffer.String(), nil
}
