package main

import (
	"fmt"
	"net/http"
	forms "snippetbox/pkg/forms"
	"snippetbox/pkg/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.database.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippets: s}
	app.render(w, r, "home.page.tmpl", *data)

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.database.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	//Because html/template allow only pass one data
	//we wrap snippets in templateData struct

	data := &templateData{Snippet: s}
	app.render(w, r, "show.page.tmpl", *data)

}

// Add a new createSnippetForm handler, which for now returns a placeholder res
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	//To allow sending multipart data > 10MB
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", templateData{Form: form})
		return
	}

	id, err := app.database.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

}
