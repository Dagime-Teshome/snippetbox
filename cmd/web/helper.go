package main

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
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
	template, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	return template, nil
}
