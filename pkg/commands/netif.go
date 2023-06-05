package commands

import (
	"flag"
	"fmt"

	"github.com/Code-Hex/vz/v3"
)

const netifCommandName = "netif"

type Netif struct {
	*flag.FlagSet
}

func NewNetifCommand() *Netif {
	return &Netif{
		flag.NewFlagSet(netifCommandName, flag.ContinueOnError),
	}
}

func (c *Netif) Run([]string) error {
	fmt.Println("The bridged network interfaces that you may use in your virtual machine:")

	networkInterfaces := vz.NetworkInterfaces()
	if len(networkInterfaces) == 0 {
		return fmt.Errorf("no network interfaces found")
	}

	for _, v := range networkInterfaces {
		fmt.Printf("- identifier: %s, name: %s\n", v.Identifier(), v.LocalizedDisplayName())
	}

	return nil
}

func (c *Netif) Name() string {
	return netifCommandName
}

func (c *Netif) flags() *flag.FlagSet {
	return c.FlagSet
}
