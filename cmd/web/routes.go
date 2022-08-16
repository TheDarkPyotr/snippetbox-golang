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
	cookieMiddleware := alice.New(app.session.Enable, app.authenticate)

	mux := pat.New()

	mux.Get("/", cookieMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", cookieMiddleware.Append(app.requiredAuthenticatedUser).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", cookieMiddleware.Append(app.requiredAuthenticatedUser).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", cookieMiddleware.ThenFunc(app.showSnippet))

	mux.Get("/user/signup", cookieMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", cookieMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", cookieMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", cookieMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", cookieMiddleware.Append(app.requiredAuthenticatedUser).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)

}
