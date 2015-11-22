package fs

import (
	"os"
	"path"

	"github.com/ncw/directio"
)

//OSClient a filesystem based on OSClient
type OSClient struct {
	Direct  bool
	RootDir string
}

//NewOSClient returns a new OSClient
func NewOSClient(rootDir string, direct bool) *OSClient {
	return &OSClient{RootDir: rootDir, Direct: direct}
}

//Create creates a new OSFile
func (c *OSClient) Create(filename string) (File, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
	return c.open(filename, flags)
}

//Open opens a new OSFile
func (c *OSClient) Open(filename string) (File, error) {
	flags := os.O_RDONLY
	return c.open(filename, flags)
}

func (c *OSClient) open(filename string, flags int) (File, error) {
	fullpath := path.Join(c.RootDir, filename)

	var f *os.File
	var err error
	if c.Direct {
		f, err = directio.OpenFile(fullpath, flags, 0666)
	} else {
		f, err = os.OpenFile(fullpath, flags, 0644)
	}

	if err != nil {
		return nil, err
	}

	return &OSFile{
		BaseFile: BaseFile{filename: fullpath},
		file:     f,
	}, nil
}

func (c *OSClient) Stat(filename string) (FileInfo, error) {
	fullpath := path.Join(c.RootDir, filename)

	return os.Stat(fullpath)
}

type OSFile struct {
	file *os.File
	BaseFile
}

func (f *OSFile) Read(p []byte) (int, error) {
	return f.file.Read(p)
}

func (f *OSFile) Write(p []byte) (int, error) {
	return f.file.Write(p)
}

func (f *OSFile) Close() error {
	f.closed = true

	//f.file.Sync()
	return f.file.Close()
}
