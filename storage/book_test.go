package storage

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"sort"
	"testing"

	// . "github.com/go-jet/jet/v2/sqlite"
	"github.com/kencx/teal"
	// . "github.com/kencx/teal/storage/sqlite/table"
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

	err = db.ExecFile("./schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	// init test data
	err = db.ExecFile("./testdata.sql")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func TestRetrieveBookWithID(t *testing.T) {
	got, err := db.RetrieveBookWithID(testBook1.ID)
	checkErr(t, err)

	want := testBook1
	if !compareBook(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestRetrieveBookWithISBN(t *testing.T) {
	got, err := db.RetrieveBookWithISBN(testBook1.ISBN)
	checkErr(t, err)

	want := testBook1
	if !compareBook(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestRetrieveBookWithTitle(t *testing.T) {
	got, err := db.RetrieveBookWithTitle(testBook2.Title)
	checkErr(t, err)

	want := testBook2
	if !compareBook(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestRetrieveBookNotExists(t *testing.T) {
	result, err := db.RetrieveBookWithID(-1)
	if err == nil {
		t.Fatalf("expected error")
	}

	if result != nil {
		t.Fatalf("got %v, want nil", result)
	}
}

func TestRetrieveAllBooks(t *testing.T) {
	got, err := db.RetrieveAllBooks()
	checkErr(t, err)

	sort.Slice(got, func(i, j int) bool {
		return got[i].ID < got[j].ID
	})

	want := []*teal.Book{testBook1, testBook2, testBook3}

	if len(got) != len(want) {
		t.Fatalf("got %d books, want %d books", len(got), len(want))
	}

	for i := 0; i < len(got); i++ {
		if !compareBook(got[i], want[i]) {
			t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
		}
	}
}

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name string
		want *teal.Book
	}{{
		name: "book with minimal data",
		want: &teal.Book{
			Title:  "1984",
			ISBN:   "1001",
			Author: []string{"George Orwell"},
		},
	}, {
		name: "book with all data",
		want: &teal.Book{
			Title:      "World War Z",
			ISBN:       "1002",
			Author:     []string{"Max Brooks"},
			NumOfPages: 100,
			Rating:     10,
			State:      "read",
		},
	}, {
		name: "book with two authors",
		want: &teal.Book{
			Title:      "Pro Git",
			ISBN:       "1003",
			Author:     []string{"Scott Chacon", "Ben Straub"},
			NumOfPages: 100,
			Rating:     10,
			State:      "read",
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.CreateBook(tt.want)
			checkErr(t, err)

			got, err := db.RetrieveBookWithISBN(tt.want.ISBN)
			checkErr(t, err)

			if !compareBook(got, tt.want) {
				t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(tt.want))
			}
		})
	}
}

func TestCreateBookExistingISBN(t *testing.T) {
	err := db.CreateBook(testBook2)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestCreateBookExistingAuthor(t *testing.T) {
	want := &teal.Book{
		Title:      "New Book",
		ISBN:       "1004",
		Author:     []string{"John Doe"},
		NumOfPages: 100,
		Rating:     10,
		State:      "unread",
	}

	err := db.CreateBook(want)
	checkErr(t, err)

	tx, err := db.db.Beginx()
	if err != nil {
		t.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	// check for number of entries in authors
	var dest []string

	stmt := `SELECT name FROM authors WHERE name=$1`
	if err := tx.Select(&dest, stmt, want.Author[0]); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(dest) != 1 {
		t.Error("more than one author inserted")
	}
}

func TestCreateBookNewAndExistingAuthor(t *testing.T) {
	want := &teal.Book{
		Title:      "Thinking Fast and Slow",
		ISBN:       "1005",
		Author:     []string{"John Doe", "Newest Author"},
		NumOfPages: 100,
		Rating:     10,
		State:      "unread",
	}

	err := db.CreateBook(want)
	checkErr(t, err)

	tx, err := db.db.Beginx()
	if err != nil {
		t.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	// check for number of entries in authors
	var dest []string

	stmt := `SELECT name FROM authors WHERE name=$1 OR name=$2`
	if err := tx.Select(&dest, stmt, want.Author[0], want.Author[1]); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(dest) != 2 {
		t.Error("more than one author inserted")
	}
}

// func TestUpdate(t *testing.T) {
// 	db := setup()
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
// 	db := setup()
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
// 	db := setup()
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

func TestDeleteBook(t *testing.T) {

	err := db.DeleteBook(testBook1.ID)
	checkErr(t, err)

	_, err = db.RetrieveBookWithID(testBook1.ID)
	if err == nil {
		t.Errorf("expected error, book %d not deleted", testBook1.ID)
	}

	tx, err := db.db.Beginx()
	if err != nil {
		t.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	// check entries deleted from books_authors
	var dest []int
	stmt := `SELECT author_id FROM books_authors WHERE books_authors.book_id=$1`
	if err := tx.Select(&dest, stmt, testBook1.ID); err != nil {
		t.Errorf("")
	}

	if len(dest) != 0 {
		t.Errorf("no rows deleted from books_authors for book %d", testBook1.ID)
	}
}

func TestDeleteBookNotExists(t *testing.T) {

	err := db.DeleteBook(-1)
	if err == nil {
		t.Fatalf("expected error: book not exists")
	}
}

func compareBook(a, b *teal.Book) bool {
	author := reflect.DeepEqual(a.Author, b.Author)
	return (a.Title == b.Title &&
		a.ISBN == b.ISBN &&
		a.NumOfPages == b.NumOfPages &&
		a.State == b.State &&
		a.Rating == b.Rating &&
		a.Read == b.Read && author)
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
