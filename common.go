package fsbench

import (
	"math/rand"
	"time"

	"github.com/dripolles/histogram"
)

type Status struct {
	Duration time.Duration // Time period covered by the statistics
	Bytes    int64         // Total number of bytes transferred
	Samples  int64         // Total number of samples taken
	AvgRate  int64         // Average transfer rate (Bytes / Duration)
	PeakRate int64         // Maximum instantaneous transfer rate
	Files    int           // Number of files transferred
	Errors   int           // Number of errors
}

type AggregatedStatus struct {
	Status
	HistogramAvgRate  *histogram.Histogram
	HistogramDuration *histogram.Histogram
}

func NewAggregatedStatus() *AggregatedStatus {
	return &AggregatedStatus{
		HistogramAvgRate:  histogram.NewHistogram(),
		HistogramDuration: histogram.NewHistogram(),
	}
}

func (s *AggregatedStatus) Add(a *Status) {
	s.Files += a.Files
	s.Errors += a.Errors
	s.Bytes += a.Bytes
	s.Duration += a.Duration
	s.Samples += a.Samples
	s.AvgRate = int64(float64(s.Bytes) / s.Duration.Seconds())

	if a.PeakRate > s.PeakRate {
		s.PeakRate = s.PeakRate
	}

	s.HistogramAvgRate.Add(int(a.AvgRate))
	s.HistogramDuration.Add(int(a.Duration))
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// radomString generates a random string of any length, code extracted from:
// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func randomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
