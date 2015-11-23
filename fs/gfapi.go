package fs

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"

	"github.com/kshlm/gogfapi/gfapi"
)

//GFAPIClient a filesystem based on GlusterFS
type GFAPIClient struct {
	RootDir string
	volume  *gfapi.Volume
}

//NewGFAPIClient returns a new GFAPIClient
func NewGFAPIClient(server, datastore, root string) (*GFAPIClient, error) {
	v := &gfapi.Volume{}
	if r := v.Init(server, datastore); r != 0 {
		return nil, errors.New(fmt.Sprintf("Unable to Init volume %q", server))
	}

	v.SetLogging("log", gfapi.LogInfo)
	if r := v.Mount(); r != 0 {
		return nil, errors.New(fmt.Sprintf("Unable to Mount volume %q", server))
	}

	return &GFAPIClient{
		RootDir: root,
		volume:  v,
	}, nil
}

//Create creates a new GFAPIFile
func (c *GFAPIClient) Create(filename string) (File, error) {
	fullpath := path.Join(c.RootDir, filename)
	if filepath.Dir(filename) != "." {
		dir := filepath.Dir(fullpath)
		if err := c.volume.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	f, err := c.volume.Create(fullpath)
	if err != nil {
		return nil, err
	}

	return &GFAPIFile{
		BaseFile: BaseFile{filename: filename},
		file:     f,
	}, nil
}

func (c *GFAPIClient) Open(filename string) (File, error) {
	fullpath := path.Join(c.RootDir, filename)
	f, err := c.volume.Open(fullpath)

	if err != nil {
		return nil, err
	}

	return &GFAPIFile{
		BaseFile: BaseFile{filename: filename},
		file:     f,
	}, nil
}

type GFAPIFile struct {
	file *gfapi.File
	BaseFile
}

func (f *GFAPIFile) Write(p []byte) (int, error) {
	return f.file.Write(p)
}

func (f *GFAPIFile) Close() error {
	f.closed = true
	return f.file.Close()
}
