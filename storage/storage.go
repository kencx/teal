package storage

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	SQLITE   = "sqlite3"
	POSTGRES = "postgres"
)

type Store struct {
	Books   *BookStore
	Authors *AuthorStore
	Users   *UserStore
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		Books:   &BookStore{db},
		Authors: &AuthorStore{db},
		Users:   &UserStore{db},
	}
}

func (s *Store) GetDB() *sqlx.DB {
	return s.Books.db
}

func Open(path string) (*sqlx.DB, error) {
	if path == "" {
		return nil, fmt.Errorf("db: connection string required")
	}

	db, err := sqlx.Open(SQLITE, path)
	if err != nil {
		return nil, fmt.Errorf("db: failed to open: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db: failed to connect: %w", err)
	}
	if err := createTable(db); err != nil {
		return nil, err
	}
	return db, nil
}

func Close(db *sqlx.DB) error {
	if db != nil {
		return db.Close()
	}
	return fmt.Errorf("db: db is nil")
}

// TODO replace with migration file
func createTable(db *sqlx.DB) error {
	if err := ExecFile(db, "../migrations/schema.sql"); err != nil {
		return fmt.Errorf("%v: %s", err, CREATE_TABLES)
	}
	return nil
}

// TODO replace with migration file
func dropTables(db *sqlx.DB) error {
	if err := ExecFile(db, "../migrations/dropall.sql"); err != nil {
		return fmt.Errorf("%v: %s", err, DROP_ALL)
	}
	return nil
}

func ExecFile(db *sqlx.DB, filePath string) error {
	query, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("db: cannot read sql file %q: %v", filePath, err)
	}

	if _, err := db.Exec(string(query)); err != nil {
		return fmt.Errorf("%v", err)
	}
	log.Printf("File %s loaded", filePath)
	return nil
}
