package fs

import (
	"bytes"

	. "gopkg.in/check.v1"
)

func (s *WritersSuite) TestMemoryClient_Create(c *C) {
	client := NewMemoryClient()
	f, _ := client.Create("foo")

	c.Assert(f.(*MemoryFile).filename, Equals, "foo")
	c.Assert(client.Files["foo"], Equals, f)
}

func (s *WritersSuite) TestMemoryClient_Purge(c *C) {
	client := NewMemoryClient()

	f, _ := client.Create("foo")
	c.Assert(client.Files, HasLen, 1)

	f.Close()
	client.Purge()
	c.Assert(client.Files, HasLen, 0)
}

func (s *WritersSuite) TestMemoryFile_Write(c *C) {
	f := &MemoryFile{Content: bytes.NewBuffer(nil)}

	i, err := f.Write([]byte("foo"))
	c.Assert(i, Equals, 3)
	c.Assert(err, IsNil)
}

func (s *WritersSuite) TestMemoryFile_Close(c *C) {
	f := &MemoryFile{Content: bytes.NewBuffer(nil)}
	f.Close()

	c.Assert(f.closed, Equals, true)
}

func (s *WritersSuite) TestMemoryFile_CloseWrite(c *C) {
	f := &MemoryFile{Content: bytes.NewBuffer(nil)}
	f.Close()

	i, err := f.Write([]byte("foo"))
	c.Assert(i, Equals, 0)
	c.Assert(err, Equals, ClosedFileError)
}
