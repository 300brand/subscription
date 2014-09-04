package subscription

import (
	"fmt"
	"github.com/300brand/logger"
	"net/http"
	"net/url"
)

type CookieJar struct {
	cookies map[string]http.Cookie
}

var _ http.CookieJar = new(CookieJar)

func NewCookieJar() (jar *CookieJar) {
	return &CookieJar{
		cookies: make(map[string]http.Cookie),
	}
}

func (jar *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	logger.Debug.Printf("CookieJar.SetCookie: %s (%d)", u, len(cookies))
	for _, c := range cookies {
		key := fmt.Sprintf("%s|%s|%s", c.Domain, c.Path, c.Name)
		jar.cookies[key] = *c
	}
}

func (jar *CookieJar) Cookies(u *url.URL) (cookies []*http.Cookie) {
	logger.Debug.Printf("CookieJar.Cookies: %s", u)
	cookies = make([]*http.Cookie, 0, len(jar.cookies))
	for i, c := range jar.cookies {
		logger.Debug.Printf("[%s] %s=%s; Expires=%s; Path=%s", i, c.Name, c.Value, c.RawExpires, c.Path)
		cookies = append(cookies, &c)
	}
	return
}
