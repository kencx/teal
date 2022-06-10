package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/teal"
)

type contextKey string

type BookAuthorDest struct {
	*teal.Book
	Author_string string
}

func (s *Store) RetrieveBookWithID(id int) (*teal.Book, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest BookAuthorDest
	stmt := `SELECT b.*, GROUP_CONCAT(a.name) AS author_string
		FROM books b
		INNER JOIN books_authors ba ON ba.book_id=b.id
		INNER JOIN authors a ON ba.author_id=a.id
		WHERE b.id=$1
		GROUP BY b.id;`

	err = tx.QueryRowx(stmt, id).StructScan(&dest)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve book id %d failed: %v", id, err)
	}

	dest.Author = strings.Split(dest.Author_string, ",")
	return dest.Book, nil
}

func (s *Store) RetrieveBookWithISBN(isbn string) (*teal.Book, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest BookAuthorDest
	stmt := `SELECT b.*, GROUP_CONCAT(a.name) AS author_string
		FROM books b
		INNER JOIN books_authors ba ON ba.book_id=b.id
		INNER JOIN authors a ON ba.author_id=a.id
		WHERE b.isbn=$1
		GROUP BY b.isbn;`

	err = tx.QueryRowx(stmt, isbn).StructScan(&dest)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve book isbn %q failed: %v", isbn, err)
	}

	dest.Author = strings.Split(dest.Author_string, ",")
	return dest.Book, nil
}

func (s *Store) RetrieveBookWithTitle(title string) (*teal.Book, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest BookAuthorDest
	stmt := `SELECT b.*, GROUP_CONCAT(a.name) AS author_string
		FROM books b
		INNER JOIN books_authors ba ON ba.book_id=b.id
		INNER JOIN authors a ON ba.author_id=a.id
		WHERE b.title=$1
		GROUP BY b.title;`

	err = tx.QueryRowx(stmt, title).StructScan(&dest)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve book title %q failed: %v", title, err)
	}

	dest.Author = strings.Split(dest.Author_string, ",")
	return dest.Book, nil
}

func (s *Store) RetrieveAllBooks() ([]*teal.Book, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest []BookAuthorDest
	stmt := `SELECT b.*, GROUP_CONCAT(a.name) AS author_string
		FROM books b
		INNER JOIN books_authors ba ON ba.book_id=b.id
		INNER JOIN authors a ON ba.author_id=a.id
		GROUP BY b.title
		ORDER BY b.id;`

	err = tx.Select(&dest, stmt)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve all books failed: %v", err)
	}

	var books []*teal.Book
	for _, row := range dest {
		row.Author = strings.Split(row.Author_string, ",")
		books = append(books, row.Book)
	}
	return books, nil
}

// Create book entry in books, author entries in authors
func (s *Store) CreateBook(ctx context.Context, b *teal.Book) (*teal.Book, error) {
	if err := s.Tx(func(tx *sqlx.Tx) error {

		book, err := insertBook(tx, b)
		if err != nil {
			return err
		}
		ctx = context.WithValue(ctx, contextKey("book_id"), book)

		// create authors
		authors := parseAuthors(b.Author)
		a_ids, err := insertOrGetAuthors(tx, authors)
		if err != nil {
			return err
		}

		// establish new book author relationship
		for _, a_id := range a_ids {
			err := linkBookToAuthor(tx, int64(book.ID), a_id)
			if err != nil {
				return err
			}
		}
		return nil

	}, &sql.TxOptions{}); err != nil {
		return nil, err
	}

	book, ok := ctx.Value(contextKey("book_id")).(*teal.Book)
	if !ok {
		return nil, fmt.Errorf("db: failed to cast book")
	}
	return book, nil
}

// Update book details.
// For authors, a new author row is created for each new author
// No authors are deleted, unless it has no relationship with any books
func (s *Store) UpdateBook(id int, b *teal.Book) error {
	if err := s.Tx(func(tx *sqlx.Tx) error {

		err := updateBook(tx, id, b)
		if err != nil {
			return err
		}

		book, err := s.RetrieveBookWithID(id)
		if err != nil {
			return err
		}
		// if author is changed
		if !reflect.DeepEqual(book.Author, b.Author) {

			// Add new or get existing authors
			// Renaming an author should not update the same author row for other books
			// Always create a new author row, never update the original in this case
			authors := parseAuthors(b.Author)
			a_ids, err := insertOrGetAuthors(tx, authors)
			if err != nil {
				return err
			}

			// establish NEW or EXISTING book author relationships
			for _, a_id := range a_ids {
				err := linkBookToAuthor(tx, int64(id), a_id)
				if err != nil {
					return err
				}
			}

			// remove broken book author relationships
			stmt := `DELETE FROM books_authors
				WHERE book_id=? AND author_id NOT IN (?);`
			query, args, err := sqlx.In(stmt, id, a_ids)
			if err != nil {
				return fmt.Errorf("db: unlink book %d to authors %v in book_authors failed: %v", id, a_ids, err)
			}
			_, err = tx.Exec(query, args...)
			if err != nil {
				return fmt.Errorf("db: update book %d in books_authors failed: %v", id, err)
			}

			// delete authors with no books
			err = deleteAuthorsWithNoBooks(tx)
			if err != nil {
				return err
			}
		}
		return nil

	}, &sql.TxOptions{}); err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteBook(id int) error {
	if err := s.Tx(func(tx *sqlx.Tx) error {

		err := deleteBook(tx, id)
		if err != nil {
			return err
		}

		// delete all book entries from booksAuthors table
		stmt := `DELETE FROM books_authors WHERE book_id=$1;`
		res, err := tx.Exec(stmt, id)
		if err != nil {
			return fmt.Errorf("db: delete book %d from books_authors failed: %v", id, err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("db: delete book %d from books_authors failed: %v", id, err)
		}

		if count == 0 {
			return errors.New("no rows deleted from books_authors table")
		}

		deleteAuthorsWithNoBooks(tx)
		return nil

	}, &sql.TxOptions{}); err != nil {
		return err
	}
	return nil
}

// insert book entry to books table
func insertBook(tx *sqlx.Tx, b *teal.Book) (*teal.Book, error) {
	var book teal.Book

	stmt := `INSERT INTO books
		(title, description, isbn, numOfPages, rating, state, dateAdded, dateUpdated, dateCompleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;`
	err := tx.QueryRowx(stmt,
		b.Title,
		b.Description,
		b.ISBN,
		b.NumOfPages,
		b.Rating,
		b.State,
		b.DateAdded,
		b.DateUpdated,
		b.DateCompleted).StructScan(&book)

	if err != nil {
		return nil, fmt.Errorf("db: insert to books table failed: %v", err)
	}
	book.Author = b.Author
	return &book, nil
}

func updateBook(tx *sqlx.Tx, id int, b *teal.Book) error {

	stmt := `UPDATE books
			SET title=$1,
			description=$2,
			isbn=$3,
			numOfPages=$4,
			rating=$5,
			state=$6,
			dateAdded=$7,
			dateUpdated=$8,
			dateCompleted=$9
			WHERE id=$10;`
	res, err := tx.Exec(stmt,
		b.Title,
		b.Description,
		b.ISBN,
		b.NumOfPages,
		b.Rating,
		b.State,
		b.DateAdded,
		b.DateUpdated,
		b.DateCompleted,
		b.ID)

	if err != nil {
		return fmt.Errorf("db: update book %d failed: %v", id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: update book %d failed: %v", id, err)
	}

	if count == 0 {
		return errors.New("db: no books updated")
	}
	return nil
}

// delete book entry from books table
func deleteBook(tx *sqlx.Tx, id int) error {

	stmt := `DELETE from books WHERE books.id=$1;`
	res, err := tx.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("db: delete book %d failed: %v", id, err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: delete book %d failed: %v", id, err)
	}
	if count == 0 {
		return errors.New("db: no books removed")
	}
	return nil
}

// functional Tx helper for Exec statements
func (s *Store) Tx(fn func(tx *sqlx.Tx) error, opts *sql.TxOptions) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	if err = fn(tx); err != nil {
		return err
	}
	return nil
}

// Tx rollback and commit helper, use with defer
func endTx(tx *sqlx.Tx, err error) error {
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p)
	} else if err != nil {
		tx.Rollback()
		return nil
	} else {
		return tx.Commit()
	}
}
