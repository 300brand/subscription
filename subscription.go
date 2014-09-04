package subscription

import (
	"fmt"
	"github.com/300brand/logger"
	"github.com/elazarl/goproxy"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Config struct {
	Ident string
	Post  url.Values
	URL   struct {
		PrimeCookies string
		Form         string
		FormAction   string
		BadLogin     string
		GoodLogin    string
	}
}

type Handler struct {
	CookieJar *CookieJar
}

type Subscription struct {
	*goproxy.ProxyHttpServer
	Handler *Handler
}

var _ goproxy.ReqHandler = new(Handler)
var _ goproxy.HttpsHandler = new(Handler)
var Start = time.Now()

func New() (s *Subscription) {
	server := goproxy.NewProxyHttpServer()
	server.Logger = logger.Trace
	server.Verbose = false

	h := &Handler{
		CookieJar: NewCookieJar(),
	}
	server.OnRequest(s.WatchingOrigin()).Do(h)
	server.OnRequest(s.WatchingOrigin()).HandleConnect(h)

	s = &Subscription{
		Handler:         h,
		ProxyHttpServer: server,
	}

	return
}

func (s *Subscription) WatchingOrigin() goproxy.ReqConditionFunc {
	watching := []string{
		"washingtonpost.com",
		"google.com",
	}
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		host := req.Host
		if strings.ContainsRune(host, ':') {
			var err error
			if host, _, err = net.SplitHostPort(host); err != nil {
				logger.Error.Printf("[%03d] net.SplitHostPort: %s", ctx.Session, err)
			}
		}
		for _, w := range watching {
			if strings.HasSuffix(host, w) {
				logger.Info.Printf("[%03d] Matched %s against %s", ctx.Session, req.Host, w)
				return true
			}
		}
		return false
	}
}

// ReqHandler will "tamper" with the request coming to the proxy server
//
// If Handle returns req,nil the proxy will send the returned request to the
// destination server.
//
// If it returns nil,resp the proxy will skip sending any requests, and will
// simply return the response `resp` to the client.
func (h *Handler) Handle(reqIn *http.Request, ctx *goproxy.ProxyCtx) (reqOut *http.Request, respOut *http.Response) {
	logger.Info.Printf("[%04d] Handle %s", ctx.Session, reqIn.URL)

	path := fmt.Sprintf("/tmp/subscription/%s", Start.Format("2006-01-02T15.04.05"))
	if err := os.MkdirAll(path, 0755); err != nil {
		logger.Error.Fatalf("os.MkdirAll: %s", err)
	}
	out, err := os.Create(fmt.Sprintf("%s/%04d", path, ctx.Session))
	if err != nil {
		logger.Error.Fatalf("os.Create: %s", err)
	}

	client := new(http.Client)
	client.Jar = h.CookieJar
	respOut, err = client.Get(reqIn.URL.String())
	if err != nil {
		logger.Error.Printf("[%03d] ERROR %s - %s", ctx.Session, reqIn.URL, err)
		respOut = goproxy.NewResponse(reqIn, goproxy.ContentTypeText, http.StatusForbidden, err.Error())
		return
	}

	respOut.Header.Add("X-Subscription", fmt.Sprint(ctx.Session))
	tee := &TeeReaderCloser{respOut.Body, out}
	respOut.Body = tee
	return
}

func (h *Handler) HandleConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	logger.Info.Printf("[%03d] %s", ctx.Session, host)
	return goproxy.MitmConnect, host
}
