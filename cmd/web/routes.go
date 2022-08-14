package main

import "net/http"

func (app *application) routes(cfc *Config) *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	fileServer := http.FileServer(http.Dir(cfc.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux

}
