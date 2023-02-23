package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Books interface {
		Insert(book *Book) error
		Get(id int64) (*Book, error)
		Update(book *Book) error
		Delete(id int64) error
		GetAll(title string, author string, genres []string, filters Filters) ([]*Book, Metadata, error)
	}
	Users  UserModel
	Tokens TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Books:  BookModel{DB: db},
		Tokens: TokenModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
