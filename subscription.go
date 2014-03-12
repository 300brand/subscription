package main

import (
	"flag"
	"github.com/300brand/subscription/authorize"
	"github.com/300brand/subscription/config"
	"github.com/300brand/subscription/samplesite"
	"github.com/gorilla/handlers"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	UseSampleSite = flag.Bool("sample", false, "Use the internal sample site")
	Addr          = flag.String("addr", ":8082", "Listen address")
	ConfigFile    = flag.String("config", "/tmp/subscription.json", "Subscription config")
)

func handler(w http.ResponseWriter, r *http.Request) {
	alias := r.Host[:strings.IndexByte(r.Host, '.')]
	domain, err := config.Get(alias)
	if err != nil {
		http.Error(w, "Invalid alias: "+alias, http.StatusBadRequest)
		return
	}

	auth, err := authorize.Get(domain.LoginType)
	if err != nil {
		http.Error(w, "Invalid authorizer: "+domain.LoginType, http.StatusBadRequest)
		return
	}

	switch loggedIn, err := auth.LoggedIn(domain); true {
	case err != nil:
		http.Error(w, "Error checking logged in state: "+err.Error(), http.StatusInternalServerError)
		return
	case !loggedIn:
		if err := auth.Login(domain); err != nil {
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

	resp, err := domain.Client().Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			log.Printf("Adding header %s = %s", key, value)
			w.Header().Add(key, value)
		}
	}
	w.Header().Add("X-Remote-URL", remoteURL.String())
	w.WriteHeader(resp.StatusCode)

	n, _ := io.Copy(w, resp.Body)
	log.Printf("%s -> %s %d %d", r.URL, remoteURL, resp.StatusCode, n)
}

func main() {
	flag.Parse()

	if *UseSampleSite {
		config.Add(&config.Domain{
			Alias:     "sample",
			Domain:    samplesite.Start(),
			LoginType: "standard",
			Username:  [2]string{"username", samplesite.Username},
			Password:  [2]string{"password", samplesite.Password},
			URLs: config.URLs{
				Form:    "/standard/login.form",
				Do:      "/standard/login.do",
				Success: "/standard/success",
			},
		})
	}

	http.HandleFunc("/", handler)
	h := handlers.CombinedLoggingHandler(os.Stdout, http.DefaultServeMux)
	log.Fatal(http.ListenAndServe(*Addr, h))
}
