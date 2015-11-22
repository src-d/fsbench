package main

import "github.com/src-d/fsbench/fs"

const (
	OSCommandDescription     = "Filesystem based benchmark"
	OSCommandLongDescription = OSCommandDescription + `
	
- Path (-p): Path where the test files will be written.

- Cache (--cache): Enable the usage of cache, by default fsbench make the writes
  using O_DIRECT flag skipping write and read from cache.
` + BaseLongDescription
)

type OSCommand struct {
	Path  string `short:"p" default:"/tmp/" description:"Filesystem path."`
	Cache bool   `long:"cache" description:"Allow cache writers."`
	BaseCommand
}

func (c *OSCommand) Execute(args []string) error {
	fs := fs.NewOSClient(c.Path, !c.Cache)
	c.BaseCommand.fs = fs

	return c.BaseCommand.Execute(args)
}
