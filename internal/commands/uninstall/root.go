package uninstall

import (
	"fast-wireguard/internal/system"
	"fast-wireguard/internal/wireguard"
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
			// 1. Remove sysctl configuration file if it exists
			if err := system.RestoreIPForwarding(); err != nil {
				fmt.Println("Error restoring IP forwarding settings:", err)
				// Continue uninstalling even if restoring IP forwarding fails
			}

			// 2. Delete WireGuard service files
			confirmed := utils.PromptConfirm("Do you want to delete the WireGuard configuration files?", false)
			if confirmed {
				if err := wireguard.DeleteAllWGConfigs(); err != nil {
					fmt.Println("Error deleting WireGuard configuration files:", err)
					// Continue uninstalling even if deleting configs fails
				} else {
					fmt.Println("✅ WireGuard configuration files deleted successfully.")
				}
			} else {
				fmt.Println("Skipping deletion of WireGuard configuration files.")
			}

			// 4. Remove the binary file
			binaryPath := "/usr/local/bin/fwg"
			if err := os.Remove(binaryPath); err != nil {
				if !os.IsNotExist(err) {
					fmt.Printf("Error removing binary file: %v\n", err)
					fmt.Println("Please remove it manually: sudo rm " + binaryPath)
					os.Exit(1)
				}
			} else {
				fmt.Println("✅ Binary file removed successfully.")
			}

			fmt.Println("✅ Fast-Wireguard uninstallation completed.")
		},
	}
	return uninstallCmd
}
