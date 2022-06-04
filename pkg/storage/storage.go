package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	SQLITE      = "sqlite3"
	POSTGRES    = "postgres"
	CREATE_STMT = `CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		isbn TEXT NOT NULL UNIQUE
	);`
	DROP_STMT = `DROP TABLE IF EXISTS books;`
)

type DB struct {
	db     *sql.DB
	Driver string
	DSN    string
}

func NewDB(driver string) *DB {
	return &DB{
		Driver: driver,
	}
}

func (db *DB) Open(path string) error {
	if db.Driver == "" {
		return fmt.Errorf("Database driver required")
	}

	db.DSN = path
	if db.DSN == "" {
		return fmt.Errorf("Database connection string required")
	}

	var err error
	if db.db, err = sql.Open(db.Driver, db.DSN); err != nil {
		return fmt.Errorf("[ERROR] Failed to open database: %w", err)
	}
	if err = db.db.Ping(); err != nil {
		return fmt.Errorf("[ERROR] Failed to connect to database: %w", err)
	}
	if err := db.createTable(); err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

func (db *DB) createTable() error {
	_, err := db.db.Exec(CREATE_STMT)
	if err != nil {
		return fmt.Errorf("%q: %s\n", err, CREATE_STMT)
	}
	return nil
}

func (db *DB) dropTable() error {
	if _, err := db.db.Exec(DROP_STMT); err != nil {
		return fmt.Errorf("%q: %s\n", err, DROP_STMT)
	}
	return nil
}
