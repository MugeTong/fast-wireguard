package uninstall

import (
	"fast-wireguard/internal/system"
	"fast-wireguard/pkg/utils"
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

func CreateUninstallCmd() *cobra.Command {
	var uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall this application and restore modifications",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Prompt the user for confirmation before uninstalling
			if os.Getenv("USER_CONFIRM") != "1" {
				if !utils.PromptConfirm("Are you sure you want to uninstall Fast-Wireguard?", false) {
					fmt.Println("Uninstallation cancelled.")
					os.Exit(0)
				}
			}
			// If the user confirmed, ensure we have root privileges
			utils.EnsureRoot()
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Uninstalling Fast-Wireguard...")
			// Remove sysctl configuration file if it exists
			if err := system.RestoreIPForwarding(); err != nil {
				fmt.Println("Error restoring IP forwarding settings:", err)
				// Continue uninstalling even if restoring IP forwarding fails
			}

			// Remove the binary file
			binaryPath := "/usr/local/bin/fwg"
			if err := os.Remove(binaryPath); err != nil {
				if !os.IsNotExist(err) {
					fmt.Printf("Error removing binary file: %v\n", err)
					fmt.Println("Please remove it manually: sudo rm " + binaryPath)
					os.Exit(1)
				}
			} else {
				fmt.Println("Binary file removed successfully.")
			}

			fmt.Println("Fast-Wireguard uninstallation completed.")
		},
	}
	return uninstallCmd
}
