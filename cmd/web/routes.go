package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {

	mux := pat.New()
	mux.Get("/", app.Session.Enable(http.HandlerFunc(app.handleHome)))
	mux.Get("/snippet/create", app.Session.Enable(http.HandlerFunc(app.createSnippetForm)))
	mux.Post("/snippet/create", app.Session.Enable(http.HandlerFunc(app.createSnippet)))
	mux.Get("/snippets/list", app.Session.Enable(http.HandlerFunc(app.listSnippets)))
	mux.Get("/snippet/:id", app.Session.Enable(http.HandlerFunc(app.showSnippet)))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return app.recoverPanic(app.logRequest(securityMiddleware(mux)))
}
