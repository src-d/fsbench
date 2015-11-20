package fs

import (
	"os"
	"path"
)

//OSClient a filesystem based on OSClient
type OSClient struct {
	RootDir string
}

//NewOSClient returns a new OSClient
func NewOSClient(rootDir string) *OSClient {
	return &OSClient{RootDir: rootDir}
}

//Create creates a new OSFile
func (c *OSClient) Create(filename string) (File, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL | os.O_SYNC
	return c.open(filename, flags)
}

//Open opens a new OSFile
func (c *OSClient) Open(filename string) (File, error) {
	flags := os.O_RDONLY | os.O_SYNC
	return c.open(filename, flags)
}

func (c *OSClient) open(filename string, flags int) (File, error) {
	fullpath := path.Join(c.RootDir, filename)
	f, err := os.OpenFile(fullpath, flags, 0644)
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
