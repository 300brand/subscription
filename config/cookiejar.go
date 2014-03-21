package config

import (
	"github.com/300brand/logger"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var _ http.CookieJar = new(CookieJar)

type CookieJar struct {
	jar *cookiejar.Jar
}

func NewCookieJar() (jar *CookieJar) {
	jar = new(CookieJar)
	jar.jar, _ = cookiejar.New(nil)
	return
}

func (jar *CookieJar) Cookies(u *url.URL) []*http.Cookie {
	cookies := jar.jar.Cookies(u)
	for i, c := range cookies {
		logger.Info.Printf("<- [%d] U:%s %s=%v", i, u.Host, c.Name, c.Value)
	}
	return cookies
}

func (jar *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for i, c := range cookies {
		logger.Warn.Printf("-> [%d] U:%s %s", i, u.Host, c.Raw)
	}
	jar.jar.SetCookies(u, cookies)
}
