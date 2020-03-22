package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/sauravgsh16/api-doorway/domain"
)

// Proxy struct
type Proxy struct {
	proxy   *httputil.ReverseProxy
	service *domain.MicroService
}

func newProxy(s *domain.MicroService) (*Proxy, error) {
	var err error

	p := &Proxy{service: s}
	p.proxy, err = singleHostReverseProxy(s)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func singleHostReverseProxy(s *domain.MicroService) (*httputil.ReverseProxy, error) {
	host, err := url.Parse(s.Host)
	if err != nil {
		return nil, err
	}

	director := func(req *http.Request) {
		target := host
		targetQuery := host.RawQuery

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		// We are expected to receive incoming request in the form:
		// http://ip:port/service_identifier/actual_path
		// We need to remove the service_identifier and join the req.URL.Path
		// with the actual_path requested.
		// TODO: Better way to accomplish this
		req.URL.Path = fmt.Sprintf("/%s", strings.Join(strings.Split(req.URL.Path, "/")[2:], "/"))

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	return &httputil.ReverseProxy{Director: director}, nil
}

// HandlerFunc serves the incoming request
func (p *Proxy) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Proxy", "GatewayProxy")
	// TODO: At this point we want to add middlewares
	p.proxy.ServeHTTP(w, r)
}
