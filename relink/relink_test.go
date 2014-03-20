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
		"top.tld": "alias.subscription.tld",
	},
	Tests: []relinkTest{
		{
			In:  `<a href="http://top.tld/">`,
			Out: `<a href="http://alias.subscription.tld/">`,
		},
		{
			In:  `<a href="http://top.tld">`,
			Out: `<a href="http://alias.subscription.tld">`,
		},
		{
			In:  `<a href="http://top.tld/some/page">`,
			Out: `<a href="http://alias.subscription.tld/some/page">`,
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
			Out: `<link>http://alias.subscription.tld/</link>`,
		},
		{
			In:  `<link>http://top.tld</link>`,
			Out: `<link>http://alias.subscription.tld</link>`,
		},
		{
			In:  `<link>http://top.tld/some/page</link>`,
			Out: `<link>http://alias.subscription.tld/some/page</link>`,
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
	inBuf := bytes.NewBuffer(make([]byte, 0, 1024))
	sep := []byte{'\n'}

	for _, test := range s.Tests {
		inBuf.WriteString(test.In)
		inBuf.WriteByte(sep[0])
	}

	out := bytes.NewBuffer(make([]byte, 0, 16))
	in := New(inBuf, s.DomainMap)

	_, err := io.Copy(out, in)
	c.Assert(err, gocheck.IsNil)

	outLines := bytes.Split(out.Bytes(), sep)
	for i := range s.Tests {
		c.Check(string(outLines[i]), gocheck.Equals, s.Tests[i].Out)
	}
}
