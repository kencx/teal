package storage

import (
	"reflect"
	"testing"

	"github.com/kencx/teal/pkg"
)

func setup(t *testing.T) *DB {
	db := NewDB("sqlite3")
	if err := db.Open("./test.db"); err != nil {
		t.Fatal(err)
	}

	if err := db.createTable(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.dropTable()
		db.Close()
	})
	return db
}

func TestGetBook(t *testing.T) {
	db := setup(t)

	// populate db
	expected := &pkg.Book{
		Title:  "FooBar",
		Author: "John Doe",
		ISBN:   "12345",
	}
	id, err := db.CreateBook(expected)
	checkErr(t, err)

	result, err := db.GetBook(id)
	checkErr(t, err)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestGetBookNotExists(t *testing.T) {
	db := setup(t)

	expected := &pkg.Book{
		Title:  "FooBar",
		Author: "John Doe",
		ISBN:   "12345",
	}
	_, err := db.CreateBook(expected)
	checkErr(t, err)

	res, err := db.GetBook(2)
	if res != nil {
		t.Fatalf("got %v, want nil", res)
	}
}

func TestGetBookByTitle(t *testing.T) {
	db := setup(t)

	// populate db
	expected := &pkg.Book{
		Title:  "FooBar",
		Author: "John Doe",
		ISBN:   "12345",
	}
	_, err := db.CreateBook(expected)
	checkErr(t, err)

	result, err := db.GetBookByTitle("FooBar")
	checkErr(t, err)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestGetAllBooks(t *testing.T) {
	db := setup(t)

	expected := []*pkg.Book{
		&pkg.Book{
			Title:  "FooBar",
			Author: "John Doe",
			ISBN:   "12345",
		},
		&pkg.Book{
			Title:  "BarBaz",
			Author: "Ben Adams",
			ISBN:   "54678",
		},
	}
	for _, b := range expected {
		_, err := db.CreateBook(b)
		checkErr(t, err)
	}

	result, err := db.GetAllBooks()
	checkErr(t, err)

	if len(result) != len(expected) {
		t.Errorf("got %v, want %v", result, expected)
	}

	// check if all elems in slice match
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestCreateBook(t *testing.T) {
	db := setup(t)

	expected := &pkg.Book{
		Title:  "FooBar",
		Author: "John Doe",
		ISBN:   "12345",
	}
	_, err := db.CreateBook(expected)
	checkErr(t, err)

	books, err := db.GetAllBooks()
	checkErr(t, err)

	if len(books) != 1 {
		t.Errorf("unexpected number of books: %d", len(books))
	}
	if !reflect.DeepEqual(books[0], expected) {
		t.Errorf("got %v, want %v", books[0], expected)
	}
}

func TestCreateNotUniqueBook(t *testing.T) {
	db := setup(t)

	expected := &pkg.Book{
		Title:  "FooBar",
		Author: "John Doe",
		ISBN:   "12345",
	}
	_, err := db.CreateBook(expected)
	checkErr(t, err)

	_, err = db.CreateBook(expected)
}

func TestUpdateBook(t *testing.T) {
	db := setup(t)

	b := &pkg.Book{
		Title:  "FooBar",
		Author: "John Doe",
		ISBN:   "12345",
	}
	id, err := db.CreateBook(b)
	checkErr(t, err)

	expected := &pkg.Book{
		Title:  "BarBaz",
		Author: "Ben Adams",
		ISBN:   "1022",
	}

	err = db.UpdateBook(id, expected)
	checkErr(t, err)

	result, err := db.GetBook(id)
	checkErr(t, err)

	if result.Title != expected.Title {
		t.Errorf("got %v, want %v", result.Title, expected.Title)
	}
	if result.Author != expected.Author {
		t.Errorf("got %v, want %v", result.Author, expected.Author)
	}
	if result.ISBN != expected.ISBN {
		t.Errorf("got %v, want %v", result.ISBN, expected.ISBN)
	}
}

func TestDeleteBook(t *testing.T) {
	db := setup(t)

	b := &pkg.Book{
		Title:  "FooBar",
		Author: "John Doe",
		ISBN:   "12345",
	}
	id, err := db.CreateBook(b)
	checkErr(t, err)

	err = db.DeleteBook(id)
	checkErr(t, err)

	res, err := db.GetBook(id)
	if res != nil {
		t.Fatalf("got %v, want nil", res)
	}
}

func checkErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
