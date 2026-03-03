package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(app.handleHome))

	mux.Get("/snippets/list", http.HandlerFunc(app.listSnippets))
	mux.Get("/snippet/create", app.authenticate(app.isLoggedIn(http.HandlerFunc(app.createSnippetForm))))
	mux.Post("/snippet/create", app.authenticate(app.isLoggedIn(http.HandlerFunc(app.createSnippet))))

	mux.Get("/user/login", app.authenticate(http.HandlerFunc(app.loginUserForm)))
	mux.Get("/user/signup", app.authenticate(http.HandlerFunc(app.signupUserForm)))
	mux.Post("/user/login", app.authenticate(http.HandlerFunc(app.loginUser)))
	mux.Post("/user/signup", http.HandlerFunc(app.signupUser))
	mux.Post("/user/logout", app.authenticate(app.isLoggedIn(http.HandlerFunc(app.logoutUser))))

	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))
	mux.Get("/ping", http.HandlerFunc(ping))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return app.recoverPanic(app.logRequest(securityMiddleware(app.Session.Enable(noSurf(mux)))))
}
