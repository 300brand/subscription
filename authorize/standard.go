package authorize

import (
	"github.com/300brand/subscription/config"
)

type Standard struct{}

var _ Authorizer = new(Standard)

func (a *Standard) Login(cfg *config.Config) (err error) {
	return
}

func (a *Standard) LoggedIn(cfg *config.Config) (loggedIn bool, err error) {
	return
}
