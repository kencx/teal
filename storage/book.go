package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"unsafe"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/kencx/teal"
	"github.com/kencx/teal/storage/sqlite/model"
	. "github.com/kencx/teal/storage/sqlite/table"
)

var (
	FULL_SELECT = SELECT(Books.AllColumns, Authors.AllColumns).
		FROM(Books.
			INNER_JOIN(BooksAuthors, BooksAuthors.BookID.EQ(Books.ID)).
			INNER_JOIN(Authors, BooksAuthors.AuthorID.EQ(Authors.ID)))
)

func (s *Store) RetrieveBookWithID(id int) (*teal.Book, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest BookDest
	if err := FULL_SELECT.WHERE(Books.ID.EQ(Int(int64(id)))).Query(tx, &dest); err != nil {
		return nil, fmt.Errorf("db: retrieve book %q failed: %v", id, err)
	}
	res := modelToBook(dest)
	return res, nil
}

func (s *Store) RetrieveBookWithTitle(title string) (*teal.Book, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest BookDest
	if err := FULL_SELECT.WHERE(Books.Title.EQ(String(title))).Query(tx, &dest); err != nil {
		return nil, fmt.Errorf("db: retrieve book %q failed: %v", title, err)
	}
	res := modelToBook(dest)
	return res, nil
}

// Create book entry in books, author entries in authors
func (s *Store) CreateBook(b *teal.Book) error {
	if err := s.Tx(func(tx *sqlx.Tx) error {

		bookDest := bookToModel(b)
		book := bookDest.Books
		authors := bookDest.Authors

		b_id, err := insertBook(tx, book)
		if err != nil {
			return err
		}
		a_ids, err := insertAuthors(tx, authors)
		if err != nil {
			return err
		}

		// insert into BookAuthors table for each author
		for _, a := range a_ids {
			_, err := BooksAuthors.INSERT(BooksAuthors.BookID, BooksAuthors.AuthorID).
				VALUES(b_id, a).Exec(tx)
			if err != nil {
				return fmt.Errorf("db: insert to book_authors table failed: %v", err)
			}
		}
		return nil

	}, &sql.TxOptions{}); err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateBook(id int, b *teal.Book) error {
	if err := s.Tx(func(tx *sqlx.Tx) error {
		bookDest := bookToModel(b)
		err := updateBook(tx, id, bookDest.Books)
		if err != nil {
			return err
		}

		// update authors
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

		// delete authors?
		// delete entry from booksAuthors table
		return nil

	}, &sql.TxOptions{}); err != nil {
		return err
	}
	return nil
}

// insert book entry to books table
func insertBook(tx *sqlx.Tx, b model.Books) (*int32, error) {

	var book model.Books
	err := Books.INSERT(Books.MutableColumns).MODEL(b).RETURNING(Books.ID).Query(tx, &book)
	if err != nil {
		return nil, fmt.Errorf("db: insert to books table failed: %v", err)
	}
	return book.ID, nil
}

func updateBook(tx *sqlx.Tx, id int, b model.Books) error {

	res, err := Books.UPDATE(Books.AllColumns).
		MODEL(b).
		WHERE(Books.ID.EQ(Int(int64(id)))).Exec(tx)

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

	res, err := Books.DELETE().WHERE(Books.ID.EQ(Int(int64(id)))).Exec(tx)
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

type BookDest struct {
	model.Books
	Authors []model.Authors
}

func modelToBook(dest BookDest) *teal.Book {

	var authors teal.Authors
	for _, a := range dest.Authors {
		res := teal.Author{
			ID:   int(*a.ID),
			Name: a.Name,
		}
		authors = append(authors, res)
	}

	return &teal.Book{
		ID:         int(*dest.ID),
		Title:      dest.Title,
		ISBN:       dest.Isbn,
		Author:     authors,
		NumOfPages: int(*dest.NumOfPages),
		Rating:     int(*dest.Rating),
		State:      dest.State,
	}
}

func bookToModel(b *teal.Book) BookDest {

	books := model.Books{
		Title: b.Title,
		Isbn:  b.ISBN,
		// TODO replace unsafe casting
		NumOfPages: (*int32)(unsafe.Pointer(&b.NumOfPages)),
		Rating:     (*int32)(unsafe.Pointer(&b.Rating)),
		State:      b.State,
	}

	var authors []model.Authors
	for _, a := range b.Author {
		author := model.Authors{
			Name: a.Name,
		}
		authors = append(authors, author)
	}
	return BookDest{books, authors}
}

// func getBook(s *Store, id int) (*teal.Book, error) {
// 	tx, err := s.db.Beginx()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer endTx(tx, err)
//
// 	var b teal.Book
// 	if err := tx.Get(&b, `SELECT
// 		id,
// 		title,
// 		description,
// 		isbn,
// 		numOfPages,
// 		rating,
// 		state,
// 		dateAdded,
// 		dateUpdated,
// 		dateCompleted
// 		FROM books WHERE id = $1`, id); err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		} else if err != nil {
// 			return nil, fmt.Errorf("db: unable to fetch book %d: %w", id, err)
// 		}
// 	}
// 	return &b, nil
// }
