package main

import (
	"net/http"
	"snippetbox/config"

	"github.com/bmizerany/pat"

	"github.com/justinas/alice"
)

func (app *application) routes(cfc *config.Specification) http.Handler {

	//Standard middleware for request without cookie
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	//Use the cookieMiddleware by wrapping the handlers functions
	cookieMiddleware := alice.New(app.session.Enable)
	mux := pat.New()

	mux.Get("/", cookieMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", cookieMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", cookieMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", cookieMiddleware.ThenFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)

}
