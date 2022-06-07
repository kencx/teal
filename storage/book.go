package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/teal"
)

func (r *Store) GetBook(id int) (*teal.Book, error) {
	tx, err := r.db.Begin()
func getBook(s *Store, id int) (*teal.Book, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer endTx(tx, err)

	var b teal.Book
	if err := tx.Get(&b, `SELECT
		id,
		title,
		description,
		isbn,
		numOfPages,
		rating,
		state,
		dateAdded,
		dateUpdated,
		dateCompleted
		FROM books WHERE id = $1`, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else if err != nil {
			return nil, fmt.Errorf("db: unable to fetch book %d: %w", id, err)
		}
	}
	return &b, nil
}

func (r *Store) GetBookByTitle(title string) (*teal.Book, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer endTx(tx, err)

	var b teal.Book
	if err := tx.QueryRow(`SELECT
		id,
		title,
		author,
		isbn
	FROM books WHERE title = $1`, title).Scan(
		&b.ID, &b.Title, &b.Author, &b.ISBN); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		} else if err != nil {
			return nil, fmt.Errorf("db: unable to fetch book %q: %w", title, err)
		}
	}

	return &b, nil
}

func (r *Store) GetAllBooks() ([]*teal.Book, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer endTx(tx, err)

	var books []*teal.Book
	rows, err := tx.Query("SELECT id, title, author, isbn FROM books")
	if err != nil {
		return nil, fmt.Errorf("db: unable to fetch books: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var b teal.Book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN)
		if err != nil {
			return nil, err
		}
		books = append(books, &b)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("db: unable to fetch books: %w", err)
	}
	return books, nil
}

func insertBook(s *Store, b *teal.Book) (int, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return -1, err
	}
	defer endTx(tx, err)

	res, err := tx.Exec(`INSERT into books (
		title,
		description,
		isbn,
		numOfPages,
		rating,
		state,
		dateAdded,
		dateUpdated,
		dateCompleted)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
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
		return -1, fmt.Errorf("db: unable to create book %q: %w", b.Title, err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("db: unable to create book %q: %w", b.Title, err)
	}
	b.ID = int(lastId)

	return int(lastId), nil
}

func (r *Store) UpdateBook(id int, b *teal.Book) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer endTx(tx, err)

	result, err := tx.Exec("UPDATE books SET title=$1, author=$2, isbn=$3 WHERE id=$4", b.Title, b.Author, b.ISBN, id)
	if err != nil {
		return fmt.Errorf("db: unable to update book %d: %w", id, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to update book %d: %w", id, err)
	}

	if count == 0 {
		return errors.New("db: no books updated")
	}

	return nil
}

func deleteBook(s *Store, id int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer endTx(tx, err)

	result, err := tx.Exec("DELETE FROM books WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("db: unable to delete book %d: %w", id, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to delete book %d: %w", id, err)
	}

	if count == 0 {
		return errors.New("db: no books removed")
	}
	return nil
}

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
