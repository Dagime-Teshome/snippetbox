package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/Dagime-Teshome/snippetbox/pkg/forms"
	"github.com/Dagime-Teshome/snippetbox/pkg/models"
)

type templateData struct {
	User        *models.User
	CurrentYear int
	CSRFToken   string
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        *forms.Form
	Flash       string
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		tem, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		tem, err = tem.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		tem, err = tem.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = tem

	}
	return cache, nil
}

func humanDate(t time.Time) string {

	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}
