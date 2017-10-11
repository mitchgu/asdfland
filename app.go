package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB DB
}

func (a *App) InitRedis(redis_addr, redis_pass string, redis_dbnum int) {
	rdb := RedisDB{}
	rdb.Init(redis_addr, redis_pass, redis_dbnum)
	a.DB = &rdb

	routes := a.GetRoutes()
	a.Router = mux.NewRouter().StrictSlash(true)
	for _, route := range *routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		a.Router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
}

func (a *App) Run(addr string) {
    srv := &http.Server{
        Handler:      a.Router,
        Addr:         addr,
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    log.Println("Starting server")
    log.Fatal(srv.ListenAndServe())
}