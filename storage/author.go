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

func authorToModel(a *teal.Author) *model.Authors {
	return nil
}

// Create one author entry without books
func (s *Store) CreateAuthor(a *teal.Author) error {
	if err := s.Tx(func(tx *sqlx.Tx) error {

		author := authorToModel(a)
		_, err := insertAuthor(tx, *author)
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

func insertAuthor(tx *sqlx.Tx, a model.Authors) (*int32, error) {

	var author model.Authors
	err := Authors.INSERT(Authors.MutableColumns).MODEL(a).RETURNING(Authors.ID).Query(tx, &author)
	if err != nil {
		return nil, fmt.Errorf("db: insert to authors table failed: %v", err)
	}

	return author.ID, nil
}

func insertAuthors(tx *sqlx.Tx, a []model.Authors) ([]*int32, error) {

	var authors []model.Authors
	err := Authors.INSERT(Authors.MutableColumns).MODELS(a).RETURNING(Authors.ID).Query(tx, &authors)
	if err != nil {
		return nil, fmt.Errorf("db: insert to authors table failed: %v", err)
	}

	var ids []*int32
	for _, v := range authors {
		ids = append(ids, v.ID)
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
