package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dripolles/histogram"
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
	Workers        int    `short:"w" default:"4" description:"Number of workers to run concurrently."`
	Files          int    `short:"f" default:"100" description:"Number of files to write."`
	BlockSize      int64  `short:"b" default:"0" description:"Size of the block"`
	FixedFileSize  int64  `short:"s" default:"1048576" description:"Size of the files to be written."`
	DirectoryDepth int    `short:"d" default:"0" description:"Directory depth"`
	Output         string `short:"o" default:"fsbench_%s" description:"Output filename"`

	b  *fsbench.Benchmark
	pb *pb.ProgressBar
	fs fs.Client
}

func (c *BaseCommand) Execute(args []string) error {
	c.init()
	go c.updateProgressBar()

	status := c.b.Run()
	c.pb.Set(status.WStatus.Files + status.RStatus.Files)
	c.pb.Finish()
	c.printStatus(status)

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

	c.pb = pb.StartNew(c.Files * 2)
	c.pb.ShowTimeLeft = true
	c.pb.Format(" ▓▒░ ")
}

func (c *BaseCommand) updateProgressBar() {
	total := c.Files * 2
	for {
		status := c.b.Status()
		count := status.WStatus.Files + status.RStatus.Files
		c.pb.Set(count)
		time.Sleep(time.Millisecond)
		if count >= total {
			break
		}
	}
}

func (c *BaseCommand) printStatus(s *fsbench.BenchmarkStatus) {
	fmt.Println("\nWrite Stats\n==========")
	c.printAggregatedStatus(s.WStatus)
	fmt.Println("\nRead Stats\n==========")
	c.printAggregatedStatus(s.RStatus)

	if c.Output == "" {
		return
	}

	base := fmt.Sprintf(c.Output, time.Now().Format("2006-01-02T150405"))
	c.saveAggregatedSStatus(s.WStatus, base+".w.json")
	c.saveAggregatedSStatus(s.RStatus, base+".r.json")
}

var percentiles = []float64{.99, .95, .90, .75, .50, .25, .10}

func (c *BaseCommand) printAggregatedStatus(s *fsbench.AggregatedStatus) {
	secs := s.Duration.Seconds() / float64(c.Workers)

	fmt.Printf(
		"Summary:\n  - Files: %d\n  - Errors: %d\n  - Size: %s\n  - Speed: %s/s\n\n",
		s.Files, s.Errors,
		humanize.Bytes(uint64(s.Bytes)),
		humanize.Bytes(uint64(float64(s.Bytes)/secs)),
	)

	fmt.Println("Avg. Speed Percentiles")
	for _, p := range percentiles {
		value := s.HistogramAvgRate.GetAtPercentile(p)
		fmt.Printf("  - p%.0f) %s/s\n", p*100, humanize.Bytes(uint64(value)))
	}
}

func (c *BaseCommand) saveAggregatedSStatus(s *fsbench.AggregatedStatus, filename string) {
	data := map[string]interface{}{}
	data["Config"] = c
	data["General"] = s
	data["AvgRatePercentiles"] = c.getPercentiles(s.HistogramAvgRate)

	js, _ := json.MarshalIndent(data, "", "  ")

	err := ioutil.WriteFile(filename, js, 0644)
	if err != nil {
		fmt.Printf("Error saving %q: %s", err, filename)
	}
}

func (c *BaseCommand) getPercentiles(h *histogram.Histogram) map[string]int {
	table := map[string]int{}
	for _, p := range percentiles {
		table[fmt.Sprintf("p%.0f", p*100)] = h.GetAtPercentile(p)
	}

	return table
}
