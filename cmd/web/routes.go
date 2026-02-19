package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {

	mux := pat.New()
	mux.Get("/", app.Session.Enable(http.HandlerFunc(app.handleHome)))
	mux.Get("/snippets/list", app.Session.Enable(http.HandlerFunc(app.listSnippets)))
	// create snippets
	mux.Get("/snippet/create", app.Session.Enable(app.isLoggedIn(http.HandlerFunc(app.createSnippetForm))))
	mux.Post("/snippet/create", app.Session.Enable(app.isLoggedIn(http.HandlerFunc(app.createSnippet))))
	// login
	mux.Get("/user/login", app.Session.Enable(http.HandlerFunc(app.loginUserForm)))
	mux.Post("/user/login", app.Session.Enable(http.HandlerFunc(app.loginUser)))
	// sign up
	mux.Get("/user/signup", app.Session.Enable(http.HandlerFunc(app.signupUserForm)))
	mux.Post("/user/signup", app.Session.Enable(http.HandlerFunc(app.signupUser)))
	// log out
	mux.Post("/user/logout", app.Session.Enable(app.isLoggedIn(http.HandlerFunc(app.logoutUser))))
	// show snippet
	mux.Get("/snippet/:id", app.Session.Enable(app.isLoggedIn(http.HandlerFunc(app.showSnippet))))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return app.recoverPanic(app.logRequest(securityMiddleware(mux)))
}
