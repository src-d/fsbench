package fs

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type WritersSuite struct{}

var _ = Suite(&WritersSuite{})
