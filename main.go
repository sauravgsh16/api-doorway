package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sauravgsh16/api-doorway/domain"
	"github.com/sauravgsh16/api-doorway/store"

	"github.com/sauravgsh16/api-doorway/db"
)

type Contact struct {
	ContactID   int `gorm:"primary_key"`
	CountryCode int
	MobileNo    uint
	CustId      int
}

type Customer struct {
	CustomerID   int `gorm:"primary_key"`
	CustomerName string
	Contacts     []Contact `gorm:"ForeignKey:CustId"` //you need to do like this
}

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

func main() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PWD", "postgres")
	os.Setenv("DB_NAME", "gateway")
	os.Setenv("DB_TYPE", "postgres")

	db, err := db.NewDB()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.Close()

	s := store.NewMicroServiceStore(db)
	eps := []*domain.Endpoint{
		&domain.Endpoint{Method: "get", Path: "/foo"},
		&domain.Endpoint{Method: "post", Path: "/bar"},
	}

	_, err = s.AddService("name", "host", "desc", "path", eps)
	if err != nil {
		log.Fatalf(err.Error())
	}

	service, err := s.GetServices()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%+v\n", service)

}
