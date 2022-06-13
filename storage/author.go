package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/teal"
)

func parseAuthors(authors []string) []*teal.Author {
	var result []*teal.Author
	for _, a := range authors {
		result = append(result, &teal.Author{
			Name: a,
		})
	}
	return result
}

func (s *Store) RetrieveAuthorWithID(id int) (*teal.Author, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest teal.Author
	stmt := `SELECT * FROM authors WHERE id=$1;`
	err = tx.QueryRowx(stmt, id).StructScan(&dest)
	if err == sql.ErrNoRows {
		return nil, teal.ErrDoesNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve author %d failed: %v", id, err)
	}
	return &dest, nil
}

func (s *Store) RetrieveAuthorWithName(name string) (*teal.Author, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest teal.Author
	stmt := `SELECT * FROM authors WHERE name=$1;`
	err = tx.QueryRowx(stmt, name).StructScan(&dest)
	if err == sql.ErrNoRows {
		return nil, teal.ErrDoesNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve author %q failed: %v", name, err)
	}
	return &dest, nil
}

func (s *Store) RetrieveAllAuthors() ([]*teal.Author, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest []*teal.Author
	stmt := `SELECT * FROM authors;`
	err = tx.Select(&dest, stmt)
	if err == sql.ErrNoRows {
		return nil, teal.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve all authors failed: %v", err)
	}
	return dest, nil
}

func (s *Store) RetrieveAllAuthorNames() ([]string, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest []string
	stmt := `SELECT name FROM authors;`
	err = tx.Select(&dest, stmt)
	if err == sql.ErrNoRows {
		return nil, teal.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve all authors names failed: %v", err)
	}
	return dest, nil
}

func (s *Store) CreateAuthor(ctx context.Context, a *teal.Author) (*teal.Author, error) {
	if err := s.Tx(ctx, func(tx *sqlx.Tx) error {

		id, err := insertOrGetAuthor(tx, a)
		if err != nil {
			return err
		}
		// save id to context for querying later
		ctx = context.WithValue(ctx, authorCtxKey, id)
		return nil

	}, &sql.TxOptions{}); err != nil {
		return nil, err
	}

	// TODO implement separate context package with type safe getters and setters
	id, ok := ctx.Value(authorCtxKey).(int64)
	if !ok {
		return nil, fmt.Errorf("db: failed to cast author")
	}

	// query author after transaction committed
	author, err := s.RetrieveAuthorWithID(int(id))
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (s *Store) UpdateAuthor(ctx context.Context, id int, a *teal.Author) (*teal.Author, error) {
	if err := s.Tx(ctx, func(tx *sqlx.Tx) error {

		err := updateAuthor(tx, id, a)
		if err != nil {
			return err
		}
		return nil

	}, &sql.TxOptions{}); err != nil {
		return nil, err
	}

	author, err := s.RetrieveAuthorWithID(int(id))
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (s *Store) DeleteAuthor(ctx context.Context, id int) error {
	if err := s.Tx(ctx, func(tx *sqlx.Tx) error {

		err := deleteAuthor(tx, id)
		if err != nil {
			return err
		}

		// delete all author entries from booksAuthors table
		stmt := `DELETE FROM books_authors WHERE author_id=$1;`
		res, err := tx.Exec(stmt, id)
		if err != nil {
			return fmt.Errorf("db: delete author %d from books_authors failed: %v", id, err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("db: delete author %d from books_authors failed: %v", id, err)
		}

		if count == 0 {
			return errors.New("no rows deleted from books_authors table")
		}
		return nil

	}, &sql.TxOptions{}); err != nil {
		return err
	}
	return nil
}

// insert author. If already exists, return author id
func insertOrGetAuthor(tx *sqlx.Tx, a *teal.Author) (int64, error) {

	stmt := `INSERT OR IGNORE INTO authors (name) VALUES ($1);`
	res, err := tx.Exec(stmt, a.Name)
	if err != nil {
		return -1, fmt.Errorf("db: insert to authors table failed: %v", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("db: insert to authors table failed: %v", err)
	}

	// no rows inserted, query to get existing id
	if n == 0 {
		// authors.name is unique
		var id int64
		stmt := `SELECT id FROM authors WHERE name=$1;`
		err := tx.Get(&id, stmt, a.Name)
		if err != nil {
			return -1, fmt.Errorf("db: query existing author failed: %v", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("db: query existing author failed: %v", err)
		}
		return id, nil
	}
}

func insertOrGetAuthors(tx *sqlx.Tx, a []*teal.Author) ([]int64, error) {

	var ids []int64
	for _, author := range a {
		id, err := insertOrGetAuthor(tx, author)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func updateAuthor(tx *sqlx.Tx, id int, a *teal.Author) error {

	stmt := `UPDATE authors SET name=$1 WHERE id=$2`

	res, err := tx.Exec(stmt, a.Name, id)
	if err != nil {
		return fmt.Errorf("db: update author %d failed: %v", id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: update author %d failed: %v", id, err)
	}

	if count == 0 {
		return errors.New("db: no authors updated")
	}
	return nil
}

func deleteAuthor(tx *sqlx.Tx, id int) error {

	stmt := `DELETE FROM authors WHERE id=$1;`
	res, err := tx.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("db: unable to delete author %d: %w", id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to delete author %d: %w", id, err)
	}

	if count == 0 {
		return fmt.Errorf("db: author %d not removed", id)
	}
	return nil
}

func deleteAuthorsWithNoBooks(tx *sqlx.Tx) error {

	stmt := `DELETE FROM authors WHERE id NOT IN
				(SELECT author_id FROM books_authors);`
	res, err := tx.Exec(stmt)
	if err != nil {
		return fmt.Errorf("db: delete author from authors table failed: %v", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: delete author from authors table failed: %v", err)
	}

	if count != 0 {
		// TODO log author removed
		return nil
	}
	return nil
}
