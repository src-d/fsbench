package fsbench

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type CommonSuite struct{}

var _ = Suite(&CommonSuite{})

func (s *CommonSuite) TestAdd(c *C) {
	status := NewAggregatedStatus()
	status.Add(&Status{AvgRate: 42, Bytes: 84})
	status.Add(&Status{AvgRate: 42, Bytes: 84})
	status.Add(&Status{AvgRate: 42, Bytes: 84})
	status.Add(&Status{AvgRate: 42, Bytes: 84})

	c.Assert(status.HistogramAvgRate.Len(), Equals, 4)
	c.Assert(status.HistogramDuration.Len(), Equals, 4)
}
