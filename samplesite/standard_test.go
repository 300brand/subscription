package samplesite

import (
	"launchpad.net/gocheck"
	"net/http"
	"net/url"
)

type StandardSuite struct {
	URL       string
	GoodLogin url.Values
	BadLogin  url.Values
}

var _ = gocheck.Suite(new(StandardSuite))

func (s *StandardSuite) SetUpSuite(c *gocheck.C) {
	s.URL = Start()
	s.GoodLogin = url.Values{
		"username": []string{Username},
		"password": []string{Password},
	}
	s.BadLogin = url.Values{
		"username": []string{Username + "BAD"},
		"password": []string{Password + "BAD"},
	}
}

func (s *StandardSuite) TearDownSuite(c *gocheck.C) {
	Close()
}

func (s *StandardSuite) TestNoLogin(c *gocheck.C) {
	resp, err := http.Get(s.URL + "/loggedin.check")
	c.Assert(err, gocheck.IsNil)
	defer resp.Body.Close()
	buf := make([]byte, 1)
	n, err := resp.Body.Read(buf)
	c.Assert(err, gocheck.IsNil)
	c.Assert(n, gocheck.Equals, 1)
	c.Assert(buf[0], gocheck.Equals, byte('0'))
}

func (s *StandardSuite) TestBadLoginRedirect(c *gocheck.C) {
	resp, err := http.PostForm(s.URL+"/standard/login.do", s.BadLogin)
	c.Assert(err, gocheck.IsNil)
	defer resp.Body.Close()

	u, err := resp.Location()
	c.Assert(err, gocheck.IsNil)
	c.Assert(u.Path, gocheck.Equals, "/standard/failed")
}

func (s *StandardSuite) TestGoodLoginRedirect(c *gocheck.C) {
	resp, err := http.PostForm(s.URL+"/standard/login.do", s.GoodLogin)
	c.Assert(err, gocheck.IsNil)
	defer resp.Body.Close()

	u, err := resp.Location()
	c.Assert(err, gocheck.IsNil)
	c.Assert(u.Path, gocheck.Equals, "/standard/success")
}

func (s *StandardSuite) TestGoodLoginCookie(c *gocheck.C) {
	resp, err := http.PostForm(s.URL+"/standard/login.do", s.GoodLogin)
	c.Assert(err, gocheck.IsNil)
	resp.Body.Close()

	cookies := resp.Cookies()
	c.Assert(len(cookies), gocheck.Equals, 1)

	req, err := http.NewRequest("GET", s.URL+"/loggedin.check", nil)
	c.Assert(err, gocheck.IsNil)
	req.AddCookie(cookies[0])

	resp, err = http.DefaultClient.Do(req)
	c.Assert(err, gocheck.IsNil)
	defer resp.Body.Close()
	buf := make([]byte, 1)
	n, err := resp.Body.Read(buf)
	c.Assert(err, gocheck.IsNil)
	c.Assert(n, gocheck.Equals, 1)
	c.Assert(buf[0], gocheck.Equals, byte('1'))
}

func (s *StandardSuite) TestBadLoginCookie(c *gocheck.C) {
	resp, err := http.PostForm(s.URL+"/standard/login.do", s.BadLogin)
	c.Assert(err, gocheck.IsNil)
	resp.Body.Close()

	cookies := resp.Cookies()
	c.Assert(len(cookies), gocheck.Equals, 1)

	req, err := http.NewRequest("GET", s.URL+"/loggedin.check", nil)
	c.Assert(err, gocheck.IsNil)
	req.AddCookie(cookies[0])

	resp, err = http.DefaultClient.Do(req)
	c.Assert(err, gocheck.IsNil)
	defer resp.Body.Close()
	buf := make([]byte, 1)
	n, err := resp.Body.Read(buf)
	c.Assert(err, gocheck.IsNil)
	c.Assert(n, gocheck.Equals, 1)
	c.Assert(buf[0], gocheck.Equals, byte('0'))
}
