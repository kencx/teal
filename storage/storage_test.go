package storage

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

var (
	db = setup()
	ts = NewStore(db)
)

func TestMain(m *testing.M) {
	defer func() {
		dropTable(db)
		Close(db)
		os.Remove("./test.db")
	}()
	os.Exit(m.Run())
}

func setup() *sqlx.DB {
	db, err := Open("./test.db")
	if err != nil {
		log.Fatal(err)
	}

	initTestDB(db)
	return db
}

func initTestDB(db *sqlx.DB) {
	err := ExecFile(db, "./schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	err = ExecFile(db, "./testdata.sql")
	if err != nil {
		log.Fatal(err)
	}
}

func resetDB(db *sqlx.DB) {
	initTestDB(db)
}

func TestOpen(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		db, err := Open("./test.db")
		if err != nil {
			t.Error(err)
		}
		defer db.Close()
	})

	t.Run("no DSN", func(t *testing.T) {
		_, err := Open("")
		if err == nil {
			t.Error("expected error: connection string required")
		}
	})
}
