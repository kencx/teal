package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/teal"
)

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
		WHERE b.id=$1`

	if err = tx.QueryRowx(stmt, id).StructScan(&dest); err != nil {
		return nil, fmt.Errorf("db: retrieve book %d failed: %v", id, err)
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
		WHERE b.isbn=$1`

	if err = tx.QueryRowx(stmt, isbn).StructScan(&dest); err != nil {
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
		WHERE b.title=$1`

	if err = tx.QueryRowx(stmt, title).StructScan(&dest); err != nil {
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
		ORDER BY b.title`

	if err = tx.Select(&dest, stmt); err != nil {
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
func (s *Store) CreateBook(b *teal.Book) error {
	if err := s.Tx(func(tx *sqlx.Tx) error {

		b_id, err := insertBook(tx, b)
		if err != nil {
			return err
		}

		// create authors
		authors := parseAuthors(b.Author)
		a_ids, err := insertAuthors(tx, authors)
		if err != nil {
			return err
		}

		// establish book author relationship
		for _, a_id := range a_ids {
			stmt := `INSERT INTO books_authors (book_id, author_id) VALUES ($1, $2)`
			_, err := tx.Exec(stmt, b_id, a_id)
			if err != nil {
				return fmt.Errorf("db: insert to books_authors table failed: %v", err)
			}
		}
		return nil

	}, &sql.TxOptions{}); err != nil {
		return err
	}
	return nil
}

// func (s *Store) UpdateBook(id int, b *teal.Book) error {
// 	if err := s.Tx(func(tx *sqlx.Tx) error {
//
// 		err := updateBook(tx, id, b)
// 		if err != nil {
// 			return err
// 		}
//
// 		authors := parseAuthors(b.Author)
// 		// add new authors
// 		// remove authors
// 		// delete authors with no books
// 		// rename -> add new name and delete old name
// 		// don't delete if there still exists books with old name
//
// 		// update authors table
// 		err := updateAuthor(tx, id, authors)
// 		if err != nil {
// 			return err
// 		}
//
// 		// if add/remove author, need to add/delete row from booksAuthors
// 		// if change author, need to delete row from booksAuthors table and add new row
// 		// change: book change to another existing author
// 		for _, a := range authors {
// 			stmt := `UPDATE books_authors
// 			SET author_id=$1, book_id=$2
// 			WHERE books_authors.author_id=authors.id`
// 			res, err := tx.Exec(stmt, a.ID, id)
// 			if err != nil {
// 				return fmt.Errorf("db: update to books_authors table failed: %v", err)
// 			}
// 			// check for rows affected
// 		}
// 		return nil
//
// 	}, &sql.TxOptions{}); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (s *Store) DeleteBook(id int) error {
	if err := s.Tx(func(tx *sqlx.Tx) error {

		err := deleteBook(tx, id)
		if err != nil {
			return err
		}

		// delete entries from booksAuthors table
		stmt := `DELETE FROM books_authors WHERE books_authors.book_id=$1`
		res, err := tx.Exec(stmt, id)
		if err != nil {
			return fmt.Errorf("db: delete book %d from books_authors table failed: %v", id, err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("db: delete book %d from books_authors table failed: %v", id, err)
		}

		if count == 0 {
			return errors.New("no rows deleted from books_authors table")
		}

		// check for entries in authors with no books
		stmt = `DELETE FROM authors WHERE id NOT IN
				(SELECT author_id FROM books_authors)`
		res, err = tx.Exec(stmt, id)
		if err != nil {
			return fmt.Errorf("db: delete author from authors table failed: %v", err)
		}

		count, err = res.RowsAffected()
		if err != nil {
			return fmt.Errorf("db: delete author from authors table failed: %v", err)
		}

		if count == 0 {
			// TODO what should this return?
			return nil
		}
		return nil

	}, &sql.TxOptions{}); err != nil {
		return err
	}
	return nil
}

// insert book entry to books table
func insertBook(tx *sqlx.Tx, b *teal.Book) (int64, error) {

	stmt := `INSERT INTO books
		(title, description, isbn, numOfPages, rating, state, dateAdded, dateUpdated, dateCompleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	res, err := tx.Exec(stmt,
		b.Title,
		b.Description,
		b.ISBN,
		b.NumOfPages,
		b.Rating,
		b.State,
		b.DateAdded,
		b.DateUpdated,
		b.DateCompleted)

	if err != nil {
		return -1, fmt.Errorf("db: insert to books table failed: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("db: insert to books table failed: %v", err)
	}
	return id, nil
}

func updateBook(tx *sqlx.Tx, id int, b *teal.Book) error {

	stmt := `UPDATE books b
			SET title=$1, description=$2, isbn=$3, numOfPages=$4,
			rating=$5, state=$6, dateAdded=$7, dateUpdated=$8, dateCompleted=$9
			WHERE b.id=$10`
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
		return fmt.Errorf("db: unable to update book %d: %w", id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to update book %d: %w", id, err)
	}

	if count == 0 {
		return errors.New("db: no books updated")
	}
	return nil
}

// delete book entry from books table
func deleteBook(tx *sqlx.Tx, id int) error {

	stmt := `DELETE from books WHERE books.id=$1`
	res, err := tx.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("db: unable to delete book %d: %w", id, err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to delete book %d: %w", id, err)
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
