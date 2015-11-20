package fsbench

import (
	"runtime"

	"github.com/src-d/fsbench/fs"

	. "gopkg.in/check.v1"
)

type BenchmarkSuite struct{}

var _ = Suite(&BenchmarkSuite{})

func (s *BenchmarkSuite) TestInit(c *C) {
	fs := fs.NewMemoryClient()
	b := NewBenchmark(fs, &Config{Workers: 3, Files: 10})
	b.Init()

	c.Assert(b.w, HasLen, 3)
	c.Assert(b.w[0].c.Files, Equals, 3)
	c.Assert(b.w[1].c.Files, Equals, 3)
	c.Assert(b.w[2].c.Files, Equals, 4)
}

func (s *BenchmarkSuite) TestRun(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fs := fs.NewMemoryClient()
	b := NewBenchmark(fs, &Config{
		Workers:       100,
		Files:         10000,
		BlockSize:     512,
		FixedFileSize: 1024 * 100,
	})

	b.Init()
	status := b.Run()

	c.Assert(status.WStatus.Files, Equals, 10000)
	c.Assert(status.WStatus.Errors, Equals, 0)
	c.Assert(status.WStatus.Bytes, Equals, int64(1024000000))
}
