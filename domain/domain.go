package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
)

var (
	errInvalidServiceName     = errors.New("invalid service name")
	errInvalidBaseURL         = errors.New("invalid base URL")
	errInvalidEndpoint        = errors.New("invalid endpoint url")
	errNoMatchBaseAndEndpoint = errors.New("no match - base URL and endpoint URL")
)

// EndPoint struct
type EndPoint struct {
	Method string `json:"method" sql:"type:varchar(6);not null"`
	Path   string `json:"path" sql:"type:varchar(50)"`
}

// MicroService struct
type MicroService struct {
	ID          string      `json:"id" gorm:"primay_key"`
	Name        string      `json:"name" sql:"varchar(50);index;unique;not null"`
	Endpoints   []*EndPoint `json:"end_points" sql:"varchar(250);not null" gorm:"foreignkey:Path"`
	Host        string      `json:"host" sql:"type:varchar(250);unique;not null"`
	Description string      `json:"description" sql:"type:varchar(250)"`
	Running     bool        `json:"running" sql:"type:boolean"`

	// TODO: Support for multiple instances running on different ports
	// TODO: Add slice - containing list of running instances
}

// NewMicroService returns a pointer to a new microservice
func NewMicroService(name, host, desc string, eps []*EndPoint) *MicroService {
	return &MicroService{
		ID:          fmt.Sprintf("%s", uuid.Must(uuid.NewV4())),
		Host:        host,
		Name:        name,
		Description: desc,
		Endpoints:   eps,
	}
}

// Validate the service struct
func (ms *MicroService) Validate() error {
	if len(strings.TrimSpace(ms.Name)) == 0 {
		return errInvalidServiceName
	}
	// TODO : Need to update validation

	/*
		b, err := url.Parse(ms.BaseURL)
		if err != nil {
			return ErrInvalidBaseURL
		}

		for _, epURL := range ms.Endpoints {
			ep, err := url.Parse(epURL)
			if err != nil {
				return ErrInvalidEndpoint
			}
			if b.Host != ep.Host {
				return ErrNoMatchBaseAndEndpoint
			}
		}
	*/
	return nil
}

// UpdateStatus updates the status of a microservice
func (ms *MicroService) UpdateStatus(status bool) {
	ms.Running = status
}

// Status returns the status of the servicel
func (ms *MicroService) Status() bool {
	return ms.Running
}
