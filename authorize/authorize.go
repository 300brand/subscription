package authorize

import (
	"github.com/300brand/subscription/config"
)

type Authorizer interface {
	Login(cfg *config.Config) (err error)
	LoggedIn(cfg *config.Config) (loggedIn bool, err error)
}
