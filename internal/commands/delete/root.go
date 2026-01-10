package delete

import (
	"fast-wireguard/internal/wireguard"
	"fast-wireguard/pkg/utils"
	"fmt"
	"github.com/spf13/cobra"
)

func CreateDeleteCmd() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "delete [interface]",
		Short: "Delete the spicific interface",
		Args:  cobra.MaximumNArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			utils.EnsureRoot()
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Parse the arg
			interfaceName := "wg0"
			if len(args) > 0 {
				interfaceName = args[0]
			}
			if err :=wireguard.DeleteWGConfig(interfaceName); err != nil {
				fmt.Printf("Error in deleting the WireGuard server: %v\n", err)
			}

		},
	}
	return deleteCmd
}
