package authorize

import (
	"errors"
	"github.com/300brand/subscription/config"
)

type Authorizer interface {
	Login(cfg *config.Domain) (err error)
	LoggedIn(cfg *config.Domain) (loggedIn bool, err error)
}

var authorizers = make(map[string]Authorizer)

var (
	ErrInvalidAuthorizer = errors.New("Invalid authorizer name")
	ErrInvalidRedirect   = errors.New("Invalid redirect after login")
	ErrUnauthorized      = errors.New("Unauthorized")
)

func Register(name string, a Authorizer) {
	if _, ok := authorizers[name]; ok {
		panic("Authorizer " + name + " already exists")
	}
	authorizers[name] = a
}

func Get(name string) (a Authorizer, err error) {
	a, ok := authorizers[name]
	if !ok {
		return nil, ErrInvalidAuthorizer
	}
	return
}
