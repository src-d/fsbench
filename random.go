package fsbench

import "crypto/rand"

var (
	RandomSampleSize int64 = 10 * MB
	randomSample     []byte
)

func init() {
	randomSample = make([]byte, RandomSampleSize)
	_, err := rand.Reader.Read(randomSample)
	if err != nil {
		panic(err)
	}
}

type RandomReader struct {
	pos   int64
	size  int64
	bytes []byte
}

func NewRandomReader() *RandomReader {
	return &RandomReader{bytes: randomSample, size: RandomSampleSize}
}
func (r *RandomReader) Read(b []byte) (int, error) {
	size := len(b)
	needed := size
	for {
		copied := copy(b[size-needed:], r.bytes[r.pos:])
		r.pos += int64(copied)
		if r.pos >= r.size {
			r.pos = 0
		}

		needed -= copied
		if needed <= 0 {
			break
		}
	}

	return size - needed, nil
}
