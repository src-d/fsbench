package main

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/rakyll/pb"
	"github.com/src-d/fsbench"
	"github.com/src-d/fsbench/fs"
)

const (
	MaxBlockSize        int64  = 1024 * 1024 * 1024
	BaseLongDescription string = `
- Workers (-w): Number of workers to run concurrently.

- Files (-f): Number of files to write.

- SizeBlock (-b): The writes are made on blocks of the given size, by default
  the value of is the value fixed file size (-s), max. size of the block is 1GB.

- FixedFileSize (-s): Size of the files to be written. If this value is set all
  the files written by the test are of the given size.

- DirectoryDepth (-d): if the directory depth number is different of 0 the files
  are written on directories, the directories are created using the first two
  chars from the file. Example: if deep is 2 a file name "abefghif.rand" i
  transformed on: "ab/ef/ghif.rand".
`
)

type BaseCommand struct {
	Workers        int   `short:"w" default:"4" description:"Number of workers to run concurrently."`
	Files          int   `short:"f" default:"100" description:"Number of files to write."`
	BlockSize      int64 `short:"b" default:"0" description:"Size of the block"`
	FixedFileSize  int64 `short:"s" default:"1048576" description:"Size of the files to be written."`
	DirectoryDepth int   `short:"d" default:":0" description:"Directory depth"`

	b  *fsbench.Benchmark
	pb *pb.ProgressBar
	fs fs.Client
}

func (c *BaseCommand) Execute(args []string) error {
	c.init()
	go c.updateProgressBar()

	status := c.b.Run()
	c.pb.Set(status.WStatus.Files)
	c.pb.Finish()
	c.printStatus(status.WStatus)
	c.printStatus(status.RStatus)

	return nil
}

func (c *BaseCommand) init() {
	if c.BlockSize == 0 {
		c.BlockSize = c.FixedFileSize
	}

	if c.BlockSize > MaxBlockSize {
		c.BlockSize = MaxBlockSize
	}

	c.b = fsbench.NewBenchmark(c.fs, &fsbench.Config{
		Workers:       c.Workers,
		Files:         c.Files,
		BlockSize:     c.BlockSize,
		FixedFileSize: c.FixedFileSize,
	})

	c.b.Init()

	c.pb = pb.StartNew(c.Files)
	c.pb.ShowTimeLeft = true
	c.pb.Format(" ▓▒░ ")
}

func (c *BaseCommand) updateProgressBar() {
	for {
		status := c.b.Status()
		c.pb.Set(status.WStatus.Files)
		time.Sleep(time.Millisecond)
		if status.WStatus.Files >= c.Files {
			break
		}
	}
}

func (c *BaseCommand) printStatus(s *fsbench.AggregatedStatus) {
	secs := s.Duration.Seconds() / float64(c.Workers)

	fmt.Printf(
		"Summary:\n  - Files: %d\n  - Errors: %d\n  - Size: %s\n  - Speed: %s/s\n",
		s.Files, s.Errors,
		humanize.Bytes(uint64(s.Bytes)),
		humanize.Bytes(uint64(float64(s.Bytes)/secs)),
	)
}
