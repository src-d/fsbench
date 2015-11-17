package fs

import (
	"io/ioutil"
	"os"

	. "gopkg.in/check.v1"
)

func (s *WritersSuite) TestOSClient_Create(c *C) {
	temp := getTempDir()
	path := getTempDir()

	client := NewOSClient(path, temp)
	f, _ := client.Create("foo")
	c.Assert(f.(*OSFile).file.Name(), Not(Equals), f.GetFilename())
}

func (s *WritersSuite) TestOSClient_Write(c *C) {
	temp := getTempDir()
	path := getTempDir()

	client := NewOSClient(path, temp)
	f, _ := client.Create("foo")

	l, err := f.Write([]byte("foo"))
	c.Assert(l, Equals, 3)
	c.Assert(err, IsNil)

	wrote, _ := ioutil.ReadFile(f.(*OSFile).file.Name())
	c.Assert(wrote, DeepEquals, []byte("foo"))
}

func (s *WritersSuite) TestOSClient_Close(c *C) {
	temp := getTempDir()
	path := getTempDir()

	client := NewOSClient(path, temp)
	f, _ := client.Create("foo")

	f.Write([]byte("foo"))
	c.Assert(f.Close(), IsNil)

	//The temp file should be removed
	_, err := ioutil.ReadFile(f.(*OSFile).file.Name())
	c.Assert(err, Not(IsNil))

	//and the final file should be there
	wrote, _ := ioutil.ReadFile(f.GetFilename())
	c.Assert(wrote, DeepEquals, []byte("foo"))
}

func getTempDir() string {
	dir, _ := ioutil.TempDir(os.TempDir(), "--OSClientTest--")
	return dir
}
