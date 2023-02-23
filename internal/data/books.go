package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type Book struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"-"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	Year       int32     `json:"year,omitempty"`
	Genres     []string  `json:"genres,omitempty"`
	ReleasedAt int32     `json:"released_at"`
}

type BookModel struct {
	DB *sql.DB
}

func (b BookModel) Insert(book *Book) error {
	query := `
INSERT INTO books (title, author, year, genres, released_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at`

	args := []any{book.Title, book.Author, book.Year, pq.Array(book.Genres), book.ReleasedAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return b.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.CreatedAt)

}
func (b BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, nil
	}

	query := `
SELECT id, created_at, title, author,  year, genres, released_at
FROM books
WHERE id = $1`

	var book Book
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.Title,
		&book.Author,
		&book.Year,
		pq.Array(&book.Genres),
		&book.ReleasedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &book, nil

}
func (b BookModel) Update(book *Book) error {
	query := `
UPDATE books
SET title = $1, author = $2, year = $3, genres = $4, released_at = $5
WHERE id = $6
RETURNING title,author,year,genres,released_at`
	args := []interface{}{
		book.Title,
		book.Author,
		book.Year,
		book.ReleasedAt,
		pq.Array(book.Genres),
		book.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, args...).Scan(&book.Title, &book.Author, &book.Year, &book.Genres, &book.ReleasedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (b BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
DELETE FROM books
WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := b.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}
