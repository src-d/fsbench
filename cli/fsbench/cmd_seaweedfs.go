package main

import "github.com/src-d/fsbench/fs"

type SeeweedFSCommand struct {
	Server string `short:"m" default:"127.0.0.1:9333" description:"SeaweedFS master."`
	Path   string `short:"p" default:"/tmp/" description:"Filesystem path."`
	BaseCommand
}

func (c *SeeweedFSCommand) Execute(args []string) error {
	fs := fs.NewSeaweedFSClient(c.Server, c.Path)
	c.BaseCommand.fs = fs

	return c.BaseCommand.Execute(args)
}
