package subscription

import (
	"github.com/300brand/logger"
	"net/http"
	"net/url"
)

type CookieJar struct {
	cookies []http.Cookie
}

var _ http.CookieJar = new(CookieJar)

func NewCookieJar() (jar *CookieJar) {
	return &CookieJar{
		cookies: make([]http.Cookie, 0, 4096),
	}
}

func (jar *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	logger.Debug.Printf("CookieJar.SetCookie: %s (%d)", u, len(cookies))
	for _, cookie := range cookies {
		jar.cookies = append(jar.cookies, *cookie)
	}
}

func (jar *CookieJar) Cookies(u *url.URL) (cookies []*http.Cookie) {
	logger.Debug.Printf("CookieJar.Cookies: %s", u)
	for i, c := range jar.cookies {
		logger.Debug.Printf("[%03d] %s=%s; Expires=%s; Path=%s", i, c.Name, c.Value, c.Expires, c.Path)
	}
	return []*http.Cookie{}
}
