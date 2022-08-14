package main

import (
	"net/http"
	"snippetbox/config"

	"github.com/justinas/alice"
)

func (app *application) routes(cfc *config.Specification) http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	middlewareChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	fileServer := http.FileServer(http.Dir(cfc.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//Standard way of chaining middleware without "github.com/justinas/alice"
	//secureMiddleware := secureHeaders(mux)
	//requestLogMiddleware := app.logRequest(secureMiddleware)
	//recoverPanicMiddlware := app.recoverPanic(requestLogMiddleware)

	return middlewareChain.Then(mux)

}
