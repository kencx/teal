package storage

import (
	"database/sql"
	"errors"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/kencx/teal"
	"github.com/kencx/teal/storage/sqlite/model"
	. "github.com/kencx/teal/storage/sqlite/table"
)

type AuthorDest struct {
	model.Authors
	Books []model.Books
}

func authorToModel(a *teal.Author) AuthorDest {
	return AuthorDest{}
}

func modelToAuthor(a AuthorDest) *teal.Author {
	return nil
}

// should this return struct or []string
func (s *Store) RetrieveAllAuthorNames() (*teal.Author, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest *teal.Author
	if err := SELECT(Authors.AllColumns).FROM(Authors).Query(tx, &dest); err != nil {
		return nil, fmt.Errorf("db: retrieve all authors failed: %v", err)
	}
	return dest, nil
}

func (s *Store) RetrieveAuthorWithID(id int) (*teal.Author, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest AuthorDest
	if err := SELECT(Authors.AllColumns, Books.AllColumns).
		FROM(Authors.
			INNER_JOIN(BooksAuthors, BooksAuthors.AuthorID.EQ(Authors.ID)).
			INNER_JOIN(Books, BooksAuthors.BookID.EQ(Books.ID))).
		WHERE(Authors.ID.EQ(Int(int64(id)))).Query(tx, &dest); err != nil {
		return nil, fmt.Errorf("db: retrieve author %d failed: %v", id, err)
	}
	res := modelToAuthor(dest)
	return res, nil
}

func (s *Store) RetrieveAuthorWithName(name string) (*teal.Author, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest AuthorDest
	if err := SELECT(Authors.AllColumns, Books.AllColumns).
		FROM(Authors.
			INNER_JOIN(BooksAuthors, BooksAuthors.AuthorID.EQ(Authors.ID)).
			INNER_JOIN(Books, BooksAuthors.BookID.EQ(Books.ID))).
		WHERE(Authors.Name.EQ(String(name))).Query(tx, &dest); err != nil {
		return nil, fmt.Errorf("db: retrieve author %q failed: %v", name, err)
	}
	res := modelToAuthor(dest)
	return res, nil
}

func (s *Store) RetrieveAllAuthors() ([]*teal.Author, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest []AuthorDest
	if err := SELECT(Authors.AllColumns, Books.AllColumns).
		FROM(Authors.
			INNER_JOIN(BooksAuthors, BooksAuthors.AuthorID.EQ(Authors.ID)).
			INNER_JOIN(Books, BooksAuthors.BookID.EQ(Books.ID))).
		Query(tx, &dest); err != nil {
		return nil, fmt.Errorf("db: retrieve all authors failed: %v", err)
	}

	var authors []*teal.Author
	for _, a := range dest {
		authors = append(authors, modelToAuthor(a))
	}
	return authors, nil
}

// Create one author entry without books
func (s *Store) CreateAuthor(a *teal.Author) error {
	if err := s.Tx(func(tx *sqlx.Tx) error {

		author := authorToModel(a)
		_, err := insertAuthor(tx, author.Authors)
		if err != nil {
			return err
		}
		return nil

	}, &sql.TxOptions{}); err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateAuthor(id int) error {
	return nil
}

func (s *Store) DeleteAuthor(id int) error {
	if err := s.Tx(func(tx *sqlx.Tx) error {
		err := deleteAuthor(tx, id)
		if err != nil {
			return err
		}

		// delete entry from booksAuthors table
		return nil

	}, &sql.TxOptions{}); err != nil {
		return err
	}
	return nil
}

func insertAuthor(tx *sqlx.Tx, a model.Authors) (int32, error) {

	var author model.Authors
	_, err := RawStatement(`INSERT OR IGNORE INTO authors (name) VALUES ($authorName);`,
		RawArgs{"$authorName": a.Name}).
		Exec(tx)
	if err != nil {
		return -1, fmt.Errorf("db: insert to authors table failed: %v", err)
	}

	if err := SELECT(Authors.AllColumns).
		FROM(Authors).
		WHERE(Authors.Name.EQ(String(a.Name))).
		Query(tx, &author); err != nil {
		return -1, fmt.Errorf("db: get last insert id from authors table failed: %v", err)
	}
	return author.ID, nil
}

func insertAuthors(tx *sqlx.Tx, a []model.Authors) ([]int32, error) {

	var ids []int32
	for _, author := range a {
		id, err := insertAuthor(tx, author)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func updateAuthor(tx *sqlx.Tx, id int) error {
	return nil
}

func deleteAuthor(tx *sqlx.Tx, id int) error {
	res, err := Authors.DELETE().WHERE(Authors.ID.EQ(Int(int64(id)))).Exec(tx)
	if err != nil {
		return fmt.Errorf("db: unable to delete author %d: %w", id, err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to delete author %d: %w", id, err)
	}
	if count == 0 {
		return errors.New("db: no authors removed")
	}
	return nil
}
