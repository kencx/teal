package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/kencx/teal/pkg"
)

func (db *DB) GetBook(id int) (*pkg.Book, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}
	defer endTx(tx, err)

	var b pkg.Book
	if err := tx.QueryRow(`SELECT
		id,
		title,
		author,
		isbn
	FROM books WHERE id = $1`, id).Scan(
		&b.ID, &b.Title, &b.Author, &b.ISBN); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		} else if err != nil {
			return nil, fmt.Errorf("db: unable to fetch book %d: %w", id, err)
		}
	}

	return &b, nil
}

func (db *DB) GetBookByTitle(title string) (*pkg.Book, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}
	defer endTx(tx, err)

	var b pkg.Book
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

func (db *DB) GetAllBooks() ([]*pkg.Book, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}
	defer endTx(tx, err)

	var books []*pkg.Book
	rows, err := tx.Query("SELECT id, title, author, isbn FROM books")
	if err != nil {
		return nil, fmt.Errorf("db: unable to fetch books: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var b pkg.Book
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

func (db *DB) CreateBook(b *pkg.Book) (int, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return -1, err
	}
	defer endTx(tx, err)

	res, err := tx.Exec("INSERT INTO books (title, author, isbn) VALUES ($1, $2, $3)", b.Title, b.Author, b.ISBN)
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

func (db *DB) UpdateBook(id int, b *pkg.Book) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer endTx(tx, err)

	_, err = tx.Exec("UPDATE books SET title=$1, author=$2, isbn=$3 WHERE id=$4", b.Title, b.Author, b.ISBN, id)
	if err != nil {
		return fmt.Errorf("db: unable to update book %d: %w", id, err)
	}

	return nil
}

func (db *DB) DeleteBook(id int) error {
	tx, err := db.db.Begin()
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

func endTx(tx *sql.Tx, err error) error {
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

// func ReadBook(id int) error {
// 	i := findIndexByBookId(id)
// 	if i == -1 {
// 		return ErrBookNotFound
// 	}
//
// 	b := bookList[i]
// 	b.Read = "read"
// 	// update date completed
// 	return nil
// }
//
// func ReadingBook(id int) error {
// 	i := findIndexByBookId(id)
// 	if i == -1 {
// 		return ErrBookNotFound
// 	}
//
// 	b := bookList[i]
// 	b.Read = "reading"
// 	// update date updated
// 	return nil
// }
//
// func UnreadBook(id int) error {
// 	i := findIndexByBookId(id)
// 	if i == -1 {
// 		return ErrBookNotFound
// 	}
//
// 	b := bookList[i]
// 	b.Read = "unread"
// 	return nil
// }
//
// func TagBook(id int) error {
// 	return nil
// }
//
// func UntagBook(id int) error {
// 	return nil
// }
//
// func AddBookToCategory(id int) error {
// 	return nil
// }
//
// func RemoveBookFromCategory(id int) error {
// 	return nil
// }
