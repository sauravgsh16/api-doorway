package service

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/sauravgsh16/api-doorway/client"
	"github.com/sauravgsh16/api-doorway/domain"
	"github.com/sauravgsh16/api-doorway/store"
)

var (
	errNilChan = errors.New("channel is nil")
)

// ProxyService interface
type ProxyService interface {
	AddService(req *client.RegisterRequest) (*client.RegisterResponse, error)
	LoadServices() error
	GetServices() map[string]*domain.MicroService
	GetNotificationChan() (<-chan *domain.MicroService, error)
}

type service struct {
	store  store.MicroServiceStore
	msMap  map[string]*domain.MicroService
	notify chan *domain.MicroService
	mux    sync.RWMutex
}

// NewService returns a new service
func NewService(s store.MicroServiceStore) ProxyService {
	srv := &service{
		store:  s,
		msMap:  make(map[string]*domain.MicroService, 0),
		notify: make(chan *domain.MicroService),
	}

	if err := srv.LoadServices(); err != nil {
		log.Fatalf(err.Error())
	}
	return srv
}

// AddService is called when a new service requests it to be added.
// Along with the db, we also add the service to the in-memory storage
func (s *service) AddService(req *client.RegisterRequest) (*client.RegisterResponse, error) {
	var eps []*domain.Endpoint
	for _, ep := range req.Endpoints {
		dEp := &domain.Endpoint{
			Method: ep.Method,
			Path:   ep.Path,
		}
		eps = append(eps, dEp)
	}

	serv, err := s.store.AddService(req.Name, req.Host, req.Description, req.Path, eps)
	if err != nil {
		return nil, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()
	s.msMap[serv.ID] = serv

	// Notify new service was added to the handler
	// for it to register it to it's proxy map
	go func() {
		select {
		case s.notify <- serv:
		}
	}()

	return client.NewRegisterResponse(serv.ID, serv.Name), nil
}

// LoadService is called once when upon initial start-up of the application.
func (s *service) LoadServices() error {
	var err error
	var once sync.Once

	once.Do(func() {
		services, err := s.store.GetServices()
		if err != nil {
			err = fmt.Errorf("failed to load services: %s", err.Error())
		}

		s.mux.Lock()
		defer s.mux.Unlock()

		for _, srv := range services {
			if _, ok := s.msMap[srv.ID]; ok {
				// TODO : Add logic to check if service info have been refreshed.
				// If refreshed, change service map
				continue
			}
			// Add service to mapl
			s.msMap[srv.ID] = &srv
		}
	})

	return err
}

func (s *service) GetServices() map[string]*domain.MicroService {
	return s.msMap
}

func (s *service) GetNotificationChan() (<-chan *domain.MicroService, error) {
	if s.notify == nil {
		return nil, errNilChan
	}

	return s.notify, nil
}
