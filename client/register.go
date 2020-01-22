package client

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/gofrs/uuid"
)

var (
	errInvalidHost    = errors.New("host name invalid - empty")
	errEmptyEndpoints = errors.New("need endpoints for to register service")
)

// RegisterRequest struct
type RegisterRequest struct {
	Name        string   `json:"name"`
	Host        string   `json:"host"`
	Endpoints   []string `json:"end_points"`
	Description string   `json:"description"`
}

// Validate the request structure
func (r *RegisterRequest) Validate() error {
	if r.Name == "" {
		r.Name = fmt.Sprintf("%s", uuid.Must(uuid.NewV4()))
	}

	if r.Host == "" {
		return errInvalidHost
	}

	if _, err := url.Parse(r.Host); err != nil {
		return fmt.Errorf("host url invalid, %s", err.Error())
	}

	if len(r.Endpoints) == 0 {
		return errEmptyEndpoints
	}

	return nil
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
