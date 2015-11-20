package fsbench

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type CommonSuite struct{}

var _ = Suite(&CommonSuite{})

func (s *CommonSuite) TestAdd(c *C) {
	status := NewAggregatedStatus()
	status.Add(Status{AvgRate: 42, Duration: time.Duration(1)})
	status.Add(Status{AvgRate: 42, Duration: time.Duration(21)})
	status.Add(Status{AvgRate: 42, Duration: time.Duration(42)})
	status.Add(Status{AvgRate: 42, Duration: time.Duration(84)})

	c.Assert(status.HistogramAvgRate.Len(), Equals, 4)
	c.Assert(status.HistogramDuration.Len(), Equals, 4)
	c.Assert(status.HistogramDuration.GetAtPercentile(.9), Equals, 84)
}
