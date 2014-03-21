package main

import (
	"github.com/300brand/logger"
	"github.com/300brand/subscription/authorize"
	"github.com/300brand/subscription/config"
	"github.com/300brand/subscription/relink"
	"io"
	"net/http"
	"strings"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	alias := r.Host[:strings.IndexByte(r.Host, '.')]
	logger.Debug.Printf("Alias: %s", alias)
	domain, err := config.Get(alias)
	if err != nil {
		http.Error(w, "Invalid alias: "+alias, http.StatusBadRequest)
		return
	}
	logger.Debug.Printf("Domain: %s", domain.Domain)

	auth, err := authorize.Get(domain.LoginType)
	if err != nil {
		http.Error(w, "Invalid authorizer: "+domain.LoginType, http.StatusBadRequest)
		return
	}
	logger.Debug.Printf("Authorizer: %s", domain.LoginType)

	switch loggedIn, err := auth.LoggedIn(domain); true {
	case err != nil:
		logger.Error.Printf("Error checking logged in state: %s", err)
		http.Error(w, "Error checking logged in state: "+err.Error(), http.StatusInternalServerError)
		return
	case !loggedIn:
		logger.Info.Printf("Logging in %s:%s@%s", domain.Username[1], domain.Password[1], domain.URL().Host)
		if err := auth.Login(domain); err != nil {
			logger.Error.Printf("Error logging in: %s", err)
			http.Error(w, "Error logging in: "+err.Error(), http.StatusUnauthorized)
			return
		}
	}

	remoteURL := domain.ResolveReference(r.RequestURI)

	// Create request for remote site
	req, err := http.NewRequest(r.Method, remoteURL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Body != nil {
		defer r.Body.Close()
	}
	req.Header = r.Header
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1622.0 Safari/537.36")

	resp, err := domain.Client().Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Set up relinker
	domainMap := make(map[string]string, len(domain.Rewrite))
	for _, rw := range domain.Rewrite {
		domainMap[rw] = r.Host
	}
	logger.Debug.Printf("domainMap: %+v", domainMap)
	relinker := relink.New(resp.Body, domainMap)

	for key, values := range resp.Header {
		for _, value := range values {
			logger.Trace.Printf("Adding header %s = %s", key, value)
			w.Header().Add(key, value)
		}
	}
	w.Header().Add("X-Remote-URL", remoteURL.String())
	w.Header().Del("Content-Length")
	w.Header().Del("Set-Cookie")
	w.WriteHeader(resp.StatusCode)

	n, _ := io.Copy(w, relinker)
	logger.Info.Printf("%s -> %s %d %d", r.URL, remoteURL, resp.StatusCode, n)
}
