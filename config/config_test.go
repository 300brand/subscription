package config

import (
	"github.com/300brand/subscription/samplesite"
	"launchpad.net/gocheck"
	"testing"
)

type ConfigSuite struct{}

var _ = gocheck.Suite(new(ConfigSuite))

func Test(t *testing.T) { gocheck.TestingT(t) }

func (s *ConfigSuite) TestResolveReference(c *gocheck.C) {
	cfg := Domain{Domain: "http://test.com"}
	c.Assert(cfg.ResolveReference("/login").String(), gocheck.Equals, "http://test.com/login")
	// Ensure the reference didn't stick
	c.Assert(cfg.ResolveReference(".").String(), gocheck.Equals, "http://test.com/")
}

func (s *ConfigSuite) TestClient(c *gocheck.C) {
	url := samplesite.Start()
	defer samplesite.Close()

	cfg := Domain{Domain: url}

	resp, err := cfg.Client().Get(url + "/loggedin.check")
	c.Assert(err, gocheck.IsNil)
	defer resp.Body.Close()

	buf := make([]byte, 1)
	n, err := resp.Body.Read(buf)
	c.Assert(err, gocheck.IsNil)
	c.Assert(n, gocheck.Equals, len(buf))
	c.Assert(buf[0], gocheck.Equals, byte('0'))
}
