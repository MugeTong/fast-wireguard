package create

import (
	"fast-wireguard/internal/system"
	"fast-wireguard/internal/wireguard"
	"fast-wireguard/pkg/utils"
	"fmt"
	"github.com/spf13/cobra"
)

/*
CreateCreateCmd represents the create command to install the application and dependencies.
It sets up the initial configuration and runs the service.
*/
func CreateCreateCmd() *cobra.Command {
	opts := &wireguard.ServerOptions{}
	var createCmd = &cobra.Command{
		Use:     "create [interface]",
		Aliases: []string{"install", "setup"},
		Short:   "Create initial configuration and run the service.",
		Long: `Set up a new WireGuard interface.
If no interface name is provided, it defaults to 'wg0'.
You can specify a custom name like 'wg1' or 'myvpn' to create multiple instances.`,
		Args: cobra.MaximumNArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			utils.EnsureRoot()
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Check the package installation and IP forwarding
			if err := wireguard.InstallWireGuard(); err != nil {
				fmt.Printf("Error installing WireGuard: %v\n", err)
				return
			}
			if err := system.EnableIPForwarding(); err != nil {
				fmt.Printf("Error enabling IP forwarding: %v\n", err)
				return
			}

			// If the user choose to just config the WireGuard environment, exit
			if opts.DryRun {
				return
			}

			// Setup the server configuration
			interfaceName := "wg0"
			if len(args) > 0 {
				interfaceName = args[0]
			}
			if err := wireguard.CreateServer(interfaceName, opts); err != nil {
				fmt.Printf("Error in creating the WireGuard server: %v\n", err)
				return
			}

			// Start the service and enable it to start automatically on boot
			if err := wireguard.EnableServiceAutoStart(interfaceName); err != nil {
				fmt.Printf("Error in enabling the service %s: %v\n", interfaceName, err)
				return
			}
			if err := wireguard.StartService(interfaceName); err != nil {
				fmt.Printf("Error in starting the service %s: %v\n", interfaceName, err)
			}
		},
	}

	createCmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "initial the configuration for WireGuard without creating interface")
	createCmd.Flags().IntVarP(&opts.ListenPort, "port", "p", 51820, "listening port for WireGuard server")
	createCmd.Flags().StringVarP(&opts.IPAdressLocal, "address", "a", "10.0.0.1/24, fd00::1/64", "local IP address assigned to WireGuard server")
	createCmd.Flags().IntVarP(&opts.MTU, "mtu", "m", 1420, "the length of MTU")
	createCmd.Flags().StringVarP(&opts.PeerName, "peer-name", "n", "default-peer[n]", "name of the WireGuard client peer")
	createCmd.Flags().StringVarP(&opts.IPAdressLocalClient, "address-client", "c", "10.0.0.[n+1]/32, fd00::[n+1]/128", "local IP address assigned to WireGuard client")
	createCmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "force re-setup even if already configured")

	return createCmd
}
