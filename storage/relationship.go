package storage

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/teal"
)

// get list of author names from given isbn
func (s *Store) RetrieveAuthorsFromBook(book_isbn string) ([]string, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

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

func (s *Store) RetrieveBooksFromAuthor(name string) ([]*teal.Book, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)
	var dest []BookAuthorDest

	stmt := `SELECT b.*, group_concat(a.name) as author_string
		FROM books_authors ba
		JOIN books b ON b.id=ba.book_id
		JOIN authors a ON a.id=ba.author_id
		WHERE b.id IN (SELECT ba.book_id
			FROM books_authors ba
			JOIN authors a ON a.id=ba.author_id
			WHERE a.name=$1)
		GROUP BY b.id`

	if err := tx.Select(&dest, stmt, name); err != nil {
		return nil, err
	}

	var result []*teal.Book
	for _, v := range dest {
		r := v.Book
		r.Author = strings.Split(v.Author_string, ",")
		result = append(result, r)
	}

	return result, nil
}

func linkBookToAuthor(tx *sqlx.Tx, book_id, author_id int64) error {
	stmt := `INSERT or IGNORE INTO books_authors (book_id, author_id) VALUES ($1, $2);`

	_, err := tx.Exec(stmt, book_id, author_id)
	if err != nil {
		return fmt.Errorf("db: link book %d to author %d in books_authors failed: %v", book_id, author_id, err)
	}
	return nil
}

func linkBookToAuthors(tx *sqlx.Tx, book_id int64, author_ids []int64) error {
	type value struct {
		Book_id   int64
		Author_id int64
	}

	var args = []*value{}
	for _, a := range author_ids {
		args = append(args, &value{
			Book_id:   book_id,
			Author_id: a,
		})
	}

	stmt := `INSERT or IGNORE INTO books_authors (book_id, author_id) VALUES (:book_id, :author_id);`
	_, err := tx.NamedExec(stmt, args)
	if err != nil {
		return fmt.Errorf("db: link book %d to authors %d in books_authors failed: %v", book_id, author_ids, err)
	}
	return nil
}

func unlinkBookFromAuthor(tx *sqlx.Tx, book_id, author_id int64) error {
	stmt := `DELETE FROM books_authors WHERE book_id=? AND author_id!=?;`
	_, err := tx.Exec(stmt, book_id, author_id)
	if err != nil {
		return fmt.Errorf("db: unlink author %v from book %v in book_authors failed: %v", author_id, book_id, err)
	}
	return nil
}

func unlinkBookFromAuthors(tx *sqlx.Tx, book_id int64, author_ids []int64) error {
	stmt := `DELETE FROM books_authors WHERE book_id=? AND author_id NOT IN (?);`
	query, args, err := sqlx.In(stmt, book_id, author_ids)
	if err != nil {
		return fmt.Errorf("db: unlink authors %v from book %v in book_authors failed: %v", author_ids, book_id, err)
	}
	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("db: unlink authors %v from book %v in book_authors failed: %v", author_ids, book_id, err)
	}
	return nil
}
