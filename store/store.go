package store

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/sauravgsh16/api-doorway/domain"
)

var (
	ErrServiceNotFound     = errors.New("service not found")
	ErrServiceAlreadyExist = errors.New("service already exists")
	ErrFailedToGetServices = errors.New("failed to get services from db")
)

// MicroServiceStore interface
type MicroServiceStore interface {
	FindServiceByName(string) (*domain.MicroService, error)
	AddService(name, bURL, host, desc string, eps []string) (*domain.MicroService, error)
	GetServices() ([]domain.MicroService, error)
}

type microserviceStore struct {
	db *gorm.DB
}

// NewMicroServiceStore returns a new microservice store
func NewMicroServiceStore(db *gorm.DB) MicroServiceStore {
	return &microserviceStore{db: db}
}

func (mss *microserviceStore) FindServiceByName(name string) (*domain.MicroService, error) {
	var ms domain.MicroService

	notFound := mss.db.Where("name = ?", name).First(&ms).RecordNotFound()
	if notFound {
		return nil, ErrServiceNotFound
	}

	return &ms, nil
}

func (mss *microserviceStore) AddService(name, bURL, host, desc string, eps []string) (*domain.MicroService, error) {
	_, err := mss.FindServiceByName(name)
	if err != nil {
		return nil, ErrServiceAlreadyExist
	}

	ms := domain.NewMicroService(name, bURL, host, desc, eps)
	if err := ms.Validate(); err != nil {
		return nil, err
	}

	tx := mss.db.Begin()

	if err := mss.db.Create(ms).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return ms, nil
}

func (mss *microserviceStore) GetServices() ([]domain.MicroService, error) {
	var services []domain.MicroService

	if err := mss.db.Find(&services).Error; err != nil {
		return nil, ErrFailedToGetServices
	}
	return services, nil
}
