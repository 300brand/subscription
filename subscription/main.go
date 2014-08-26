package main

import (
	"flag"
	"github.com/300brand/logger"
	"github.com/300brand/subscription"
	"net/http"
)

var addr = flag.String("listen", "0.0.0.0:8000", "Listen address")

func main() {
	flag.Parse()
	server := subscription.New()
	logger.Error.Fatal(http.ListenAndServe(*addr, server))
}
