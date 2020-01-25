package store

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/sauravgsh16/api-doorway/domain"
)

var (
	errServiceNotFound     = errors.New("service not found")
	errServiceAlreadyExist = errors.New("service already exists")
	errFailedToGetServices = errors.New("failed to get services from db")
)

// MicroServiceStore interface
type MicroServiceStore interface {
	FindServiceByName(string) (*domain.MicroService, error)
	AddService(name, host, desc, path string, eps []*domain.Endpoint) (*domain.MicroService, error)
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

	notFound := mss.db.Where("name = ?", name).Preload("Endpoints").First(&ms).RecordNotFound()
	if notFound {
		return nil, errServiceNotFound
	}

	return &ms, nil
}

func (mss *microserviceStore) AddService(name, host, desc, path string, eps []*domain.Endpoint) (*domain.MicroService, error) {
	_, err := mss.FindServiceByName(name)
	if err == nil {
		return nil, errServiceAlreadyExist
	}

	ms := domain.NewMicroService(name, host, desc, path, eps)

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
		return nil, errFailedToGetServices
	}
	return services, nil
}
