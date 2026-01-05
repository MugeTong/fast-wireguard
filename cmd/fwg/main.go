package main

import (
	"fast-wireguard/internal/commands"
)

var version = "dev"

func main() {
	commands.CreateRootCmd(version).Execute()
}
