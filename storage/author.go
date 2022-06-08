package storage

import (
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

func (s *Store) RetrieveAllAuthorNames() ([]string, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest []string
	stmt := `SELECT name FROM authors`
	if err = tx.Select(&dest, stmt); err != nil {
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

	var dest teal.Author
	stmt := `SELECT * FROM authors WHERE id=$1`
	if err = tx.QueryRowx(stmt, id).StructScan(&dest); err != nil {
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
	stmt := `SELECT * FROM authors WHERE name=$1`
	if err = tx.QueryRowx(stmt, name).StructScan(&dest); err != nil {
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
	stmt := `SELECT * FROM authors`
	if err = tx.Select(&dest, stmt); err != nil {
		return nil, fmt.Errorf("db: retrieve all authors failed: %v", err)
	}
	return dest, nil
}

func (s *Store) CreateAuthor(a *teal.Author) error {
	if err := s.Tx(func(tx *sqlx.Tx) error {

		_, err := insertAuthor(tx, a)
		if err != nil {
			return err
		}
		return nil

	}, &sql.TxOptions{}); err != nil {
		return err
	}
	return nil
}

// func (s *Store) UpdateAuthor(id int) error {
// 	return nil
// }
//
// func (s *Store) DeleteAuthor(id int) error {
// 	if err := s.Tx(func(tx *sqlx.Tx) error {
// 		err := deleteAuthor(tx, id)
// 		if err != nil {
// 			return err
// 		}
//
// 		// delete entry from booksAuthors table
// 		return nil
//
// 	}, &sql.TxOptions{}); err != nil {
// 		return err
// 	}
// 	return nil
// }

func insertAuthor(tx *sqlx.Tx, a *teal.Author) (int64, error) {

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
		stmt := `SELECT id FROM authors WHERE name=$1`
		err := tx.Get(&id, stmt, a.Name)
		if err != nil {
			return -1, fmt.Errorf("db: get last insert id from authors table failed: %v", err)
		}
		return id, nil

	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("db: get last insert id from authors table failed: %v", err)
		}
		return id, nil
	}
}

func insertAuthors(tx *sqlx.Tx, a []*teal.Author) ([]int64, error) {

	var ids []int64
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

	stmt := `DELETE FROM authors WHERE id=$1`
	res, err := tx.Exec(stmt, id)
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
