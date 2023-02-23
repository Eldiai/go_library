package main

import (
	"errors"
	"net/http"

	"github.com/Eldiai/go_library/internal/data"
	"github.com/Eldiai/go_library/internal/validator"
)

func (app *application) createBook(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title      string   `json:"title"`
		Author     string   `json:"author"`
		Year       int32    `json:"year"`
		Genres     []string `json:"genres"`
		ReleasedAt int32    `json:"released_at"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	book := &data.Book{
		Title:      input.Title,
		Author:     input.Author,
		Year:       input.Year,
		Genres:     input.Genres,
		ReleasedAt: input.ReleasedAt,
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

func (app *application) listBook(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listBooks(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Author string
		Genres []string
		data.Filters
	}
	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Title = app.readString(qs, "author", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "-id", "-title", "-year"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	books, metadata, err := app.models.Books.GetAll(input.Title, input.Author, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"books": books, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateBook(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title      *string  `json:"title"`
		Author     *string  `json:"author"`
		Year       *int32   `json:"year"`
		Genres     []string `json:"genres"`
		ReleasedAt *int32   `json:"released_at"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Title != nil {
		book.Title = *input.Title
	}

	if input.Year != nil {
		book.Year = *input.Year
	}
	if input.ReleasedAt != nil {
		book.ReleasedAt = *input.ReleasedAt
	}
	if input.Genres != nil {
		book.Genres = input.Genres
	}

	v := validator.New()
	if data.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Books.Update(book)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Books.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "book was successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
