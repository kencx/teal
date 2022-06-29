package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/teal"
)

type UserStore struct {
	db *sqlx.DB
}

func (s *UserStore) Get(id int64) (*teal.User, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var user teal.User
	stmt := `SELECT * FROM users WHERE id=$1;`
	err = tx.QueryRowx(stmt, id).StructScan(&user)
	if err == sql.ErrNoRows {
		return nil, teal.ErrDoesNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve user %d failed: %v", id, err)
	}
	return &user, nil
}

func (s *UserStore) GetByUsername(name string) (*teal.User, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var user teal.User
	stmt := `SELECT * FROM users WHERE username=$1;`
	err = tx.QueryRowx(stmt, name).StructScan(&user)
	if err == sql.ErrNoRows {
		return nil, teal.ErrDoesNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve user %q failed: %v", name, err)
	}
	return &user, nil
}

func (s *UserStore) GetAll() ([]*teal.User, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	var users []*teal.User
	stmt := `SELECT * FROM users;`
	err = tx.Select(&users, stmt)
	if err == sql.ErrNoRows {
		return nil, teal.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve all users failed: %v", err)
	}
	return users, nil
}

func (s *UserStore) Create(u *teal.User) (*teal.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	stmt := `INSERT INTO users
	(name, username, hashed_password, email, token, role, dateAdded)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	err = tx.QueryRowx(stmt,
		u.Name,
		u.Username,
		u.HashedPassword.Hash,
		u.Email,
		u.Token,
		u.Role,
		u.DateAdded).StructScan(u)

	// TODO check for duplicate username, email violation
	if err != nil {
		return nil, fmt.Errorf("db: insert to authors table failed: %v", err)
	}

	return u, nil
}

func (s *UserStore) Update(id int64, u *teal.User) (*teal.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	stmt := `UPDATE users
	SET name=$1,
	username=$2,
	hashed_password=$3,
	email=$4,
	token=$5,
	role=$6,
	dateAdded=$7
	WHERE id=$8`

	res, err := tx.Exec(stmt,
		u.Name,
		u.Username,
		u.HashedPassword.Hash,
		u.Email,
		u.Token,
		u.Role,
		u.DateAdded,
		id)

	if err != nil {
		return nil, fmt.Errorf("db: update user %d failed: %v", id, err)
	}
	// TODO check for duplicate username, email violation

	count, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("db: update user %d failed: %v", id, err)
	}
	if count == 0 {
		return nil, errors.New("db: no users updated")
	}
	return u, nil
}

func (s *UserStore) Delete(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

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
	if err != nil {
		return err
	}
	return nil
}
