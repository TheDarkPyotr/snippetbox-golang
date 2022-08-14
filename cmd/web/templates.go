package main

import (
	"path/filepath"
	"snippetbox/pkg/models"
	"text/template"
	"time"
)

type templateData struct {
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	CurrentYear int
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")

}

//Map to inject the humanDate function into template language
var functions = template.FuncMap{
	"humanDate": humanDate,
}

//Function to avoid reloading for every page the template from disk
func newTemplateCache(dir string) (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	//filepath.Glob: to get a slice of all filepath with extension 'page.tmpl'
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		//Parse page template file in to a template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		//Funcs inject the map containing the user-define template function
		if err != nil {
			return nil, err
		}

		//Add any 'layout' templates to the template set (only base, now)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		//Same for partial templates (e.g. footer)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		//Update cache (e.g. 'home.page.tmpl' as key)
		cache[name] = ts

	}
	return cache, nil
}