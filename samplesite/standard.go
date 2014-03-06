package samplesite

import (
	"fmt"
	"net/http"
	"path/filepath"
)

type standardLogin struct {
	CookieName string
	Title      string
}

var _ loginHandler = new(standardLogin)

func init() {
	register("/standard/", newStandardLogin(&tmpls))
}

func newStandardLogin(s *string) *standardLogin {
	*s += `{{ define "std_login.form" }}{{ template "header" . }}
	<form method="post" action="login.do">
		<div><input type="text" name="username" placeholder="Username"></div>
		<div><input type="password" name="password" placeholder="Password"></div>
		<div><button type="submit">Login</button></div>
	</form>{{ template "footer" . }}{{ end }}`
	return &standardLogin{
		CookieName: "standardLogin",
		Title:      "Standard Login",
	}
}

func (s *standardLogin) LoggedIn(r *http.Request) bool {
	c, err := r.Cookie(s.CookieName)
	return err == nil && c.Value == "true"
}

func (s *standardLogin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch f := filepath.Base(r.RequestURI); f {
	case "login.form":
		base.ExecuteTemplate(w, "std_login.form", s)
	case "login.do":
		u, p := r.PostFormValue("username"), r.PostFormValue("password")
		valid := u == Username && p == Password
		http.SetCookie(w, &http.Cookie{
			Name:  s.CookieName,
			Value: fmt.Sprint(valid),
			Path:  "/",
		})
		if valid {
			http.Redirect(w, r, "success", http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(w, r, "failed", http.StatusTemporaryRedirect)
	case "failed":
		http.Error(w, "failed", http.StatusUnauthorized)
	case "success":
		fmt.Fprint(w, "success")
	default:
		http.Error(w, "Not found: "+f, http.StatusNotFound)
	}
}
