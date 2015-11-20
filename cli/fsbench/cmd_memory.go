package main

import "github.com/src-d/fsbench/fs"

type MemoryCommand struct {
	BaseCommand
}

func (c *MemoryCommand) Execute(args []string) error {
	c.BaseCommand.fs = fs.NewMemoryClient()

	return c.BaseCommand.Execute(args)
}
