package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Dagime-Teshome/snippetbox/pkg/forms"
	"github.com/Dagime-Teshome/snippetbox/pkg/models"
)

func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippet.Latest()
	if err != nil {
		app.ServerError(w, err)
	}
	data := &templateData{
		Snippet:  nil,
		Snippets: snippets,
	}
	app.render(w, r, "home.page.tmpl", data)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)

	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}
	title, content, expires := form.Get("title"), form.Get("content"), form.Get("expires")
	id, err := app.snippet.Insert(title, content, expires)

	if err != nil {
		app.ServerError(w, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		if r.URL.Query().Get("id") == "" {
			id = 0
		}
		app.ClientError(w, 505)
		return
	}
	if id <= 0 {
		app.ServerError(w, fmt.Errorf("id cant be less than 0"))
		return
	}

	snippet, err := app.snippet.GetById(id)
	if err == models.ErrNoRecord {
		app.NotFound(w)
		return
	} else if err != nil {
		app.ServerError(w, err)
		return
	}
	data := &templateData{Snippet: snippet}
	app.render(w, r, "show.page.tmpl", data)

	w.Write([]byte(fmt.Sprintf("show snippet for %v", snippet)))
}

func (app *application) listSnippets(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippet.Latest()
	if err != nil {
		app.ServerError(w, err)
	}
	for _, snippet := range snippets {
		fmt.Fprintf(w, "%v\n", snippet)
	}
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}
