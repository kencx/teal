package storage

import (
	"os"
	"testing"
)

func TestOpen(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		db := NewDB("sqlite3")
		if err := db.Open("./test.db"); err != nil {
			t.Error(err)
		}
		defer db.Close()
	})

	t.Run("no driver", func(t *testing.T) {
		db := NewDB("")
		if err := db.Open("./test.db"); err == nil {
			t.Error("expected error: driver required")
		}
	})

	t.Run("no DSN", func(t *testing.T) {
		db := NewDB("sqlite3")
		if err := db.Open(""); err == nil {
			t.Error("expected error: connection string required")
		}
	})

	os.Remove("./test.db")
}
