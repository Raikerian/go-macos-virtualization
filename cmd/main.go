package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/raikerian/go-macos-virtualization/pkg/macos"
)

var install bool
var cpuCount uint
var memorySize uint64

func init() {
	flag.BoolVar(&install, "install", false, "run command as install mode")
	flag.UintVar(&cpuCount, "cpu", 4, "number of cpu cores")
	flag.Uint64Var(&memorySize, "memory", 16, "memory size in GB")
}

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	if install {
		return macos.Install(ctx)
	}

	m, err := macos.NewManager(cpuCount, memorySize)
	if err != nil {
		return err
	}

	return m.Run(ctx)
}
