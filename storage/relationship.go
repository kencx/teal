package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func linkBookToAuthor(tx *sqlx.Tx, book_id, author_id int64) error {
	stmt := `INSERT or IGNORE INTO books_authors (book_id, author_id) VALUES ($1, $2);`

	_, err := tx.Exec(stmt, book_id, author_id)
	if err != nil {
		return fmt.Errorf("db: link book %d to author %d in books_authors failed: %v", book_id, author_id, err)
	}
	return nil
}

// get list of author names from given isbn
func getAuthorsFromBook(tx *sqlx.Tx, book_isbn string) ([]string, error) {

	var dest []struct {
		Title string
		Name  string
	}
	stmt := `SELECT b.title, a.name
		FROM books_authors ba
		JOIN books b ON b.id=ba.book_id
		JOIN authors a ON a.id=ba.author_id
		WHERE b.isbn=$1`

	if err := tx.Select(&dest, stmt, book_isbn); err != nil {
		return nil, err
	}

	var authors []string
	for _, v := range dest {
		authors = append(authors, v.Name)
	}
	return authors, nil
}

func getBooksFromAuthor(tx *sqlx.Tx, name string) ([]string, error) {
	var dest []struct {
		Title string
		Name  string
	}

	stmt := `SELECT b.title, a.name
		FROM books_authors ba
		JOIN books b ON b.id=ba.book_id
		JOIN authors a ON a.id=ba.author_id
		WHERE a.name=$1`

	if err := tx.Select(&dest, stmt, name); err != nil {
		return nil, err
	}

	var books []string
	for _, v := range dest {
		books = append(books, v.Title)
	}
	return books, nil
}
