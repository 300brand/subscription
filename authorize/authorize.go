package authorize

import (
	"errors"
	"github.com/300brand/subscription/config"
)

type Authorizer interface {
	Login(cfg *config.Domain) (err error)
	LoggedIn(cfg *config.Domain) (loggedIn bool, err error)
}

var (
	ErrInvalidRedirect = errors.New("Invalid redirect after login")
)
