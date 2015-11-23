package main

import "github.com/src-d/fsbench/fs"

const (
	GFAPICommandDescription     = "GlusterFS (via gfapi) based benchmark"
	GFAPICommandLongDescription = GFAPICommandDescription + `
	
- Path (-p): Path where the test files will be written.
` + BaseLongDescription
)

type GFAPICommand struct {
	Server    string `short:"g" description:"GlusterFS server."`
	Datastore string `short:"d" description:"Datastore name."`
	Path      string `short:"p" default:"/" description:"Filesystem path."`
	BaseCommand
}

func (c *GFAPICommand) Execute(args []string) error {
	fs := fs.NewGFAPIClient(c.Path, !c.Cache)
	c.BaseCommand.fs = fs

	return c.BaseCommand.Execute(args)
}
