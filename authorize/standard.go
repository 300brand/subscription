package authorize

import (
	"github.com/300brand/subscription/config"
	"log"
	"net/http"
	"net/url"
)

type Standard struct{}

var _ Authorizer = new(Standard)

func (a *Standard) Login(cfg *config.Domain) (err error) {
	var baselineResp, loginResp *http.Response
	// (Re)-Establish baseline cookies
	urlForm := cfg.ResolveReference(cfg.URLs.Form)
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

	// Perform login
	urlDo := cfg.ResolveReference(cfg.URLs.Do)
	if loginResp, err = cfg.Client().PostForm(urlDo.String(), data); err != nil {
		return
	}

	// Check the redirect
	loc, err := loginResp.Location()
	if err != nil {
		return
	}

	if loc.Path != cfg.URLs.Success {
		return ErrInvalidRedirect
	}

	return
}

func (a *Standard) LoggedIn(cfg *config.Domain) (loggedIn bool, err error) {
	name, value := cfg.Cookie[0], cfg.Cookie[1]
	for i, c := range cfg.Client().Jar.Cookies(cfg.URL()) {
		log.Printf("[%d] %s", i, c)
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
