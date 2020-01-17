package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sauravgsh16/api-doorway/client"
	"github.com/sauravgsh16/api-doorway/service"
)

var (
	errNoRegisteredServices = errors.New("no registered services")
)

type GateWayHandler interface{}

type gateway struct {
	proxy service.ProxyService
}

func NewGateWayHandler(s service.ProxyService) GateWayHandler {
	return &gateway{proxy: s}
}

func (g *gateway) Register(w http.ResponseWriter, r *http.Request) {
	var req client.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Validate req structure

	resp, err := g.proxy.AddService(&req)
	if err != nil {
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeValidResponse(w, resp, http.StatusCreated)
}

func (g *gateway) GetProxyHandlers() (map[string]http.HandlerFunc, error) {
	services := g.proxy.GetServices()
	if len(services) == 0 {
		return nil, errNoRegisteredServices
	}

	proxyMap := make(map[string]http.HandlerFunc)

	// TODO: Need to modify below map.
	//

	for _, srv := range services {
		proxy, err := newProxy(srv)
		if err != nil {
			return nil, err
		}
		proxyMap[proxy.service.Path] = proxy.HandlerFunc
	}

	return proxyMap, nil
}

/*
func getHandlers(srv *domain.MicroService) func(http.ResponseWriter, *http.Request) {

}
*/
