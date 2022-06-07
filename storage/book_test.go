package storage

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/kencx/teal"
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

func TestGetBook(t *testing.T) {
	db := setup(t)

	got, err := getBook(db, 1)
	checkErr(t, err)

	want := testBook1

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestInsertBook(t *testing.T) {
	db := setup(t)

	want := &teal.Book{
		Title:      "World War Z",
		ISBN:       "45678",
		NumOfPages: 100,
		Rating:     10,
		State:      "read",
	}
	id, err := insertBook(db, want)
	checkErr(t, err)

	got, err := getBook(db, id)
	checkErr(t, err)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestInsertTwoBooks(t *testing.T) {
	db := setup(t)

	want := []teal.Book{
		{
			Title:      "Persopolis Rising",
			ISBN:       "500",
			NumOfPages: 256,
			Rating:     8,
			State:      "unread",
		},
		{
			Title:      "Dune",
			ISBN:       "501",
			NumOfPages: 1500,
			Rating:     8,
			State:      "read",
		},
	}

	var got []teal.Book
	for i, b := range want {
		id, err := insertBook(db, &b)
		checkErr(t, err)

		want[i].ID = id // to ensure ID not 0

		res, err := getBook(db, id)
		checkErr(t, err)

		got = append(got, *res)
	}

	if len(got) != len(want) {
		t.Errorf("got %v, want %v", len(got), len(want))
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestInsertNotUnique(t *testing.T) {
	db := setup(t)

	_, err := insertBook(db, testBook1)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestDelete(t *testing.T) {
	db := setup(t)

	err := deleteBook(db, 1)
	checkErr(t, err)

	got, err := getBook(db, 1)
	checkErr(t, err)
	if got != nil {
		t.Fatalf("got %v, want nil", prettyPrint(got))
	}
}

func TestDeleteNotExists(t *testing.T) {
	db := setup(t)

	err := deleteBook(db, 999)
	if err == nil {
		t.Fatalf("expected error: book not exists")
	}
}

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
// func TestGetByTitle(t *testing.T) {
// 	db := setup(t)
// 	expected := testdata["book1"].(*teal.Book)
//
// 	result, err := db.GetBookByTitle(expected.Title)
// 	checkErr(t, err)
//
// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("got %v, want %v", prettyPrint(result), prettyPrint(expected))
// 	}
// }
//
// func TestGetByTitleNotExists(t *testing.T) {
// 	db := setup(t)
//
// 	result, err := db.GetBookByTitle("asdfgh")
// 	checkErr(t, err)
// 	if result != nil {
// 		t.Fatalf("got %v, want nil", prettyPrint(result))
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
// func TestCreate(t *testing.T) {
// 	db := setup(t)
//
// 	a := &teal.Authors{
// 		{Name: "S.A. Corey"},
// 	}
//
// 	expected := &teal.Book{
// 		Title:  "Leviathan Wakes",
// 		Author: *a,
// 		ISBN:   "99999",
// 		// DateAdded:   time.Now().Format(time.RFC822),
// 		// DateUpdated: time.Now().Format(time.RFC822),
// 	}
// 	err := db.CreateBook(expected)
// 	checkErr(t, err)
//
// 	// check books table
// 	// check authors table
// 	// check books_authors table
//
// 	result, err := db.RetrieveBook(1)
// 	checkErr(t, err)
//
// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("got %v, want %v", prettyPrint(result), prettyPrint(expected))
// 	}
// }
//
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
