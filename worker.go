package fsbench

import (
	"bytes"
	"fmt"
	"math/rand"
	"path/filepath"

	"github.com/mxk/go-flowrate/flowrate"
	"github.com/src-d/fsbench/fs"
)

const (
	DirectoryLength = 2
	FilenameLength  = 32
	FileExtension   = ".zero"
)

type WorkerConfig struct {
	Files          int
	DirectoryDepth int
	BlockSize      int64
	MinFileSize    int64
	MaxFileSize    int64
}

type Worker struct {
	c     *WorkerConfig
	fs    fs.Client
	block []byte

	Status Status
}

func NewWorker(fs fs.Client, c *WorkerConfig) *Worker {
	return &Worker{
		c:     c,
		fs:    fs,
		block: bytes.Repeat([]byte{0}, int(c.BlockSize)),
	}
}

func (w *Worker) Do() error {
	for i := 0; i < w.c.Files; i++ {
		if err := w.doCreate(); err != nil {
			w.Status.Errors++
			fmt.Println(err)
			continue
		}
	}

	return nil
}

func (w *Worker) doCreate() error {
	file, err := w.fs.Create(w.getFilename())
	if err != nil {
		return err
	}

	flow := flowrate.NewWriter(file, -1)
	var size int64
	expected := w.getSize()
	for {
		s, err := flow.Write(w.block)
		if err != nil {
			return err
		}

		size += int64(s)
		if size >= expected {
			break
		}
	}

	flow.Close()

	w.done(flow.Status())
	return nil
}

func (w *Worker) done(s flowrate.Status) {
	w.Status.Bytes += s.Bytes
	w.Status.Duration += s.Duration
	w.Status.Samples += s.Samples
	w.Status.AvgRate = int64(float64(w.Status.Bytes) / w.Status.Duration.Seconds())
	w.Status.Files++

	if s.PeakRate > w.Status.PeakRate {
		w.Status.PeakRate = s.PeakRate
	}
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
	if w.c.MinFileSize == w.c.MaxFileSize {
		return w.c.MinFileSize
	}

	diff := (w.c.MaxFileSize - w.c.MinFileSize) / w.c.BlockSize
	r := rand.Int63n(diff)

	return w.c.MinFileSize + (r * w.c.BlockSize)
}
