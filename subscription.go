package main

import (
	"flag"
	"github.com/300brand/subscription/samplesite"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var (
	Sites         = make(map[string]SiteConfig)
	UseSampleSite = flag.Bool("sample", false, "Use the internal sample site")
	Addr          = flag.String("addr", ":8082", "Listen address")
)

func handler(w http.ResponseWriter, r *http.Request) {
	i := strings.IndexByte(r.Host, '.')
	alias := r.Host[:i]
	site, ok := Sites[alias]
	if !ok {
		http.Error(w, "Invalid alias: "+alias, http.StatusBadRequest)
		return
	}
	ref := new(url.URL)
	*ref = *r.URL
	ref.Host = ""
	ref.Path = ref.Path
	remoteURL := site.url.ResolveReference(ref)

	req, err := http.NewRequest(r.Method, remoteURL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Body != nil {
		defer r.Body.Close()
	}
	req.Header = r.Header

	for i, c := range site.client.Jar.Cookies(remoteURL) {
		log.Printf("[%d] Expires: %s", i, c)
	}

	resp, err := site.client.Do(req)
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
	w.WriteHeader(resp.StatusCode)

	n, _ := io.Copy(w, resp.Body)
	log.Printf("%s -> %s %d %d", r.URL, remoteURL, resp.StatusCode, n)
}

func main() {
	flag.Parse()

	if *UseSampleSite {
		Sites["sample"] = SiteConfig{
			Alias:        "sample_std",
			Siteroot:     samplesite.Start(),
			LoginType:    "standard",
			Username:     samplesite.Username,
			Password:     samplesite.Password,
			LoginForm:    "/standard/login.form",
			LoginDo:      "/standard/login.do",
			LoginSuccess: "/standard/success",
			LoginFailure: "/standard/failure",
		}
	}

	for sub, site := range Sites {
		var err error
		if site.url, err = url.Parse(site.Siteroot); err != nil {
			log.Fatalf("Error parsing %s: %s", site.Siteroot, err)
		}
		site.client = new(http.Client)
		site.client.Jar, _ = cookiejar.New(nil)
		log.Printf("Registered: %s -> %s", sub, site.Siteroot)
		Sites[sub] = site
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*Addr, nil))
}
