package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Home",
		"GET",
		"/",
		HomeHandler,
	},
	Route{
		"Admin",
		"GET",
		"/a/",
		AdminHandler,
	},
	Route{
		"DestCreate",
		"Post",
		"/a/dest",
		DestCreateHandler,
	},
	Route{
		"LinkGet",
		"GET",
		"/{key}",
		KeyHandler,
	},
}