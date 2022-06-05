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

type Repository struct {
	db     *sql.DB
	Driver string
	DSN    string
}

func NewRepository(driver string) *Repository {
	return &Repository{
		Driver: driver,
	}
}

func (r *Repository) Open(path string) error {
	if r.Driver == "" {
		return fmt.Errorf("database driver required")
	}

	if r.Driver != SQLITE && r.Driver != POSTGRES {
		return fmt.Errorf("database driver not supported")
	}

	r.DSN = path
	if r.DSN == "" {
		return fmt.Errorf("database connection string required")
	}

	var err error
	if r.db, err = sql.Open(r.Driver, r.DSN); err != nil {
		return fmt.Errorf("[ERROR] Failed to open database: %w", err)
	}
	if err = r.db.Ping(); err != nil {
		return fmt.Errorf("[ERROR] Failed to connect to database: %w", err)
	}
	if err := r.createTable(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

func (r *Repository) createTable() error {
	_, err := r.db.Exec(CREATE_STMT)
	if err != nil {
		return fmt.Errorf("%q: %s", err, CREATE_STMT)
	}
	return nil
}

func (r *Repository) dropTable() error {
	if _, err := r.db.Exec(DROP_STMT); err != nil {
		return fmt.Errorf("%q: %s", err, DROP_STMT)
	}
	return nil
}
