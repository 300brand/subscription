package samplesite

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
)

type loggedinFunc func(r *http.Request) bool

type loginHandler interface {
	LoggedIn(r *http.Request) bool
	http.Handler
}

var (
	Username = "username"
	Password = "password"

	base   *template.Template
	logins = make(map[string]loginHandler)
	server *httptest.Server
)

var tmpls = `
{{ define "header" }}<!DOCTYPE html>
<html>
<head>
	<title>{{ .Title }}</title>
</head>
<body>
	<header>
		<h1>{{ .Title }}</h1>
	</header>
{{ end }}

{{ define "footer" }}
</body>
</html>
{{ end }}`

func Close() {
	if server != nil {
		server.Close()
	}
}

func Start() (url string) {
	var err error

	mux := http.NewServeMux()
	for path, login := range logins {
		mux.Handle(path, login)
	}
	mux.HandleFunc("/loggedin.check", func(w http.ResponseWriter, r *http.Request) {
		for _, login := range logins {
			if login.LoggedIn(r) {
				fmt.Fprint(w, 1)
				return
			}
		}
		fmt.Fprint(w, 0)
	})

	server = httptest.NewServer(mux)

	base, err = template.New("base").Parse(tmpls)
	if err != nil {
		panic(err)
	}

	return server.URL
}

func register(path string, lh loginHandler) {
	logins[path] = lh
}
