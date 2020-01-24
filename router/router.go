package router

import (
	"sync"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/sauravgsh16/api-doorway/handler"
	"github.com/sauravgsh16/api-doorway/service"
	"github.com/sauravgsh16/api-doorway/store"
)

type Router struct {
	R   *mux.Router
	g   *handler.Gateway
	d   chan interface{}
	new chan string
	mux sync.Mutex
}

// New returns a new router
func New(db *gorm.DB) (*Router, error) {
	var err error

	r := &Router{
		R:   mux.NewRouter(),
		d:   make(chan interface{}),
		new: make(chan string),
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
			case <-r.d:
				break
			case n := <-r.new:
				r.addHandler(n)
			}
		}
	}()
}

func (r *Router) addHandler(pathName string) {
	r.mux.Lock()
	defer r.mux.Unlock()

	p, ok := r.g.Pm[pathName]
	if !ok {
		// TODO: better way to check why handler func is not present.
		return
	}
	path := "/" + pathName + "/*"
	r.R.HandleFunc(path, p.ServeHTTP)
}

// Init router
func (r *Router) Init() error {
	// handle the register handler - for new service registration
	r.R.HandleFunc("/register", r.g.Register)

	// load all proxies for services which are present in the db
	if err := r.g.LoadProxies(); err != nil {
		return err
	}

	// load all handlers for the proxies
	for p, h := range r.g.Pm {
		r.R.HandleFunc("/"+p+"/*", h)
	}
	return nil
}
