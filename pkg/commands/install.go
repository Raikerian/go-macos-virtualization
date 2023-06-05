package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/raikerian/go-macos-virtualization/pkg/macos"
)

const installCommandName = "install"

type Install struct {
	*flag.FlagSet
}

func NewInstallCommand() *Install {
	return &Install{
		flag.NewFlagSet(installCommandName, flag.ContinueOnError),
	}
}

func (c *Install) Name() string {
	return installCommandName
}

func (c *Install) Run([]string) error {
	fmt.Println("Performing macOS installation...")
	ctx := context.Background()
	return macos.Install(ctx)
}

func (c *Install) flags() *flag.FlagSet {
	return c.FlagSet
}
