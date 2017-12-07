package main

import (
	"log"
	"net/http"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     DB
}

func (a *App) InitRedis(redisAddr, redisPass string, redisDbnum int) {
	rdb := RedisDB{}
	rdb.Init(redisAddr, redisPass, redisDbnum)
	a.DB = &rdb
}

func (a *App) InitRouter() {
	a.Router = mux.NewRouter().StrictSlash(true)

	// Setup the static Vue.js frontend routes
	frontendServer := Logger(http.FileServer(
		&assetfs.AssetFS{
			Asset:     Asset,
			AssetDir:  AssetDir,
			AssetInfo: AssetInfo,
			Prefix:    "frontend"}), "frontend")
	a.Router.Path("/").Handler(frontendServer)
	a.Router.PathPrefix("/static").Handler(frontendServer)

	// Setup the API routes
	routes := a.GetRoutes()
	for _, route := range *routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = a.SessionMiddleware(handler, route.EnsureSession)
		handler = Logger(handler, route.Name)

		a.Router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
}

func (a *App) Run(port string) {
	srv := &http.Server{
		Handler: a.Router,
		Addr:    "localhost:" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Starting server")
	log.Fatal(srv.ListenAndServe())
}
