package commands

import (
	"fast-wireguard/internal/commands/create"
	"fast-wireguard/internal/commands/delete"
	"fast-wireguard/internal/commands/uninstall"
	"github.com/spf13/cobra"
)


func CreateRootCmd(version string) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:     "fwg",
		Version: version,
		Short:   "Quickly set up a WireGuard VPN server on your Linux machine with minimal configuration.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(create.CreateCreateCmd())
	rootCmd.AddCommand(delete.CreateDeleteCmd())
	rootCmd.AddCommand(uninstall.CreateUninstallCmd())


	rootCmd.Flags().BoolP("version", "v", false, "the version of fast-wireguard")

	return rootCmd
}
