package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

func (app *application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) NotFound(w http.ResponseWriter) {
	template, err := returnTemplate("")
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if err := template.Execute(w, nil); err != nil {
		app.ServerError(w, err)
		return
	}
	// app.ClientError(w, http.StatusNotFound)
}

func returnTemplate(page string) (*template.Template, error) {
	templatePath := ""
	switch page {
	case "home":
		templatePath = "./ui/html/home.page.tmpl"
	case "snippet":
		templatePath = "./ui/html/show.page.tmpl"
	default:
		templatePath = "./ui/html/404.page.tmpl"
	}

	files := []string{
		templatePath,
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	tm, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	return tm, nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.ServerError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Execute the template set, passing in the dynamic data.
	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.CurrentYear = time.Now().Year()
	td.Flash = app.Session.PopString(r, "flash")
	td.UserID = app.authenticatedUser(r)

	return td
}

func (app *application) authenticatedUser(r *http.Request) int {
	return app.Session.GetInt(r, "userID")
}
