package main

import "github.com/src-d/fsbench/fs"

const (
	OSCommandDescription     = "Filesystem based benchmark"
	OSCommandLongDescription = OSCommandDescription + `
	
- Path (-p): Path where the test files will be written.
` + BaseLongDescription
)

type OSCommand struct {
	Path string `short:"p" default:"/tmp/" description:"Filesystem path."`
	BaseCommand
}

func (c *OSCommand) Execute(args []string) error {
	fs := fs.NewOSClient(c.Path)
	c.BaseCommand.fs = fs

	return c.BaseCommand.Execute(args)
}
