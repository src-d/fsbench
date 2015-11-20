package fsbench

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"path/filepath"

	"github.com/mxk/go-flowrate/flowrate"
	"github.com/src-d/fsbench/fs"
)

const (
	DirectoryLength = 2
	FilenameLength  = 32
	FileExtension   = ".rand"
)

type WorkerConfig struct {
	Files          int
	DirectoryDepth int
	BlockSize      int64
	FixedFileSize  int64
	MeanFileSize   int64
	StdDevFileSize float64
}

type Worker struct {
	c         *WorkerConfig
	fs        fs.Client
	r         *RandomReader
	filenames []string

	WStatus *AggregatedStatus
	RStatus *AggregatedStatus
}

func NewWorker(fs fs.Client, c *WorkerConfig) *Worker {
	return &Worker{
		c:       c,
		fs:      fs,
		r:       NewRandomReader(),
		WStatus: NewAggregatedStatus(),
		RStatus: NewAggregatedStatus(),
	}
}

func (w *Worker) Write() error {
	for i := 0; i < w.c.Files; i++ {
		if err := w.doWrite(); err != nil {
			w.WStatus.Errors++
			fmt.Println(err)
			continue
		}
	}

	return nil
}

func (w *Worker) Read() error {
	for _, filename := range w.filenames {
		if err := w.doRead(filename); err != nil {
			w.RStatus.Errors++
			fmt.Println(err)
			continue
		}

	}

	return nil
}

func (w *Worker) doWrite() error {
	file, err := w.fs.Create(w.getFilename())
	if err != nil {
		return err
	}

	flow := flowrate.NewWriter(file, -1)
	var size int64
	expected := w.getSize()
	for {
		s, err := copyN(flow, w.r, w.c.BlockSize)
		if err != nil {
			return err
		}

		size += int64(s)
		if size >= expected {
			break
		}
	}

	w.filenames = append(w.filenames, file.GetFilename())

	flow.Close()
	w.WStatus.Add(NewStatus(flow.Status()))

	return nil
}

func (w *Worker) doRead(filename string) error {
	file, err := w.fs.Open(filename)
	if err != nil {
		return err
	}

	flow := flowrate.NewReader(file, -1)
	for {
		_, err := copyN(ioutil.Discard, flow, w.c.BlockSize)
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}
	}

	flow.Close()
	w.RStatus.Add(NewStatus(flow.Status()))

	return nil
}

func (w *Worker) getFilename() string {
	r := randomString(FilenameLength)

	offset := 0
	for i := 0; i < w.c.DirectoryDepth; i++ {
		cut := offset + DirectoryLength
		if cut > len(r) {
			return r
		}

		r = filepath.Join(r[:cut], r[cut:])
		offset += DirectoryLength + 1
	}

	return r + FileExtension
}

func (w *Worker) getSize() int64 {
	if w.c.FixedFileSize != 0 {
		return w.c.FixedFileSize
	}

	return int64(normFloat64()*w.c.StdDevFileSize + float64(w.c.MeanFileSize))
}

func copyN(dst io.Writer, src io.Reader, amount int64) (int, error) {
	bytes := make([]byte, amount)
	if _, err := src.Read(bytes); err != nil {
		return 0, err
	}

	s, err := dst.Write(bytes)
	if err != nil {
		return s, err
	}

	return s, nil
}

func normFloat64() float64 {
	for {
		r := rand.NormFloat64()
		if r <= -3 || r >= 3 {
			continue
		}

		return r
	}
}
