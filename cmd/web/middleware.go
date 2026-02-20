package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func securityMiddleware(next http.Handler) http.Handler {

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handler)
}

func (app *application) logRequest(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("IP:%s , url:%s", r.RemoteAddr, r.URL.Path)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handler)
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.ServerError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handler)
}

func (app *application) isLoggedIn(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		id := app.authenticatedUser(r)
		if id == 0 {
			http.Redirect(w, r, "/user/login", 302)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handler)
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}
