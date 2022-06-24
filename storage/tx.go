package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// functional Tx helper for Exec statements
func Tx(db *sqlx.DB, ctx context.Context, fn func(tx *sqlx.Tx) error, opts *sql.TxOptions) error {
	tx, err := db.BeginTxx(ctx, opts)
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
