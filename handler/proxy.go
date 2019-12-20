package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/sauravgsh16/api-doorway/domain"
)

type Proxy struct {
	proxy   httputil.ReverseProxy
	service *domain.MicroService
}

func newProxy(s *domain.MicroService) (*Proxy, error) {
	ph := &Proxy{service: s}

	handler, err := ph.getReverseProxy()
	if err != nil {
		return nil, err
	}
	ph.proxy = *handler
	return ph, nil
}

func (ph *Proxy) getReverseProxy() (*httputil.ReverseProxy, error) {
	host, err := url.Parse(ph.service.Host)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(host), nil
}

func (ph *Proxy) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Proxy", "GatewayProxy")
	// TODO: At this point we want to add middlewares
	ph.proxy.ServeHTTP(w, r)
}
