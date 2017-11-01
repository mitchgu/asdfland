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
			"SlugDestCreate",
			"POST",
			"/api/slugdest",
			a.SlugDestCreateHandler,
		},
		Route{
			"LinkPreviewGet",
			"GET",
			"/p/{slug}",
			a.KeyPreviewHandler,
		},
		Route{
			"LinkGet",
			"GET",
			"/{slug}",
			a.KeyHandler,
		},
	}
	return &routes
}
