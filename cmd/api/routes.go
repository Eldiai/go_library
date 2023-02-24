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

	router.HandlerFunc(http.MethodGet, "/v1/books", app.requirePermission("books:read", app.listBooks))
	router.HandlerFunc(http.MethodPost, "/v1/books", app.requirePermission("books:write", app.createBook))
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", app.requirePermission("books:read", app.listBook))
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id", app.requirePermission("books:write", app.updateBook))
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", app.requirePermission("books:write", app.deleteBook))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUser)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUser)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
