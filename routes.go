package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func (a *App) GetRoutes() *Routes {
	var routes = Routes{
		Route{
			"Home",
			"GET",
			"/",
			a.HomeHandler,
		},
		Route{
			"Admin",
			"GET",
			"/a/",
			a.AdminHandler,
		},
		Route{
			"DestCreate",
			"Post",
			"/a/dest",
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