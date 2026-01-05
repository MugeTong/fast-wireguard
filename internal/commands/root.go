package commands

import (
	"fast-wireguard/internal/commands/create"
	"fast-wireguard/internal/commands/uninstall"
	"fast-wireguard/internal/commands/update"
	"fast-wireguard/internal/commands/web"
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
	rootCmd.AddCommand(uninstall.CreateUninstallCmd())
	rootCmd.AddCommand(update.CreateUpdateCmd())
	rootCmd.AddCommand(web.CreateWebCmd())

	rootCmd.Flags().BoolP("version", "v", false, "the version of fast-wireguard")

	return rootCmd
}
