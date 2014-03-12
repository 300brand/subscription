package authorize

import (
	"github.com/300brand/subscription/config"
	"github.com/300brand/subscription/samplesite"
	"launchpad.net/gocheck"
)

type StandardSuite struct {
	URL string
}

var _ = gocheck.Suite(new(StandardSuite))

func (s *StandardSuite) SetUpSuite(c *gocheck.C) {
	s.URL = samplesite.Start()
}

func (s *StandardSuite) TearDownSuite(c *gocheck.C) {
	samplesite.Close()
}

func (s *StandardSuite) TestSamplesiteLogin(c *gocheck.C) {
	cfg := &config.Domain{
		Domain:   s.URL,
		Username: [2]string{"username", "username"},
		Password: [2]string{"password", "password"},
		Cookie:   [2]string{"standardLogin", ""},
		URLs: config.URLs{
			Form:    "/standard/login.form",
			Do:      "/standard/login.do",
			Success: "/standard/success",
		},
	}
	runTest(c, new(Standard), cfg)
}
