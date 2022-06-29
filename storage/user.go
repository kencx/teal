package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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
	stmt := `SELECT id, name, username, hashed_password, role, dateAdded
	FROM users WHERE id=$1;`
	err = tx.QueryRowx(stmt, id).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.HashedPassword.Hash,
		&user.Role,
		&user.DateAdded,
	)
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
	stmt := `SELECT id, name, username, hashed_password, role, dateAdded
	FROM users WHERE username=$1;`
	err = tx.QueryRowx(stmt, name).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.HashedPassword.Hash,
		&user.Role,
		&user.DateAdded,
	)
	if err == sql.ErrNoRows {
		return nil, teal.ErrDoesNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("db: retrieve user %q failed: %v", name, err)
	}
	return &user, nil
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
	(name, username, hashed_password, role, dateAdded)
	VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	err = tx.QueryRowx(stmt,
		u.Name,
		u.Username,
		u.HashedPassword.Hash,
		u.Role,
		u.DateAdded).StructScan(u)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, teal.ErrDuplicateUsername
		} else {
			return nil, fmt.Errorf("db: insert to users table failed: %v", err)
		}
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
	role=$4,
	dateAdded=$5
	WHERE id=$6`

	res, err := tx.Exec(stmt,
		u.Name,
		u.Username,
		u.HashedPassword.Hash,
		u.Role,
		u.DateAdded,
		id)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, teal.ErrDuplicateUsername
		} else {
			return nil, fmt.Errorf("db: update user %d failed: %v", id, err)
		}
	}

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
