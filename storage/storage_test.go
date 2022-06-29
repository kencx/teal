package storage

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

// globals for tests only
var (
	testDSN        = "./test.db"
	testSchemaPath = "../migrations/schema.sql"
	testDataPath   = "../migrations/testdata.sql"
	dropSchemaPath = "../migrations/dropall.sql"
	testdb         = setup()
	ts             = NewStore(testdb)
)

func TestMain(m *testing.M) {
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() *sqlx.DB {
	testdb, err := Open(testDSN)
	if err != nil {
		log.Fatal(err)
	}

	bootstrap(testdb)
	return testdb
}

func teardown() {
	dropTables(testdb)
	Close(testdb)
	os.Remove(testDSN)
}

func bootstrap(db *sqlx.DB) {
	err := ExecFile(db, testSchemaPath)
	if err != nil {
		log.Fatalf("initTestDB: %v", err)
	}
	err = ExecFile(db, testDataPath)
	if err != nil {
		log.Fatalf("initTestDB: %v", err)
	}
}

func dropTables(db *sqlx.DB) {
	if err := ExecFile(db, dropSchemaPath); err != nil {
		log.Fatal(err)
	}
}

func resetDB(db *sqlx.DB) {
	dropTables(db)
	bootstrap(db)
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

func checkErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

// pretty prints structs for readability
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func contains(s []string, a string) bool {
	for _, b := range s {
		if a == b {
			return true
		}
	}
	return false
}
