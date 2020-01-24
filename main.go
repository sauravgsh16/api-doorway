package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func main() {
	host := "http://localhost:9000/auth/foo/bar"

	url, err := url.Parse(host)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%s\n", url.Host)
	fmt.Printf("%v\n", strings.Join(strings.Split(url.Path, "/")[2:], "/"))

	// httputil.NewSingleHostReverseProxy(url)
}
