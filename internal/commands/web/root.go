package web

import (
	"fast-wireguard/pkg/utils"
	"fmt"
	"github.com/spf13/cobra"
)

func CreateWebCmd() *cobra.Command {
	var webCmd = &cobra.Command{
		Use:   "web",
		Short: "Start the web server for parameter management",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			utils.EnsureRoot()
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting the web server...")
		},
	}
	return webCmd
}
