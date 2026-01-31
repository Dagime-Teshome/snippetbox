package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.Write([]byte("404 not found"))
		return
	}
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	template, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "couldn't parse template", 503)
		return
	}

	if err := template.Execute(w, nil); err != nil {
		http.Error(w, "couldn't execute html", 505)
		return
	}
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		w.Write([]byte("Method not allowed"))
		return
	}
	w.Write([]byte("create snippet"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		if r.URL.Query().Get("id") == "" {
			id = 0 // Default value
			fmt.Printf("Using default value: %d\n", id)
		}
		http.Error(w, "id must be a valid number", 505)
		return
	}
	if id <= 0 {
		http.Error(w, "id needs to be greater than zero", 505)
		return
	}

	w.Write([]byte(fmt.Sprintf("show snippet for %v", id)))
}
