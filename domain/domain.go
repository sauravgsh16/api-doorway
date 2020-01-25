package domain

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

var (
	errInvalidServiceName     = errors.New("invalid service name")
	errInvalidBaseURL         = errors.New("invalid base URL")
	errInvalidEndpoint        = errors.New("invalid endpoint url")
	errNoMatchBaseAndEndpoint = errors.New("no match - base URL and endpoint URL")
)

// Endpoint struct
type Endpoint struct {
	ID        int    `json:"id" gorm:"primary_key"`
	Method    string `json:"method" gorm:"type:varchar(6);not null"`
	Path      string `json:"path" gorm:"type:varchar(50)"`
	ServiceID string `json:"service_id"`
}

// MicroService struct
type MicroService struct {
	ID          string     `json:"id" gorm:"primay_key"`
	Name        string     `json:"name" gorm:"varchar(50);index;unique;not null"`
	Path        string     `json:"path" gorm:"varchar(10);unique;no null"`
	Host        string     `json:"host" gorm:"type:varchar(250);unique;not null"`
	Description string     `json:"description" gorm:"type:varchar(250)"`
	Running     bool       `json:"running" gorm:"type:boolean"`
	Endpoints   []Endpoint `json:"end_points" gorm:"ForeignKey:ServiceID"`

	// TODO: Support for multiple instances running on different ports
	// TODO: Add slice - containing list of running instances
}

// NewMicroService returns a pointer to a new microservice
func NewMicroService(name, host, desc, path string, eps []Endpoint) *MicroService {
	return &MicroService{
		ID:          fmt.Sprintf("%s", uuid.Must(uuid.NewV4())),
		Host:        host,
		Name:        name,
		Path:        path,
		Description: desc,
		Endpoints:   eps,
	}
}

// UpdateStatus updates the status of a microservice
func (ms *MicroService) UpdateStatus(status bool) {
	ms.Running = status
}

// Status returns the status of the servicel
func (ms *MicroService) Status() bool {
	return ms.Running
}
