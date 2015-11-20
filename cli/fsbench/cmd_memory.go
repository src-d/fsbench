package main

import "github.com/src-d/fsbench/fs"

const (
	MemoryCommandDescription     = "Memory based benchmark."
	MemoryCommandLongDescription = MemoryCommandDescription + "\n" + BaseLongDescription
)

type MemoryCommand struct {
	BaseCommand
}

func (c *MemoryCommand) Execute(args []string) error {
	c.BaseCommand.fs = fs.NewMemoryClient()

	return c.BaseCommand.Execute(args)
}
