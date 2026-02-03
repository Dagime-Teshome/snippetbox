package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w)
		return
	}
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	template, err := template.ParseFiles(files...)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	if err := template.Execute(w, nil); err != nil {
		app.ServerError(w, err)
		return
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		app.ClientError(w, 405)
		return
	}
	w.Write([]byte("create snippet"))
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		if r.URL.Query().Get("id") == "" {
			id = 0 // Default value
			fmt.Printf("Using default value: %d\n", id)
		}
		app.ClientError(w, 505)
		return
	}
	if id <= 0 {
		app.ServerError(w, fmt.Errorf("id cant be less than 0"))
		return
	}

	w.Write([]byte(fmt.Sprintf("show snippet for %v", id)))
}
