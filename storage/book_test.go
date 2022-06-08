package storage

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	// . "github.com/go-jet/jet/v2/sqlite"
	"github.com/kencx/teal"
	// . "github.com/kencx/teal/storage/sqlite/table"
)

// TODO replace with TestMain to speed up tests by
// running setup only once

func setup(t *testing.T) *Store {
	db := NewStore("sqlite3")
	err := db.Open("./test.db")
	checkErr(t, err)

	err = db.ExecFile("./testdata.sql")
	checkErr(t, err)

	t.Cleanup(func() {
		db.dropTable()
		db.Close()
		os.Remove("./test.db")
	})
	return db
}

func TestCreateBook(t *testing.T) {
	db := setup(t)

	// TODO implement table driven tests
	want := &teal.Book{
		Title:      "World War Z",
		ISBN:       "45678",
		Author:     teal.Authors{{Name: "Max Brooks"}, {Name: "Second Author"}},
		NumOfPages: 100,
		Rating:     10,
		State:      "read",
	}
	err := db.CreateBook(want)
	checkErr(t, err)

	got, err := db.RetrieveBookWithTitle(want.Title)
	checkErr(t, err)

	// TODO implement compare
	if !compare(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

// compare books without ids
func compare(b1, b2 *teal.Book) bool {
	return false
}

func TestCreateBookNotUnique(t *testing.T) {
	db := setup(t)

	err := db.CreateBook(testBook1)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestRetrieveBookWithID(t *testing.T) {
	db := setup(t)

	got, err := db.RetrieveBookWithID(testBook1.ID)
	checkErr(t, err)

	want := testBook1
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestRetrieveBookWithTitle(t *testing.T) {
	db := setup(t)

	got, err := db.RetrieveBookWithTitle(testBook2.Title)
	checkErr(t, err)

	want := testBook2
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

// func TestDelete(t *testing.T) {
// 	db := setup(t)
//
// 	err := deleteBook(db, 1)
// 	checkErr(t, err)
//
// 	got, err := getBook(db, 1)
// 	checkErr(t, err)
// 	if got != nil {
// 		t.Fatalf("got %v, want nil", prettyPrint(got))
// 	}
// }
//
// func TestDeleteNotExists(t *testing.T) {
// 	db := setup(t)
//
// 	err := deleteBook(db, 999)
// 	if err == nil {
// 		t.Fatalf("expected error: book not exists")
// 	}
// }

// func TestGetNotExists(t *testing.T) {
// 	db := setup(t)
//
// 	result, err := db.GetBook(999)
// 	checkErr(t, err)
// 	if result != nil {
// 		t.Fatalf("got %v, want nil", result)
// 	}
// }
//
// func TestGetAll(t *testing.T) {
// 	db := setup(t)
// 	expected := []*teal.Book{testdata["book1"].(*teal.Book), testdata["book2"].(*teal.Book)}
//
// 	result, err := db.GetAllBooks()
// 	checkErr(t, err)
//
// 	if len(result) != len(expected) {
// 		t.Errorf("got %v, want %v", result, expected)
// 	}
//
// 	// check if all elems in slice match
// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("got %v, want %v", result, expected)
// 	}
// }
//
// func TestUpdate(t *testing.T) {
// 	db := setup(t)
// 	id := testdata["book1"].(*teal.Book).ID
//
// 	expected := &teal.Book{
// 		Title:  "The Linux Command Line",
// 		Author: "William Shotts",
// 		ISBN:   "59789",
// 	}
//
// 	err := db.UpdateBook(id, expected)
// 	checkErr(t, err)
//
// 	result, err := db.GetBook(id)
// 	checkErr(t, err)
//
// 	if result.Title != expected.Title {
// 		t.Errorf("got %v, want %v", result.Title, expected.Title)
// 	}
// 	if result.Author != expected.Author {
// 		t.Errorf("got %v, want %v", result.Author, expected.Author)
// 	}
// 	if result.ISBN != expected.ISBN {
// 		t.Errorf("got %v, want %v", result.ISBN, expected.ISBN)
// 	}
// }
//
// func TestUpdateNotExists(t *testing.T) {
// 	db := setup(t)
//
// 	b := &teal.Book{
// 		Title:  "The Linux Command Line",
// 		Author: "William Shotts",
// 		ISBN:   "59789",
// 	}
// 	err := db.UpdateBook(999, b)
// 	if err == nil {
// 		t.Fatalf("expected error: book not exists")
// 	}
// }
//
// func TestUpdateISBNConstraint(t *testing.T) {
// 	db := setup(t)
// 	id := testdata["book2"].(*teal.Book).ID
//
// 	b := &teal.Book{
// 		Title:  "The Linux Command Line",
// 		Author: "William Shotts",
// 		ISBN:   "12345",
// 	}
// 	err := db.UpdateBook(id, b)
// 	if err == nil {
// 		t.Errorf("expected error")
// 	}
// }

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
