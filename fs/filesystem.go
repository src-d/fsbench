package fs

import (
	"errors"
	"io"
	"os"
)

var (
	ClosedFileError       = errors.New("File: Writing on closed file.")
	ReadOnlyFileSystemErr = errors.New("this is a read-only filesystem")
)

type Client interface {
	Create(filename string) (File, error)
	Open(filename string) (File, error)
	Stat(filename string) (FileInfo, error)
}

type File interface {
	GetFilename() string
	io.Writer
	io.Reader
	io.Closer
}

type FileInfo os.FileInfo

type BaseFile struct {
	filename string
	closed   bool
}

//GetFilename returns the filename from the File
func (f *BaseFile) GetFilename() string {
	return f.filename
}

//IsClosed returns if te file is closed
func (f *BaseFile) IsClosed() bool {
	return f.closed
}
