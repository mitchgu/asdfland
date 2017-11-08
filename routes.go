package main

import "net/http"

type Route struct {
	Name          string
	Method        string
	Pattern       string
	EnsureSession bool
	HandlerFunc   http.HandlerFunc
}

type Routes []Route

func (a *App) GetRoutes() *Routes {
	var routes = Routes{
		Route{
			"SlugReserve",
			"POST",
			"/api/slug/reserve",
			true,
			a.SlugReserveHandler,
		},
		Route{
			"SlugDestCreate",
			"POST",
			"/api/slugdest",
			true,
			a.SlugDestCreateHandler,
		},
		Route{
			"UserCreate",
			"POST",
			"/api/user/create",
			false,
			a.UserCreateHandler,
		},
		Route{
			"UserLogin",
			"POST",
			"/api/user/login",
			false,
			a.UserLoginHandler,
		},
		Route{
			"UserLogout",
			"get",
			"/api/user/logout",
			false,
			a.UserLogoutHandler,
		},
		Route{
			"SessionGet",
			"GET",
			"/api/session",
			true,
			a.SessionGetHandler,
		},
		Route{
			"LinkPreviewGet",
			"GET",
			"/p/{slug}",
			false,
			a.KeyPreviewHandler,
		},
		Route{
			"LinkGet",
			"GET",
			"/{slug}",
			false,
			a.KeyHandler,
		},
	}
	return &routes
}
