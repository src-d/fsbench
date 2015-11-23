package fs

import (
	"errors"
	"io"
	"net/http"
	"path"
	"sync"

	"github.com/ginuerzh/weedo"
)

type SeaweedFSClient struct {
	RootDir string
	client  *weedo.Client
	fids    map[string]string
	lock    sync.Mutex
}

//NewSeaweedFSClient returns a new SeaweedFSClient
func NewSeaweedFSClient(server, rootDir string) *SeaweedFSClient {
	return &SeaweedFSClient{
		RootDir: rootDir,
		client:  weedo.NewClient(server),
		fids:    make(map[string]string, 0),
	}
}

//Create creates a new SeaweedFSFile
func (c *SeaweedFSClient) Create(filename string) (File, error) {
	fullpath := path.Join(c.RootDir, filename)

	r, w := io.Pipe()

	errChan := make(chan error, 1)
	go func(filename string) {
		fid, _, err := c.client.AssignUpload(fullpath, "application/octet-stream", r)
		errChan <- err
		close(errChan)

		c.lock.Lock()
		defer c.lock.Unlock()
		c.fids[fullpath] = fid
	}(fullpath)

	return &SeaweedFSFile{
		BaseFile: BaseFile{filename: fullpath},
		w:        w,
		errChan:  errChan,
	}, nil
}

//Open opens a new SeaweedFSFile
func (c *SeaweedFSClient) Open(filename string) (File, error) {
	fullpath := path.Join(c.RootDir, filename)

	fid, ok := c.fids[fullpath]
	if !ok {
		return nil, errors.New("file not found")
	}

	url, _, err := c.client.GetUrl(fid)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return &SeaweedFSFile{
		BaseFile: BaseFile{filename: fullpath},
		r:        resp.Body,
	}, nil
}

func (c *SeaweedFSClient) Stat(filename string) (FileInfo, error) {
	return nil, nil
}

type SeaweedFSFile struct {
	BaseFile
	w       io.WriteCloser
	r       io.Reader
	errChan chan error
}

func (f *SeaweedFSFile) Read(p []byte) (int, error) {
	return f.r.Read(p)
}

func (f *SeaweedFSFile) Write(p []byte) (int, error) {
	return f.w.Write(p)
}

func (f *SeaweedFSFile) Close() error {
	f.closed = true

	if f.r != nil {
		return nil
	}

	if err := f.w.Close(); err != nil {
		return err
	}

	return <-f.errChan
}
