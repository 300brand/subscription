package relink

import (
	"bytes"
	"io"
	"launchpad.net/gocheck"
	"testing"
)

type RelinkSuite struct {
	DomainMap map[string]string
	Tests     []relinkTest
}

type relinkTest struct {
	In, Out string
}

var _ = gocheck.Suite(&RelinkSuite{
	DomainMap: map[string]string{
		"top.tld": "top.proxy.tld",
	},
	Tests: []relinkTest{
		{
			In:  `<a href="http://top.tld/">`,
			Out: `<a href="http://alias.subscription.tld/">`,
		},
		{
			In:  `<a href="http://top.tld">`,
			Out: `<a href="http://top.tld">`,
		},
		{
			In:  `<a href="http://top.tld/some/page">`,
			Out: `<a href="http://top.tld/some/page">`,
		},
		{
			In:  `<a href="http://sub.top.tld/some/page">`,
			Out: `<a href="http://sub.top.tld/some/page">`,
		},
		{
			In:  `<a href="http://sub.sub.top.tld/some/page">`,
			Out: `<a href="http://sub.sub.top.tld/some/page">`,
		},
		{
			In:  `<a href="http://sub.top.tld">`,
			Out: `<a href="http://sub.top.tld">`,
		},
		{
			In:  `<link>http://top.tld/</link>`,
			Out: `<link>http://top.tld/</link>`,
		},
		{
			In:  `<link>http://top.tld</link>`,
			Out: `<link>http://top.tld</link>`,
		},
		{
			In:  `<link>http://top.tld/some/page</link>`,
			Out: `<link>http://top.tld/some/page</link>`,
		},
		{
			In:  `<link>http://sub.top.tld/some/page</link>`,
			Out: `<link>http://sub.top.tld/some/page</link>`,
		},
		{
			In:  `<link>http://sub.sub.top.tld/some/page</link>`,
			Out: `<link>http://sub.sub.top.tld/some/page</link>`,
		},
		{
			In:  `<link>http://sub.top.tld</link>`,
			Out: `<link>http://sub.top.tld</link>`,
		},
		{
			In:  `<a href="/some/page">`,
			Out: `<a href="/some/page">`,
		},
		{
			In:  `<link>/some/page</link>`,
			Out: `<link>/some/page</link>`,
		},
	},
})

func Test(t *testing.T) { gocheck.TestingT(t) }

func (s *RelinkSuite) testBuf(c *gocheck.C) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	for _, test := range s.Tests {
		buf.WriteString(test.In)
		buf.WriteByte('\n')
	}

	c.Logf("=== Buffer ===\n%s", buf.Bytes())
}

func (s *RelinkSuite) TestRelink(c *gocheck.C) {
	in := bytes.NewBuffer(make([]byte, 0, 1024))

	for _, test := range s.Tests {
		in.WriteString(test.In)
		in.WriteByte('|')
	}

	outBuf := bytes.NewBuffer(make([]byte, 0, 16))
	out := New(outBuf, s.DomainMap)

	var err error
	written := int64(0)
	buf := make([]byte, 32)
	for {
		nr, er := in.Read(buf)
		if nr > 0 {
			nw, ew := out.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}

	c.Assert(err, gocheck.IsNil)

	c.Logf("=== In ===\n%s\n=== Out ===\n%s\n", in.Bytes(), outBuf.Bytes())
}
