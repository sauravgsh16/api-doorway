package route

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/sauravgsh16/api-doorway/db"
	"github.com/sauravgsh16/api-doorway/proxy"
	"github.com/sauravgsh16/api-doorway/service"
	"github.com/sauravgsh16/api-doorway/store"
)

var (
	r = mux.NewRouter()
)

func initRoute() {
	db, err := db.NewDB()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize db: %s", err))
	}

	store := store.NewMicroServiceStore(db)
	service := service.NewService(store)
	proxy := proxy.NewProxyHandler(service)

	_, err = proxy.GetProxyHandlers()
	if err != nil {
		log.Println(fmt.Errorf("failed to load proxies: %s", err))
	}
	/*
		for path, h := range handlers {
			reflect.
		}
	*/
	r.HandleFunc("/register", proxy.Register).Methods("POST")
}

func Run() {
	initRoute()
	srv := http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()
}
