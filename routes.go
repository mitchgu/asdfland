package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func (a *App) GetRoutes(frontendDir string) *Routes {
	var routes = Routes{
		Route{
			"SlugReserve",
			"POST",
			"/api/slug/reserve",
			a.SlugReserveHandler,
		},
		Route{
			"DestCreate",
			"POST",
			"/api/dest",
			a.DestCreateHandler,
		},
		Route{
			"LinkGet",
			"GET",
			"/{key}",
			a.KeyHandler,
		},
	}
	return &routes
}
