package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/teal"
	tcontext "github.com/kencx/teal/context"
)

type BookStore struct {
	db *sqlx.DB
}

type BookModel struct {
	*teal.Book
	AuthorString string `db:"author_string"`
}

func (bs *BookStore) Get(id int) (*teal.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := bs.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest BookModel
	stmt := `SELECT b.*, GROUP_CONCAT(a.name) AS author_string
		FROM books b
		INNER JOIN books_authors ba ON ba.book_id=b.id
		INNER JOIN authors a ON ba.author_id=a.id
		WHERE b.id=$1
		GROUP BY b.id;`

	err = tx.QueryRowxContext(ctx, stmt, id).StructScan(&dest)
	if err == sql.ErrNoRows {
		return nil, teal.ErrDoesNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve book id %d failed: %v", id, err)
	}

	dest.Author = strings.Split(dest.AuthorString, ",")
	return dest.Book, nil
}

func (bs *BookStore) GetByISBN(isbn string) (*teal.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := bs.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest BookModel
	stmt := `SELECT b.*, GROUP_CONCAT(a.name) AS author_string
		FROM books b
		INNER JOIN books_authors ba ON ba.book_id=b.id
		INNER JOIN authors a ON ba.author_id=a.id
		WHERE b.isbn=$1
		GROUP BY b.isbn;`

	err = tx.QueryRowxContext(ctx, stmt, isbn).StructScan(&dest)
	if err == sql.ErrNoRows {
		return nil, teal.ErrDoesNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve book isbn %q failed: %v", isbn, err)
	}

	dest.Author = strings.Split(dest.AuthorString, ",")
	return dest.Book, nil
}

func (bs *BookStore) GetByTitle(title string) (*teal.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := bs.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest BookModel
	stmt := `SELECT b.*, GROUP_CONCAT(a.name) AS author_string
		FROM books b
		INNER JOIN books_authors ba ON ba.book_id=b.id
		INNER JOIN authors a ON ba.author_id=a.id
		WHERE b.title=$1
		GROUP BY b.title;`

	err = tx.QueryRowxContext(ctx, stmt, title).StructScan(&dest)
	if err == sql.ErrNoRows {
		return nil, teal.ErrDoesNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve book title %q failed: %v", title, err)
	}

	dest.Author = strings.Split(dest.AuthorString, ",")
	return dest.Book, nil
}

func (bs *BookStore) GetAll() ([]*teal.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := bs.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest []BookModel
	stmt := `SELECT b.*, GROUP_CONCAT(a.name) AS author_string
		FROM books b
		INNER JOIN books_authors ba ON ba.book_id=b.id
		INNER JOIN authors a ON ba.author_id=a.id
		GROUP BY b.id
		ORDER BY b.id;`

	err = tx.SelectContext(ctx, &dest, stmt)
	// sqlx Select does not seem to return sql.ErrNoRows
	// related issue: https://github.com/jmoiron/sqlx/issues/762#issuecomment-1062649063
	// if err == sql.ErrNoRows {
	// 	return nil, err
	// }
	if err != nil {
		return nil, fmt.Errorf("db: retrieve all books failed: %v", err)
	}
	if len(dest) == 0 {
		return nil, teal.ErrNoRows
	}

	var books []*teal.Book
	for _, row := range dest {
		row.Author = strings.Split(row.AuthorString, ",")
		books = append(books, row.Book)
	}
	return books, nil
}

// Create a book entry in books, author entries in authors and establishes the necessary
// book author relationships
func (bs *BookStore) Create(b *teal.Book) (*teal.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := Tx(bs.db, ctx, func(tx *sqlx.Tx) error {

		book, err := insertBook(tx, b)
		if err != nil {
			return err
		}
		// save created entity to context to extract after transaction
		ctx = tcontext.WithBook(ctx, book)

		// create authors
		authors := parseAuthors(b.Author)
		a_ids, err := insertOrGetAuthors(tx, authors)
		if err != nil {
			return err
		}

		// establish new book author relationship
		err = linkBookToAuthors(tx, int64(book.ID), a_ids)
		if err != nil {
			return err
		}
		return nil

	}); err != nil {
		return nil, err
	}

	book, err := tcontext.GetBook(ctx)
	if err != nil {
		return nil, err
	}
	return book, nil
}

// Update book details.
// For authors, a new author row is created for each new author
// No authors are deleted, unless it has no relationship with any books
func (bs *BookStore) Update(id int, b *teal.Book) (*teal.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := Tx(bs.db, ctx, func(tx *sqlx.Tx) error {

		err := updateBook(tx, id, b)
		if err != nil {
			return err
		}

		current_authors, err := bs.GetAuthorsFromBook(id)
		if err != nil {
			return err
		}

		// TODO check if order matters
		if !reflect.DeepEqual(current_authors, b.Author) {

			// Add new or get existing authors
			// Renaming an author should not update the same author row for other books
			// Always create a new author row, never update the original in this case
			authors := parseAuthors(b.Author)
			a_ids, err := insertOrGetAuthors(tx, authors)
			if err != nil {
				return err
			}

			// establish book author relationships with NEW or EXISTING authors
			linkBookToAuthors(tx, int64(id), a_ids)

			// remove broken book author relationships
			unlinkBookFromAuthors(tx, int64(id), a_ids)

			// delete authors with no books
			err = deleteAuthorsWithNoBooks(tx)
			if err != nil {
				return err
			}
		}
		return nil

	}); err != nil {
		return nil, err
	}
	return b, nil
}

func (bs *BookStore) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := Tx(bs.db, ctx, func(tx *sqlx.Tx) error {

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

	}); err != nil {
		return err
	}
	return nil
}

// insert book entry to books table
func insertBook(tx *sqlx.Tx, b *teal.Book) (*teal.Book, error) {

	stmt := `INSERT INTO books
		(title, description, isbn, numOfPages, rating, state, dateAdded, dateUpdated, dateCompleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`
	err := tx.QueryRowx(stmt,
		b.Title,
		b.Description,
		b.ISBN,
		b.NumOfPages,
		b.Rating,
		b.State,
		b.DateAdded,
		b.DateUpdated,
		b.DateCompleted).StructScan(b)

	if err != nil {
		return nil, fmt.Errorf("db: insert to books table failed: %v", err)
	}
	return b, nil
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
		id)

	if err != nil {
		return fmt.Errorf("db: update book %d failed: %v", id, err)
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

// delete book entry from books table
func deleteBook(tx *sqlx.Tx, id int) error {

	stmt := `DELETE from books WHERE id=$1;`
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
