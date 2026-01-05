package update

import (
	"fast-wireguard/pkg/utils"
	"fmt"
	"github.com/spf13/cobra"
)

func CreateUpdateCmd() *cobra.Command {
	var update = &cobra.Command{
		Use:   "update",
		Short: "Update WireGuard parameters and restart (if in running)",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			utils.EnsureRoot()
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Updating parameters...")
		},
	}
	return update
}
