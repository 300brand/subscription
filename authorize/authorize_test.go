package authorize

import (
	"github.com/300brand/subscription/config"
	"launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

func runTest(c *gocheck.C, a Authorizer, cfg *config.Domain) {
	var err error
	loggedIn, err := a.LoggedIn(cfg)
	c.Assert(err, gocheck.IsNil)
	c.Assert(loggedIn, gocheck.Equals, false)

	err = a.Login(cfg)
	c.Assert(err, gocheck.IsNil)

	loggedIn, err = a.LoggedIn(cfg)
	c.Assert(err, gocheck.IsNil)
	c.Assert(loggedIn, gocheck.Equals, true)
}
