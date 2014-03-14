package relink

import (
	"launchpad.net/gocheck"
	"strings"
	"testing"
)

type RelinkSuite struct{}

var _ = gocheck.Suite(new(RelinkSuite))

func Test(t *testing.T) { gocheck.TestingT(t) }

func (s *RelinkSuite) TestRelink(c *gocheck.C) {
	r := strings.NewReader(`
		
	`)
}
