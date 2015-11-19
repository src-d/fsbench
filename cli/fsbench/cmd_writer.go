package main

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/rakyll/pb"
	"github.com/src-d/fsbench"
)

type WriterCommand struct {
	Filesystem     string
	Workers        int   `short:"w" default:"4" description:"Number of workers to run concurrently."`
	Files          int   `short:"f" default:"100" description:"Number of files to write."`
	BlockSize      int64 `short:"b" default:"4096" description:"Size of the block, the writes are done on blocks of the given size"`
	MinFileSize    int64 `short:"s" default:"1048576" description:"Minimun size of the files to be written."`
	MaxFileSize    int64 `short:"S" default:"1048576" description:"Maximum size of the files to be written. If the min size and max size are different, the size of the files are random numbers on the range."`
	DirectoryDepth int   `short:"d" default:":0" description:"Number of directories to be created for each file. Avoid having large amounts of files on the same dir."`

	b  *fsbench.Benchmark
	pb *pb.ProgressBar
}

func (c *WriterCommand) Execute(args []string) error {
	c.init()
	go c.updateProgressBar()

	status := c.b.Run()
	c.pb.Set(status.Files)
	c.pb.Finish()
	c.printStatus(status)

	return nil
}

func (c *WriterCommand) init() {
	c.b = fsbench.NewBenchmark(&fsbench.Config{
		Workers:     c.Workers,
		Files:       c.Files,
		BlockSize:   c.BlockSize,
		MinFileSize: c.MinFileSize,
		MaxFileSize: c.MaxFileSize,
	})

	c.b.Init()

	c.pb = pb.StartNew(c.Files)
	c.pb.ShowTimeLeft = true
	c.pb.Format(" ▓▒░ ")
}

func (c *WriterCommand) updateProgressBar() {
	for {
		status := c.b.Status()
		c.pb.Set(status.Files)
		time.Sleep(time.Millisecond)

		if status.Files >= c.Files {
			break
		}
	}
}

func (c *WriterCommand) printStatus(s *fsbench.Status) {
	secs := s.Duration.Seconds() / float64(c.Workers)

	fmt.Printf(
		"Summary:\n  - Files: %d\n  - Errors: %d\n  - Size: %s\n  - Speed: %s/s\n",
		s.Files, s.Errors,
		humanize.Bytes(uint64(s.Bytes)),
		humanize.Bytes(uint64(float64(s.Bytes)/secs)),
	)
}
