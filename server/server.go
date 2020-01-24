package server

import (
	"log"
	"net/http"

	"github.com/sauravgsh16/api-doorway/config"
	"github.com/sauravgsh16/api-doorway/db"
	"github.com/sauravgsh16/api-doorway/router"
)

// Run server
func Run() {
	db, err := db.NewDB()
	if err != nil {
		log.Fatalf(err.Error())
	}

	r, err := router.New(db)
	if err != nil {
		log.Fatalf(err.Error())
	}
	r.Init()

	server := &http.Server{
		Addr:    config.Addr,
		Handler: r.R,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf(err.Error())
	}
}
