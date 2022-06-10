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
	db     *sqlx.DB
	Driver string
	DSN    string
}

func NewStore(driver string) *Store {
	return &Store{
		Driver: driver,
	}
}

func (r *Store) Open(path string) error {
	if r.Driver == "" {
		return fmt.Errorf("db: driver required")
	}

	if r.Driver != SQLITE && r.Driver != POSTGRES {
		return fmt.Errorf("db: driver not supported")
	}

	r.DSN = path
	if r.DSN == "" {
		return fmt.Errorf("db: connection string required")
	}

	var err error
	if r.db, err = sqlx.Open(r.Driver, r.DSN); err != nil {
		return fmt.Errorf("db: failed to open: %w", err)
	}
	if err = r.db.Ping(); err != nil {
		return fmt.Errorf("db: failed to connect: %w", err)
	}
	if err := r.createTable(); err != nil {
		return err
	}
	return nil
}

func (r *Store) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

func (r *Store) createTable() error {
	_, err := r.db.Exec(CREATE_TABLES)
	if err != nil {
		return fmt.Errorf("%v: %s", err, CREATE_TABLES)
	}
	return nil
}

func (r *Store) dropTable() error {

	if _, err := r.db.Exec(DROP_ALL); err != nil {
		return fmt.Errorf("%v: %s", err, DROP_ALL)
	}
	return nil
}

func (r *Store) ExecFile(path string) error {
	query, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("db: cannot read sql file %q: %v", path, err)
	}

	if _, err := r.db.Exec(string(query)); err != nil {
		return fmt.Errorf("%v", err)
	}
	log.Printf("File %s loaded", path)
	return nil
}
