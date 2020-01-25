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

// Gateway struct
type Gateway struct {
	proxy  service.ProxyService
	Pm     map[string]http.HandlerFunc
	notify <-chan *domain.MicroService
	done   chan interface{}
	mu     sync.Mutex
	new    chan<- string
}

// New returns a new gateway handler for a given service
func New(s service.ProxyService, new chan<- string) (*Gateway, error) {
	g := &Gateway{
		proxy: s,
		Pm:    make(map[string]http.HandlerFunc, 0),
		done:  make(chan interface{}),
		new:   new,
	}

	var err error
	g.notify, err = g.proxy.GetNotificationChan()
	if err != nil {
		return nil, err
	}

	g.listenNewService()

	return g, nil
}

// GetHandlers to get all the registered handlers
func (g *Gateway) loadAllProxies() error {
	services := g.proxy.GetServices()
	if len(services) == 0 {
		return errNoRegisteredServices
	}

	return g.loadProxies(services)
}

// LoadProxies from th db
func (g *Gateway) LoadProxies() error {
	return g.loadAllProxies()
}

// Register handlerFunc to register a new service
func (g *Gateway) Register(w http.ResponseWriter, r *http.Request) {
	var req client.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		writeErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := g.proxy.AddService(&req)
	if err != nil {
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeValidResponse(w, resp, http.StatusCreated)
}

func (g *Gateway) addProxy(s *domain.MicroService) error {
	p, err := newProxy(s)
	if err != nil {
		return err
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.Pm[p.service.Path]; ok {
		return nil
	}

	g.Pm[p.service.Path] = p.HandlerFunc
	return nil
}

func (g *Gateway) loadProxies(srvs map[string]*domain.MicroService) error {
	for _, s := range srvs {
		if err := g.addProxy(s); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gateway) listenNewService() {
	go func() {
		for {
			select {
			case s := <-g.notify:
				g.addProxy(s)
				// sends s.Path - service identifier
				// so that router can use this info to register new router
				g.new <- s.Path

			case <-g.done:
				break
			}
		}
	}()
}
