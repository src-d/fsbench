package fsbench

import . "gopkg.in/check.v1"

type RandomSuite struct{}

var _ = Suite(&RandomSuite{})

func (s *RandomSuite) TestInit(c *C) {
	c.Assert(randomSample, HasLen, 52428800)
}

func (s *RandomSuite) TestRead(c *C) {
	test := []byte("0123456789")
	r := &RandomReader{bytes: test, size: 10}

	ta := make([]byte, 6)
	l, _ := r.Read(ta)
	c.Assert(l, Equals, 6)
	c.Assert(string(ta), Equals, "012345")

	tb := make([]byte, 6)
	l, _ = r.Read(tb)
	c.Assert(l, Equals, 6)
	c.Assert(string(tb), Equals, "678901")

	tc := make([]byte, 32)
	l, _ = r.Read(tc)
	c.Assert(l, Equals, 32)
	c.Assert(string(tc), Equals, "23456789012345678901234567890123")
}
