package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     DB
}

func (a *App) InitRedis(redisAddr, redisPass string, redisDbnum int, frontendDir string) {
	rdb := RedisDB{}
	rdb.Init(redisAddr, redisPass, redisDbnum)
	a.DB = &rdb

	routes := a.GetRoutes(frontendDir)
	a.Router = mux.NewRouter().StrictSlash(true)
	frontendServer := Logger(http.FileServer(http.Dir(frontendDir)), "frontend")
	a.Router.Path("/").Handler(frontendServer)
	a.Router.PathPrefix("/static").Handler(frontendServer)
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
		Handler: a.Router,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Starting server")
	log.Fatal(srv.ListenAndServe())
}
