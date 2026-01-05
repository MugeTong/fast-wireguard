package wireguard

import (
	"errors"
	"fast-wireguard/pkg/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/*
IsWireGuardInstalled checks if WireGuard is installed on the system.

Returns true if installed, false otherwise.
*/
func IsWireGuardInstalled() bool {
	_, err := exec.LookPath("wg")
	return err == nil
}

/*
GetWireGuardVersion retrieves the installed WireGuard version.

Returns the version string and any error encountered.
*/
func GetWireGuardVersion() (string, error) {
	out, err := exec.Command("wg", "--version").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

/*
InstallWireGuard checks and installs WireGuard
*/
func InstallWireGuard() error {
	// Check if WireGuard is already installed
	if IsWireGuardInstalled() {
		// fmt.Println("WireGuard is already installed.")
		return nil
	}
	// Prompt user for installation
	if !utils.PromptConfirm("WireGuard is not installed. Are you sure to install it now?", true) {
		fmt.Println("Installation cancelled by user.")
		os.Exit(0)
	}

	// Install the WireGuard package using the system's package manager
	managers := map[string][]string{
		"apt":    {"apt", "install", "-y", "wireguard-tools"},
		"yum":    {"yum", "install", "-y", "wireguard-tools"},
		"dnf":    {"dnf", "install", "-y", "wireguard-tools"},
		"pacman": {"pacman", "-S", "--noconfirm", "wireguard-tools"},
	}

	for name, args := range managers {
		if _, err := exec.LookPath(name); err == nil {
			fmt.Printf("Installing WireGuard via %s...\n", name)
			if err := utils.RunAsRoot(args[0], args[1:]...); err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("No supported package manager found. Please install WireGuard manually.")
}
