package fsbench

import (
	"sync"

	"github.com/src-d/fsbench/fs"
)

type WorkerMode int

const (
	WriteMode WorkerMode = 1
	ReadMode  WorkerMode = 2
)

type Config struct {
	Mode           WorkerMode
	Workers        int
	Files          int
	Filesystem     string
	DirectoryDepth int
	BlockSize      int64
	FixedFileSize  int64
	MeanFileSize   int64
	StdDevFileSize float64
}

type Benchmark struct {
	c *Config
	w []*Worker
}

func NewBenchmark(c *Config) *Benchmark {
	return &Benchmark{c: c}
}

func (b *Benchmark) Init() {
	fs := fs.NewMemoryClient()

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
		FixedFileSize:  b.c.FixedFileSize,
		MeanFileSize:   b.c.MeanFileSize,
		StdDevFileSize: b.c.StdDevFileSize,
	}
}

func (b *Benchmark) Run() *BenchmarkStatus {
	var wg sync.WaitGroup
	for _, w := range b.w {
		wg.Add(1)
		go func(w *Worker) {
			w.Write()
			wg.Done()
		}(w)
	}

	wg.Wait()

	if b.c.Mode == WriteMode {
		return b.Status()
	}

	for _, w := range b.w {
		wg.Add(1)
		go func(w *Worker) {
			w.Read()
			wg.Done()
		}(w)
	}

	wg.Wait()

	return b.Status()
}

func (b *Benchmark) Status() *BenchmarkStatus {
	s := NewBenchmarkStatus()
	for _, w := range b.w {
		s.RStatus.Sum(w.RStatus)
		s.WStatus.Sum(w.WStatus)
	}

	return s
}

type BenchmarkStatus struct {
	RStatus *AggregatedStatus
	WStatus *AggregatedStatus
}

func NewBenchmarkStatus() *BenchmarkStatus {
	return &BenchmarkStatus{
		RStatus: NewAggregatedStatus(),
		WStatus: NewAggregatedStatus(),
	}
}
