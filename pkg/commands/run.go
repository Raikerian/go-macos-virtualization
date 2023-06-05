package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/raikerian/go-macos-virtualization/pkg/macos"
	"github.com/raikerian/go-macos-virtualization/pkg/utils"
)

const runCommandName = "run"

type Run struct {
	*flag.FlagSet

	cpuCount   uint
	memorySize uint64
}

func NewRunCommand() *Run {
	c := &Run{
		FlagSet: flag.NewFlagSet(runCommandName, flag.ExitOnError),
	}

	c.UintVar(&c.cpuCount, "cpu", utils.ComputeCPUCount(), "number of cpu cores")
	c.Uint64Var(&c.memorySize, "memory", utils.ComputeMemorySize(), "memory size, must be a multiple of a 1 megabyte (1024 * 1024 bytes)")
	return c
}

func (c *Run) Name() string {
	return runCommandName
}

func (c *Run) Run(args []string) error {
	c.Parse(args)

	fmt.Println("Running VM...")

	ctx := context.Background()
	m, err := macos.NewManager(c.cpuCount, c.memorySize)
	if err != nil {
		return err
	}
	return m.Run(ctx)
}

func (c *Run) flags() *flag.FlagSet {
	return c.FlagSet
}
