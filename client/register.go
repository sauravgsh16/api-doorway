package client

import (
	"errors"
	"net/url"
	"strings"
)

var (
	errInvalidServiceName = errors.New("invalid service name")
	errInvalidBaseURL     = errors.New("invalid base url")
	errInvalidMethod      = errors.New("invalid HTTP method")
	errInvalidEndpoints   = errors.New("invalid end points")
)

// Endpoint struct
type Endpoint struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// RegisterRequest struct
type RegisterRequest struct {
	Name        string     `json:"name"`
	Host        string     `json:"host"`
	Endpoints   []Endpoint `json:"end_points"`
	Description string     `json:"description"`
}

// Validate the request structure
func (req *RegisterRequest) Validate() error {
	if checkEmpty(req.Name) {
		return errInvalidServiceName
	}

	if checkEmpty(req.Host) {
		return errInvalidBaseURL
	}

	if checkEmpty(req.Endpoints) {
		return errInvalidEndpoints
	}

	_, err := url.ParseRequestURI(req.Host)
	if err != nil {
		return errInvalidBaseURL
	}

	for _, ep := range req.Endpoints {
		if len(strings.TrimSpace(ep.Path)) == 0 {
			ep.Method = "/"
		}
		if len(strings.TrimSpace(ep.Method)) == 0 {
			return errInvalidMethod
		}
	}

	return nil
}

func checkEmpty(v interface{}) bool {
	switch t := v.(type) {
	case string:
		if len(strings.TrimSpace(t)) == 0 {
			return true
		}
	case []Endpoint:
		if len(t) == 0 {
			return true
		}
	default:
		return true
	}
	return false
}

// RegisterResponse struct
type RegisterResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NewRegisterResponse returns a new registration response
func NewRegisterResponse(id, name string) *RegisterResponse {
	return &RegisterResponse{
		ID:   id,
		Name: name,
	}
}
