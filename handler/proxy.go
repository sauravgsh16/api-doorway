package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/sauravgsh16/api-doorway/domain"
)

// Proxy struct
type Proxy struct {
	proxy   httputil.ReverseProxy
	service *domain.MicroService
}

func newProxy(s *domain.MicroService) (*Proxy, error) {
	p := &Proxy{service: s}

	handler, err := p.getReverseProxy()
	if err != nil {
		return nil, err
	}
	p.proxy = *handler
	return p, nil
}

func (p *Proxy) getReverseProxy() (*httputil.ReverseProxy, error) {
	host, err := url.Parse(p.service.Host)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(host), nil
}

// HandlerFunc serves the incoming request
func (p *Proxy) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Proxy", "GatewayProxy")
	// TODO: At this point we want to add middlewares
	p.proxy.ServeHTTP(w, r)
}
