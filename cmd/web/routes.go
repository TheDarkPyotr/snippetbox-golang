package main

import (
	"net/http"
	"snippetbox/config"
)

func (app *application) routes(cfc *config.Specification) *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	fileServer := http.FileServer(http.Dir(cfc.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux

}
