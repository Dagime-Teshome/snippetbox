package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Dagime-Teshome/snippetbox/pkg/models"
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
		user := app.authenticatedUser(r)
		if user == nil {
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
	// csrfHandler.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(os.Stderr, "CSRF failure - reason: %v\n", nosurf.Reason(r))
	// 	http.Error(w, http.StatusText(400), 400)
	// }))

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		exists := app.Session.Exists(r, "userID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.user.Get(app.Session.GetInt(r, "userID"))
		if err == models.ErrNoRecord {
			app.Session.Remove(r, "userID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.ServerError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
