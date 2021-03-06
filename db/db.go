package db

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	// postgres driver
	_ "github.com/lib/pq"

	"github.com/sauravgsh16/api-doorway/domain"
)

func init() {
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}
}

// NewDB returns a pointer to gorm.DB
func NewDB() (*gorm.DB, error) {
	dbhost := os.Getenv("DB_HOST")
	dbport, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}
	dbuser := os.Getenv("DB_USER")
	dbpwd := os.Getenv("DB_PWD")
	dbname := os.Getenv("DB_NAME")
	dbtype := os.Getenv("DB_TYPE")

	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpwd, dbname,
	)

	db, err := gorm.Open(dbtype, connInfo)
	if err != nil {
		return nil, err
	}

	// db.DropTableIfExists(&domain.Endpoint{}, &domain.MicroService{})
	db.AutoMigrate(&domain.Endpoint{}, &domain.MicroService{})
	db.Model(&domain.Endpoint{}).AddForeignKey("service_id", "micro_services(id)", "CASCADE", "CASCADE")

	// TODO: set max idle and open connections

	db.LogMode(true)
	return db, nil
}
