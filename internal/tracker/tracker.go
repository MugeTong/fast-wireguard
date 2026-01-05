package tracker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	interfaceLogPath = "/etc/wireguard/.fwg_managed_interfaces"
)

/*
WireInterface tracks the given interface into file (if not already recorded).
*/
func AddInterfaceToLog(interfaceName string) error {
	// Check if already recorded
	exists, err := IsManagedByUs(interfaceName)
	if err != nil {
		return err
	}
	if exists {
		// Already recorded, just return
		return nil
	}

	// Open the file in [append | create | write-only] mode
	// os.O_APPEND: append content to the end
	// os.O_CREATE: create if not exists
	// os.O_WRONLY: write-only mode
	// 0644: permissions, root can read/write, others can read
	file, err := os.OpenFile(interfaceLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Cannot open tracker file: %w", err)
	}
	defer file.Close()

	// 3. Write the interface name + newline
	if _, err := file.WriteString(interfaceName + "\n"); err != nil {
		return fmt.Errorf("Cannot write to tracker file: %w", err)
	}

	return nil
}

/*
IsManagedByUs checks if the given interface is managed by us (i.e., recorded in the tracker file).

Returns true if managed, false otherwise, along with any error encountered.
*/
func IsManagedByUs(interfaceName string) (bool, error) {
	// Read the file
	file, err := os.Open(interfaceLogPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, so it is definitely not created by us
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	// Scan the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == interfaceName {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

/*
RemoveInterfaceFromLog removes the interface (used for uninstallation or service remove)
*/
func RemoveInterfaceFromLog(interfaceName string) error {
	// Read the file
	content, err := os.ReadFile(interfaceLogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string

	// Filter the interface
	found := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue // Skip the empty line
		}
		if trimmed == interfaceName {
			found = true
			continue // Skip the ovjective line
		}
		newLines = append(newLines, trimmed)
	}

	if !found {
		return nil // The interface is not tracked.
	}

	// Overwrite the log file
	output := strings.Join(newLines, "\n")
	// Ensure that at least one empty line exists
	if len(output) > 0 {
		output += "\n"
	}
	if err := os.WriteFile(interfaceLogPath, []byte(output), 0644); err != nil {
		return fmt.Errorf("Failed to the tracker file: %w", err)
	}

	return nil
}
