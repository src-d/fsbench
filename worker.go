package fsbench

import (
	"fmt"
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
	c      *WorkerConfig
	fs     fs.Client
	r      *RandomReader
	Status Status
}

func NewWorker(fs fs.Client, c *WorkerConfig) *Worker {
	return &Worker{c: c, fs: fs, r: NewRandomReader()}
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
		bytes := make([]byte, w.c.BlockSize)
		if _, err := w.r.Read(bytes); err != nil {
			return err
		}

		s, err := flow.Write(bytes)
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
	if w.c.FixedFileSize != 0 {
		return w.c.FixedFileSize
	}

	return int64(normFloat64()*w.c.StdDevFileSize + float64(w.c.MeanFileSize))
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
