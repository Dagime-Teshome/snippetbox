package main

import (
	"fmt"
	"net/http"
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
