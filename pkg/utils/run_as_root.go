package utils

import (
	"fmt"
	"os"
	"errors"
	"os/exec"
	"strings"
)

/*
IsRoot checks if the current user is root.
Returns true if running as root, false otherwise.
*/
func IsRoot() bool {
	return os.Geteuid() == 0
}

/*
RunAsRoot runs the given command with root privileges.
If the current user is not root, it attempts to re-execute the program with sudo.
*/
func RunAsRoot(command string, arg ...string) error {
	// Validate input
	var commandSlice = strings.Fields(command)
	if len(commandSlice) == 0 {
		return errors.New("No command provided to run as root.")
	}
	commandSlice = append(commandSlice, arg...)

	if os.Geteuid() == 0 {
		// Already running as root, execute the command directly
		cmd := exec.Command(commandSlice[0], commandSlice[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Command execution failed:", err)
			return err
		}
		return nil
	}
	fmt.Println("Detected non-root user, attempting to auto elevate privileges...")
	// Not running as root, re-execute with sudo
	args := []string{"FWG_ELEVATED=1"}
	args = append(args, commandSlice...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to elevate privileges:", err)
		return err
	}
	return nil
}

func RunAsRootSilent(command string, arg ...string) error {
	// Validate input
	var commandSlice = strings.Fields(command)
	if len(commandSlice) == 0 {
		return errors.New("No command provided to run as root.")
	}
	commandSlice = append(commandSlice, arg...)

	if os.Geteuid() == 0 {
		// Already running as root, execute the command directly
		cmd := exec.Command(commandSlice[0], commandSlice[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}
	// Not running as root, re-execute with sudo
	args := []string{"FWG_ELEVATED=1"}
	args = append(args, commandSlice...)
	cmd := exec.Command("sudo", args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
