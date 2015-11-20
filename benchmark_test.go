package fsbench

import (
	"runtime"

	. "gopkg.in/check.v1"
)

type BenchmarkSuite struct{}

var _ = Suite(&BenchmarkSuite{})

func (s *BenchmarkSuite) TestInit(c *C) {
	b := NewBenchmark(&Config{Workers: 3, Files: 10})
	b.Init()

	c.Assert(b.w, HasLen, 3)
	c.Assert(b.w[0].c.Files, Equals, 3)
	c.Assert(b.w[1].c.Files, Equals, 3)
	c.Assert(b.w[2].c.Files, Equals, 4)
}

func (s *BenchmarkSuite) TestRun(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	b := NewBenchmark(&Config{
		Workers:       100,
		Files:         10000,
		BlockSize:     512,
		FixedFileSize: 1024 * 100,
	})

	b.Init()
	status := b.Run()

	c.Assert(status.Files, Equals, 10000)
	c.Assert(status.Errors, Equals, 0)
	c.Assert(status.Bytes, Equals, int64(1024000000))
}
