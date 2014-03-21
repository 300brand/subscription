package authorize

import (
	"github.com/300brand/logger"
	"github.com/300brand/subscription/config"
	"net/http"
	"net/url"
)

type Standard struct{}

var _ Authorizer = new(Standard)

func init() {
	Register("standard", new(Standard))
}

func (a *Standard) Login(cfg *config.Domain) (err error) {
	var baselineResp, loginResp *http.Response
	// (Re)-Establish baseline cookies
	urlForm := cfg.ResolveReference(cfg.URLs.Form)
	logger.Trace.Printf("Standard.Login: Baseline: %s", urlForm)
	if baselineResp, err = cfg.Client().Get(urlForm.String()); err != nil {
		return
	}
	baselineResp.Body.Close()

	// Prepare POST data
	data := url.Values{
		cfg.Username[0]: []string{cfg.Username[1]},
		cfg.Password[0]: []string{cfg.Password[1]},
	}
	for k, v := range cfg.Extra {
		data[k] = v
	}
	logger.Trace.Printf("Standard.Login: Sending %s", data.Encode())

	// Perform login
	urlDo := cfg.ResolveReference(cfg.URLs.Do)
	logger.Trace.Printf("Standard.Login: Doing: %s", urlDo)
	if loginResp, err = cfg.Client().PostForm(urlDo.String(), data); err != nil {
		return
	}
	defer loginResp.Body.Close()

	for k, v := range loginResp.Header {
		logger.Trace.Printf("Standard.Login: Header[%q] = %q", k, v)
	}

	for i, c := range loginResp.Cookies() {
		logger.Trace.Printf("Standard.Login: Cookie[%d] = %s", i, c)
	}

	logger.Trace.Printf("Standard.Login: Status %d - %s", loginResp.StatusCode, loginResp.Status)

	return

	// Check the redirect, if one is expected
	if cfg.URLs.Success != "" {
		loc, err := loginResp.Location()
		if err != nil {
			return err
		}

		if loc.Path != cfg.URLs.Success {
			return ErrInvalidRedirect
		}
	}

	// Verify we're logged in
	if in, err := a.LoggedIn(cfg); !in || err != nil {
		return ErrUnauthorized
	}

	return
}

func (a *Standard) LoggedIn(cfg *config.Domain) (loggedIn bool, err error) {
	name, value := cfg.Cookie[0], cfg.Cookie[1]
	for _, c := range cfg.Client().Jar.Cookies(cfg.URL()) {
		logger.Trace.Printf("%s = %s - Cookie: %s", name, value, c)
		loggedIn = false ||
			(name == c.Name && value == c.Value) ||
			(name == c.Name && value == "") ||
			(name == "" && value == c.Value)
		if loggedIn {
			break
		}
	}
	return
}
