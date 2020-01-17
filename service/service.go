package service

import (
	"fmt"
	"sync"

	"github.com/sauravgsh16/api-doorway/client"
	"github.com/sauravgsh16/api-doorway/domain"
	"github.com/sauravgsh16/api-doorway/store"
)

// ProxyService interface
type ProxyService interface {
	AddService(req *client.RegisterRequest) (*client.RegisterResponse, error)
	LoadServices() error
	GetServices() map[string]*domain.MicroService
}

type service struct {
	mux   sync.RWMutex
	msMap map[string]*domain.MicroService
	store store.MicroServiceStore
}

// NewService returns a new service
func NewService(s store.MicroServiceStore) ProxyService {
	return &service{
		msMap: make(map[string]*domain.MicroService, 0),
		store: s,
	}
}

// AddService is called when a new service requests it to be added.
// Along with the db, we also add the service to the in-memory storage
func (s *service) AddService(req *client.RegisterRequest) (*client.RegisterResponse, error) {
	serv, err := s.store.AddService(req.Name, req.Path, req.Host, req.Description, req.Endpoints)
	if err != nil {
		return nil, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()
	s.msMap[serv.ID] = serv

	return client.NewRegisterResponse(serv.ID, serv.Name), nil
}

// LoadService is called once when upon initial start-up of the application.
// TODO: make sure to wrap call with sync.Once
func (s *service) LoadServices() error {
	services, err := s.store.GetServices()
	if err != nil {
		return fmt.Errorf("failed to load services: %s", err.Error())
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
	return nil
}

func (s *service) GetServices() map[string]*domain.MicroService {
	return s.msMap
}

/*
serv register

{"auth0": MicroService{
	host: "http://auth0:8080",
	enpoints: ["register", "authenticate"],
	},
 "items": MicroService{
	host: "http://items:8081",
	endpoints: ["create", "order"]
	},
}

user ->

*/
