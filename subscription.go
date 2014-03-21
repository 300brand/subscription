package main

import (
	"flag"
	"github.com/300brand/subscription/config"
	"github.com/300brand/subscription/samplesite"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
)

var (
	UseSampleSite = flag.Bool("sample", false, "Use the internal sample site")
	Addr          = flag.String("addr", ":8082", "Listen address")
	ConfigFile    = flag.String("config", "subscription.json", "Subscription config")
)

func main() {
	flag.Parse()

	if err := config.Load(*ConfigFile); err != nil {
		log.Fatal("Could not load config:", err)
	}

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

	h := handlers.CombinedLoggingHandler(os.Stdout, http.DefaultServeMux)
	log.Fatal(http.ListenAndServe(*Addr, h))
}
