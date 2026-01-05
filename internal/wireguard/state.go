package wireguard

import (
	"fast-wireguard/internal/system"
	"fast-wireguard/pkg/utils"
	"fmt"
)

type ServerOptions struct {
	DryRun              bool
	ListenPort          int
	IPAdressLocal       string
	IPAdressLocalClient string
	PeerName            string
	MTU                 int
	Force               bool
}

/*
StartService starts the WireGuard service for the given interface.
*/
func StartService(interfaceName string) error {
	serviceName := fmt.Sprintf("wg-quick@%s", interfaceName)
	if err := utils.RunAsRoot("systemctl", "start", serviceName); err != nil {
		return err
	}
	fmt.Printf("✅ Service %s started.\n", interfaceName)
	return nil
}

/*
StopService stops the WireGuard service for the given interface.
*/
func StopService(interfaceName string, silent bool) error {
	serviceName := fmt.Sprintf("wg-quick@%s", interfaceName)
	if err := utils.RunAsRoot("systemctl", "stop", serviceName); err != nil {
		return err
	}
	if !silent {
		fmt.Printf("✅ Service %s stopped.\n", interfaceName)
	}
	return nil
}

/*
RestartService restart the WireGuard service for the given interface.
*/
func RestartService(interfaceName string) error {
	serviceName := fmt.Sprintf("wg-quick@%s", interfaceName)
	if err := utils.RunAsRoot("systemctl", "restart", serviceName); err != nil {
		return err
	}
	fmt.Printf("✅ Service %s restarted.\n", interfaceName)
	return nil
}

/*
EnableServiceAutoStart allows the service for the given interface to start automatically on boot.
*/
func EnableServiceAutoStart(interfaceName string) error {
	serviceName := fmt.Sprintf("wg-quick@%s", interfaceName)
	if err := utils.RunAsRootSilent("systemctl", "enable", serviceName); err != nil {
		return err
	}
	fmt.Printf("✅ Service %s enabled to start automatically on boot.\n", interfaceName)
	return nil
}

/*
EnableServiceAutoStart disables the service for the given interface to start automatically on boot.
*/
func DisableServiceAutoStart(interfaceName string, silent bool) error {
	serviceName := fmt.Sprintf("wg-quick@%s", interfaceName)
	if err := utils.RunAsRootSilent("systemctl", "disable", serviceName); err != nil {
		return err
	}
	if !silent {
		fmt.Printf("✅ Service %s disabled from starting automatically on boot.\n", interfaceName)
	}
	return nil
}

/*
SetupServer setup the WireGuard server with the following steps:

  - generate Wireguard key pair
  - get the physical interface
  - get the ip adress
  - generate the WireGuard server configuration file
*/
func CreateServer(interfaceName string, opts *ServerOptions) error {
	// 1. Generate WireGuard key pair
	priKey, pubKey, err := GenerateWGKeys(interfaceName)
	if err != nil {
		return err
	}
	// 2. Get the physical interface and IP address
	PhysicalInterface, err := system.GetPhysicalInterface()
	if err != nil {
		return err
	}
	IPAdress, err := system.GetPublicIP()
	if err != nil {
		return err
	}

	// 3. Generate the WireGuard configuration file
	if err := GenerateWGConfig(
		interfaceName,
		opts.ListenPort,
		priKey,
		opts.MTU,
		opts.IPAdressLocal,
		PhysicalInterface,
		opts.Force); err != nil {
		return err
	}

	// 4. Collect peer information and add peer configuration
	pubKeyClient := utils.PromptInput("Input the public key for the peer (Enter to skip):", "", false)
	if pubKeyClient == "" {
		fmt.Println("Skipping peer configuration addition.")
		fmt.Printf("✅ WireGuard server %s configured successfully.\n", interfaceName)
		return nil
	}
	priKeyClient := utils.PromptInput("Input the private key for the peer (Not necessary, Enter to skip):", "", false)
	if priKeyClient == "" {
		priKeyClient = "<your_client_private_key>"
	}
	clientConfString, err := AddWGPeerConfig(
		interfaceName,
		opts.PeerName,
		opts.IPAdressLocalClient,
		pubKeyClient,
		priKeyClient,
		IPAdress,
		pubKey,
		opts.MTU,
		opts.ListenPort)
	if err != nil {
		return err
	}

	fmt.Printf("✅ WireGuard server %s configured successfully.\n", interfaceName)

	if clientConfString != "" {
		fmt.Println("\nClient configuration:")
		fmt.Println("------------------------------------------------")
		fmt.Println(clientConfString)
		fmt.Println("------------------------------------------------")
		// Print the client configuration QR code
		// utils.PrintQRCode("You can also scan this QR code with your WireGuard App:\n", clientConfString)
	}
	return nil
}
