package db

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	// prosgres driver
	_ "github.com/lib/pq"

	"github.com/sauravgsh16/api-doorway/config"
)

func init() {
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}
}

// NewDB returns a pointer to gorm.DB
func NewDB() (*gorm.DB, error) {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPwd,
		config.DBName,
	)

	db, err := gorm.Open(config.DBType, connInfo)
	if err != nil {
		return nil, err
	}

	// TODO: set max idle and open connections

	db.LogMode(true)
	return db, nil
}
