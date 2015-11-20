package main

import "github.com/src-d/fsbench/fs"

type OSCommand struct {
	Path string `short:"p" default:"/tmp/" description:"Path where the test files will be written."`
	BaseCommand
}

func (c *OSCommand) Execute(args []string) error {
	fs := fs.NewOSClient(c.Path)
	c.BaseCommand.fs = fs

	return c.BaseCommand.Execute(args)
}
