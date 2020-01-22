package main

import (
	"log"
	"net/http/httputil"
	"net/url"
)

func main() {
	host := "http://localhost:9000/test"

	url, err := url.Parse(host)
	if err != nil {
		log.Fatalf(err.Error())
	}

	httputil.NewSingleHostReverseProxy(url)
}
