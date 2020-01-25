package router

import (
	"fmt"
	"sync"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/sauravgsh16/api-doorway/handler"
	"github.com/sauravgsh16/api-doorway/service"
	"github.com/sauravgsh16/api-doorway/store"
)

type Router struct {
	R    *mux.Router
	g    *handler.Gateway
	done chan interface{}
	new  chan string
	mux  sync.Mutex
}

// New returns a new router
func New(db *gorm.DB) (*Router, error) {
	var err error

	r := &Router{
		R:    mux.NewRouter(),
		done: make(chan interface{}),
		new:  make(chan string),
	}

	s := store.NewMicroServiceStore(db)
	srv := service.NewService(s)
	r.g, err = handler.New(srv, r.new)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Router) routeUpdater() {
	go func() {
		for {
			select {
			case <-r.done:
				break
			case n := <-r.new:
				h := r.getHandler(n)
				r.addHandler(n, h)
			}
		}
	}()
}

func (r *Router) getHandler(path string) *handler.EndpointHandler {
	r.mux.Lock()
	defer r.mux.Unlock()

	p, ok := r.g.ProxyMap[path]
	if !ok {
		return nil
	}
	return p
}

func (r *Router) addHandler(path string, h *handler.EndpointHandler) {
	for _, ep := range h.Eps {
		path := "/" + path + fmt.Sprintf("/%s", ep.Path)
		r.R.HandleFunc(path, h.Hf.ServeHTTP).Methods(ep.Method)
	}
}

// Init router
func (r *Router) Init() error {
	// handle the register handler - for new service registration
	r.R.HandleFunc("/register", r.g.Register).Methods("POST")

	// load all proxies for services which are present in the db
	if err := r.g.LoadProxies(); err != nil {
		return err
	}

	// load all handlers for the proxies
	// TODO: check health of service before adding it to router
	for p, h := range r.g.ProxyMap {
		r.addHandler(p, h)
	}
	return nil
}
