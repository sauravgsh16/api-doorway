package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/sauravgsh16/api-doorway/client"
	"github.com/sauravgsh16/api-doorway/domain"
	"github.com/sauravgsh16/api-doorway/service"
)

var (
	errNoRegisteredServices = errors.New("no registered services")
)

type GateWayHandler interface{}

type gateway struct {
	proxy  service.ProxyService
	pm     map[string]http.HandlerFunc
	notify <-chan *domain.MicroService
	done   chan interface{}
	mu     sync.Mutex
}

// New returns a new gateway handler for a given service
func New(s service.ProxyService) (GateWayHandler, error) {
	g := &gateway{
		proxy: s,
		pm:    make(map[string]http.HandlerFunc, 0),
		done:  make(chan interface{}),
	}

	var err error
	g.notify, err = g.proxy.GetNotificationChan()
	if err != nil {
		return nil, err
	}

	g.listenNewService()

	return g, nil
}

func (g *gateway) Register(w http.ResponseWriter, r *http.Request) {
	var req client.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Validate req structure

	if err := req.Validate(); err != nil {
		writeErrResponse(w, err.Error(), http.StatusBadRequest)
	}

	resp, err := g.proxy.AddService(&req)
	if err != nil {
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeValidResponse(w, resp, http.StatusCreated)
}

func (g *gateway) addProxy(s *domain.MicroService) error {
	p, err := newProxy(s)
	if err != nil {
		return err
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	// TODO: Find proper Key for storing in map
	if _, ok := g.pm[p.service.Host]; ok {
		return nil
	}

	g.pm[p.service.Host] = p.HandlerFunc
	return nil
}

func (g *gateway) loadProxies(srvs map[string]*domain.MicroService) error {
	for _, s := range srvs {
		if err := g.addProxy(s); err != nil {
			return err
		}
	}
	return nil
}

func (g *gateway) GetProxyHandlers() error {
	services := g.proxy.GetServices()
	if len(services) == 0 {
		return errNoRegisteredServices
	}

	return g.loadProxies(services)
}

func (g *gateway) listenNewService() {
	go func() {
		for {
			select {
			case s := <-g.notify:
				go g.addProxy(s)
			case <-g.done:
				break
			}
		}
	}()
}

/*
func getHandlers(srv *domain.MicroService) func(http.ResponseWriter, *http.Request) {

}
*/
