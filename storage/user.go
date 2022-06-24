package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/teal"
)

type UserStore struct {
	db *sqlx.DB
}

func (s *UserStore) Get(id int) (*teal.User, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest teal.User
	stmt := `SELECT * FROM users WHERE id=$1;`
	err = tx.QueryRowx(stmt, id).StructScan(&dest)
	if err == sql.ErrNoRows {
		return nil, teal.ErrDoesNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve user %d failed: %v", id, err)
	}
	return &dest, nil
}

func (s *UserStore) GetByUsername(name string) (*teal.User, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest teal.User
	stmt := `SELECT * FROM users WHERE username=$1;`
	err = tx.QueryRowx(stmt, name).StructScan(&dest)
	if err == sql.ErrNoRows {
		return nil, teal.ErrDoesNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve user %q failed: %v", name, err)
	}
	return &dest, nil
}

func (s *UserStore) GetAll() ([]*teal.User, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var dest []*teal.User
	stmt := `SELECT * FROM users;`
	err = tx.Select(&dest, stmt)
	if err == sql.ErrNoRows {
		return nil, teal.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve all users failed: %v", err)
	}
	return dest, nil
}

// func (s *UserStore) Create(ctx context.Context, a *teal.User) (*teal.User, error) {
// 	if err := Tx(s.db, ctx, func(tx *sqlx.Tx) error {
//
// 		id, err := insertOrGetUser(tx, a)
// 		if err != nil {
// 			return err
// 		}
// 		// save id to context for querying later
// 		ctx = tcontext.WithUser(ctx, id)
// 		return nil
//
// 	}, &sql.TxOptions{}); err != nil {
// 		return nil, err
// 	}
//
// 	// TODO implement separate context package with type safe getters and setters
// 	id, err := tcontext.GetUser(ctx)
//
// 	// query user after transaction committed
// 	user, err := s.Get(int(id))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

func (s *UserStore) Update(ctx context.Context, id int, a *teal.User) (*teal.User, error) {
	if err := Tx(s.db, ctx, func(tx *sqlx.Tx) error {

		err := updateUser(tx, id, a)
		if err != nil {
			return err
		}
		return nil

	}, &sql.TxOptions{}); err != nil {
		return nil, err
	}

	user, err := s.Get(int(id))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserStore) Delete(ctx context.Context, id int) error {
	if err := Tx(s.db, ctx, func(tx *sqlx.Tx) error {

		err := deleteUser(tx, id)
		if err != nil {
			return err
		}
		return nil

	}, &sql.TxOptions{}); err != nil {
		return err
	}
	return nil
}

// insert user. If already exists, return author id
func insertOrGetUser(tx *sqlx.Tx, a *teal.User) (int64, error) {

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

func insertOrGetUsers(tx *sqlx.Tx, a []*teal.User) ([]int64, error) {

	var ids []int64
	for _, user := range a {
		id, err := insertOrGetUser(tx, user)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func updateUser(tx *sqlx.Tx, id int, a *teal.User) error {

	stmt := `UPDATE users SET name=$1 WHERE id=$2`

	res, err := tx.Exec(stmt, a.Name, id)
	if err != nil {
		return fmt.Errorf("db: update user %d failed: %v", id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: update user %d failed: %v", id, err)
	}

	if count == 0 {
		return errors.New("db: no users updated")
	}
	return nil
}

func deleteUser(tx *sqlx.Tx, id int) error {

	stmt := `DELETE FROM users WHERE id=$1;`
	res, err := tx.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("db: unable to delete user %d: %w", id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: unable to delete user %d: %w", id, err)
	}

	if count == 0 {
		return fmt.Errorf("db: user %d not removed", id)
	}
	return nil
}
