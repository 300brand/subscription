package config

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Config struct {
	Alias     string       // Site alias, also used for the subdomain
	Domain    string       // Site URL root
	LoginType string       // Login type, must be one of the available types
	Username  string       // Account username
	Password  string       // Account password
	URLs      URLs         // URLs used to perform authentication
	url       *url.URL     // Parsed Siteroot
	client    *http.Client // HTTP client with cookie jar already set
}

type URLs struct {
	Form    string // Path to login form (hit to set baseline cookies)
	Do      string // Path to page that does the logging in
	Success string // Path to success page
	Failure string // Path to failure page (bad login)
	Check   string // URL used to check if currently logged in
}

func (c *Config) ResolveReference(refurl string) *url.URL {
	if c.url == nil {
		c.url, _ = url.Parse(c.Domain)
	}
	ref, _ := url.Parse(refurl)
	return c.url.ResolveReference(ref)
}

func (c *Config) Client() *http.Client {
	if c.client == nil {
		c.client = new(http.Client)
		c.client.Jar, _ = cookiejar.New(nil)
	}
	return c.client
}
