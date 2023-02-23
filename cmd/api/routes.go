package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	//router.HandlerFunc(http.MethodGet, "/v1/books", app.listBooks)
	router.HandlerFunc(http.MethodPost, "/v1/books", app.createBook)
	//router.HandlerFunc(http.MethodGet, "/v1/books/:id", app.showMovieHandler)
	//router.HandlerFunc(http.MethodPatch, "/v1/books/:id", app.updateMovieHandler)
	//router.HandlerFunc(http.MethodDelete, "/v1/books/:id", app.deleteMovieHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	return app.recoverPanic(app.rateLimit(router))
}
