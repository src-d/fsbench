package fsbench

import (
	"os"
	"testing"

	"github.com/src-d/fsbench/fs"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type WorkerSuite struct{}

var _ = Suite(&WorkerSuite{})

func (s *WorkerSuite) TestGetFilename(c *C) {
	w := &Worker{c: &WorkerConfig{
		DirectoryDepth: 5,
	}}

	fn := w.getFilename()
	c.Assert(fn[2], Equals, byte(os.PathSeparator))
	c.Assert(fn[5], Equals, byte(os.PathSeparator))
	c.Assert(fn[8], Equals, byte(os.PathSeparator))
	c.Assert(fn[11], Equals, byte(os.PathSeparator))
	c.Assert(fn[14], Equals, byte(os.PathSeparator))
}

func (s *WorkerSuite) TestCreate(c *C) {
	cli := fs.NewMemoryClient()
	w := NewWorker(cli, &WorkerConfig{
		Files:       10,
		BlockSize:   512,
		MinFileSize: 1024 * 100,
		MaxFileSize: 1024 * 100,
	})

	c.Assert(w.Do(), IsNil)
	c.Assert(cli.Files, HasLen, 10)
	for fn, _ := range cli.Files {
		s, _ := cli.Stat(fn)
		c.Assert(s.Size(), Equals, int64(1024*100))
	}

	c.Assert(w.Status.Files, Equals, 10)
	c.Assert(w.Status.Bytes, Equals, int64(1024*100*10))
	c.Assert(w.Status.Errors, Equals, 0)
}

func (s *WorkerSuite) TestCreateRand(c *C) {
	cli := fs.NewMemoryClient()
	w := NewWorker(cli, &WorkerConfig{
		Files:       10,
		BlockSize:   512,
		MinFileSize: 1024 * 100,
		MaxFileSize: 1024 * 200,
	})

	c.Assert(w.Do(), IsNil)
	c.Assert(cli.Files, HasLen, 10)
	for fn, _ := range cli.Files {
		s, _ := cli.Stat(fn)
		c.Assert(s.Size() >= 1024*100, Equals, true)
		c.Assert(s.Size() <= 1024*200, Equals, true)
	}

	c.Assert(w.Status.Files, Equals, 10)
	c.Assert(w.Status.Errors, Equals, 0)
}
