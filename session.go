package main

import (
	"context"
	"net/http"
	"time"
)

func (a *App) SessionMiddleware(inner http.Handler, ensureSession bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("session")
		ctx := r.Context()
		if cookie != nil {
			// Add data to context
			valid, username, exists := a.DB.SessionLookup(cookie.Value)
			if valid {
				ctx = context.WithValue(ctx, "Username", username)
				ctx = context.WithValue(ctx, "IsRegistered", exists)
				inner.ServeHTTP(w, r.WithContext(ctx))
				return
			} else {
				ctx = context.WithValue(ctx, "Username", "")
				ctx = context.WithValue(ctx, "IsRegistered", false)
			}
		}
		if ensureSession {
			// make a session with an unregistered username as a guest account
			token := GetRandString(32)
			username := GetRandString(32)
			a.DB.SessionCreate(username, token)
			cookie := http.Cookie{Name: "session", Value: token, Path: "/", Expires: time.Now().Add(30 * 24 * time.Hour)}
			http.SetCookie(w, &cookie)
			ctx = context.WithValue(ctx, "Username", username)
		} else {
			ctx = context.WithValue(ctx, "Username", "")
		}
		ctx = context.WithValue(ctx, "IsRegistered", false)
		inner.ServeHTTP(w, r.WithContext(ctx))
	})
}
