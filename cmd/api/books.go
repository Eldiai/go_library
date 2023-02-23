package main

import (
	"net/http"

	"github.com/Eldiai/go_library/internal/data"
	"github.com/Eldiai/go_library/internal/validator"
)

func (app *application) createBook(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string   `json:"title"`
		Author string   `json:"author"`
		Year   int32    `json:"year"`
		Genres []string `json:"genres"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	book := &data.Book{
		Title:  input.Title,
		Author: input.Author,
		Year:   input.Year,
		Genres: input.Genres,
	}

	v := validator.New()

	if data.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Books.Insert(book); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err := app.writeJSON(w, http.StatusCreated, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
