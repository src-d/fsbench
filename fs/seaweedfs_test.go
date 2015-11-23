package fs

import (
	"io/ioutil"

	. "gopkg.in/check.v1"
)

func (s *WritersSuite) TestSeaweedFSClient_Create(c *C) {
	client := NewSeaweedFSClient("localhost:9333", "/")
	f, _ := client.Create("foo")

	i, err := f.Write([]byte("foo"))
	c.Assert(i, Equals, 3)
	c.Assert(err, IsNil)

	c.Assert(f.Close(), IsNil)
}

func (s *WritersSuite) TestSeaweedFSClient_Open(c *C) {
	client := NewSeaweedFSClient("localhost:9333", "/")
	f, _ := client.Create("foo")
	f.Write([]byte("foo"))
	f.Close()

	rf, err := client.Open("foo")
	c.Assert(err, IsNil)

	bytes, err := ioutil.ReadAll(rf)
	c.Assert(err, IsNil)
	c.Assert(string(bytes), Equals, "foo")

	c.Assert(f.Close(), IsNil)
}
