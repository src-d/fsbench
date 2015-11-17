package fsbench

import (
	"sync"

	"github.com/src-d/fsbench/fs"
)

type Config struct {
	Workers        int
	Files          int
	Filesystem     string
	DirectoryDepth int
	BlockSize      int64
	MinFileSize    int64
	MaxFileSize    int64
}

type Benchmark struct {
	c *Config
	w []*Worker
}

func NewBenchmark(c *Config) *Benchmark {
	return &Benchmark{c: c}
}

func (b *Benchmark) Init() {
	fs := fs.NewOSClient("/tmp/b/", "/tmp/b/tmp")

	for i := 0; i < b.c.Workers; i++ {
		c := b.getWorkerConfig()
		c.Files = b.c.Files / b.c.Workers

		if i == b.c.Workers-1 {
			c.Files += (b.c.Files % b.c.Workers)
		}

		b.w = append(b.w, NewWorker(fs, c))
	}
}

func (b *Benchmark) getWorkerConfig() *WorkerConfig {
	return &WorkerConfig{
		DirectoryDepth: b.c.DirectoryDepth,
		BlockSize:      b.c.BlockSize,
		MinFileSize:    b.c.MinFileSize,
		MaxFileSize:    b.c.MaxFileSize,
	}
}

func (b *Benchmark) Run() *Status {
	var wg sync.WaitGroup
	for _, w := range b.w {
		wg.Add(1)
		go func(w *Worker) {
			w.Do()
			wg.Done()
		}(w)
	}

	wg.Wait()
	return b.Status()
}

func (b *Benchmark) Status() *Status {
	s := &Status{}
	for _, w := range b.w {
		s.Add(&w.Status)
	}

	return s
}
