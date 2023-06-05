package main

import (
	"fmt"
	"os"

	"github.com/raikerian/go-macos-virtualization/pkg/commands"
)

var (
	cmds = []commands.Command{
		commands.NewNetifCommand(),
		commands.NewInstallCommand(),
		commands.NewRunCommand(),
	}
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("you must pass a sub-command")
	}

	subcommand := os.Args[1]
	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			return cmd.Run(os.Args[2:])
		}
	}

	return fmt.Errorf("unknown subcommand: %s, use `%s help` to get list of all available commands", subcommand, os.Args[0])
}
