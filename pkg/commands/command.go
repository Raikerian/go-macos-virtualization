package commands

import "flag"

type Flags interface {
	flags() *flag.FlagSet
}

type Command interface {
	Flags

	Run([]string) error
	Name() string
}
