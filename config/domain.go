package config

import (
	"errors"
	"net/http"
	"net/url"
	"sort"
)

type Domains []*Domain

type Domain struct {
	Alias     string       // Site alias, also used for the subdomain
	Domain    string       // Site URL root
	Rewrite   []string     // Domains to rewrite to the subscription URL
	LoginType string       // Login type, must be one of the available types
	Username  [2]string    // Account username {field, value}
	Password  [2]string    // Account password {field, value}
	Cookie    [2]string    // Cookie to look for {name, value}, blank = any
	Extra     url.Values   // Extra fields to submit
	URLs      URLs         // URLs used to perform authentication
	url       *url.URL     // Parsed Siteroot
	client    *http.Client // HTTP client with cookie jar already set
}

type URLs struct {
	Form    string // Path to login form (hit to set baseline cookies)
	Do      string // Path to page that does the logging in
	Success string // Path to success page
}

var ErrNotFound = errors.New("Not found")

var _ sort.Interface = Domains{}

func (d Domains) Len() int           { return len(d) }
func (d Domains) Less(i, j int) bool { return d[i].Alias < d[j].Alias }
func (d Domains) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func (d Domains) Get(alias string) (*Domain, error) {
	i := sort.Search(d.Len(), func(n int) bool { return d[n].Alias >= alias })
	if i < d.Len() && d[i].Alias == alias {
		return d[i], nil
	}
	return nil, ErrNotFound
}

func (d *Domain) Client() *http.Client {
	if d.client == nil {
		d.client = new(http.Client)
		d.client.Jar = NewCookieJar()
	}
	return d.client
}

func (d *Domain) ResolveReference(refurl string) *url.URL {
	ref, _ := url.Parse(refurl)
	return d.URL().ResolveReference(ref)
}

func (d *Domain) URL() *url.URL {
	if d.url == nil {
		d.url, _ = url.Parse(d.Domain)
	}
	return d.url
}
