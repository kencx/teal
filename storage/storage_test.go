package storage

import (
	"log"
	"os"
	"testing"
)

var db = setup()

func TestMain(m *testing.M) {
	defer func() {
		db.dropTable()
		db.Close()
		os.Remove("./test.db")
	}()
	os.Exit(m.Run())
}

func setup() *Store {
	db := NewStore("sqlite3")
	err := db.Open("./test.db")
	if err != nil {
		log.Fatal(err)
	}

	initTestDB(db)
	return db
}

func initTestDB(db *Store) {
	err := db.ExecFile("./schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	err = db.ExecFile("./testdata.sql")
	if err != nil {
		log.Fatal(err)
	}
}

func TestOpen(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		db := NewStore("sqlite3")
		if err := db.Open("./test.db"); err != nil {
			t.Error(err)
		}
		defer db.Close()
	})

	t.Run("no driver", func(t *testing.T) {
		db := NewStore("")
		if err := db.Open("./test.db"); err == nil {
			t.Error("expected error: driver required")
		}
	})

	t.Run("no DSN", func(t *testing.T) {
		db := NewStore("sqlite3")
		if err := db.Open(""); err == nil {
			t.Error("expected error: connection string required")
		}
	})

	os.Remove("./test.db")
}
