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

	want := allBooks

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

			// check author entry created
			for _, a := range tt.want.Author {
				got, err := db.RetrieveAuthorWithName(a)
				checkErr(t, err)

				if got.Name != a {
					t.Errorf("got %v, want %v", got.Name, a)
				}
			}

			tx, err := db.db.Beginx()
			if err != nil {
				t.Errorf("db: failed to start transaction: %v", err)
			}
			defer endTx(tx, err)

			// check books authors entry created
			var dest []struct {
				Title string
				Name  string
			}
			stmt := `SELECT b.title, a.name FROM books_authors ba
				JOIN books b ON b.id=ba.book_id
				JOIN authors a ON a.id=ba.author_id
				WHERE b.isbn=$1`
			if err = tx.Select(&dest, stmt, tt.want.ISBN); err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if len(dest) != len(tt.want.Author) {
				t.Errorf("book has wrong number of authors in books_authors table")
			}

			for i, d := range dest {
				if d.Title != tt.want.Title {
					t.Errorf("got %v, want %v", d.Title, tt.want.Title)
				}
				if d.Name != tt.want.Author[i] {
					t.Errorf("got %v, want %v", d.Name, tt.want.Author[i])
				}
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

	// check books authors table should have two entries
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

	// check books authors table should have 3 entries for John Doe, 1 for Newest Author
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
		t.Errorf("no rows deleted from books_authors for book %d", testAuthor1.ID)
	}

	if len(dest) != 0 {
		t.Errorf("no rows deleted from books_authors for book %d", testBook1.ID)
	}

	// check author entry deleted from authors
	var adest []string
	stmt = `SELECT a.id FROM authors a
		JOIN books_authors ba ON ba.author_id=a.id WHERE ba.book_id=$1`
	if err := tx.Select(&adest, stmt, testBook1.ID); err != nil {
		t.Errorf("no rows deleted from authors for book %d", testBook1.ID)
	}

	if len(dest) != 0 {
		t.Errorf("no rows deleted from authors for book %d", testBook1.ID)
	}
}

func TestDeleteBookEnsureAuthorRemainsForExistingBooks(t *testing.T) {

	err := db.DeleteBook(testBook3.ID)
	checkErr(t, err)

	// check author still exists in authors table
	got, err := db.RetrieveAuthorWithName(testBook3.Author[0])
	checkErr(t, err)

	if got.Name != testBook3.Author[0] {
		t.Errorf("got %v, want %v", got.Name, testBook3.Author[0])
	}

	tx, err := db.db.Beginx()
	if err != nil {
		t.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	// check author still linked to their other books in books_authors
	var dest []struct {
		Book_id   int
		Author_id int
	}
	stmt := `SELECT ba.book_id, ba.author_id
		FROM books_authors ba
		JOIN authors a ON a.id=ba.author_id
		JOIN books b ON b.id=ba.book_id
		WHERE a.name=$1`

	if err := tx.Select(&dest, stmt, testBook3.Author[0]); err != nil {
		t.Errorf("error")
	}

	if len(dest) < 1 {
		t.Errorf("error")
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
