package storage

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"github.com/kencx/teal"
)

func TestRetrieveBookWithID(t *testing.T) {
	initTestDB(db) // reset DB
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
			assertAuthorsExist(t, tt.want)
			// check books authors entry created
			assertBookAuthorRelationship(t, tt.want)
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
		Title:      "Morning Star",
		ISBN:       "1004",
		Author:     []string{"Pierce Brown"},
		NumOfPages: 100,
		Rating:     10,
		State:      "unread",
	}

	err := db.CreateBook(want)
	checkErr(t, err)

	assertAuthorsExist(t, want)

	tx, err := db.db.Beginx()
	if err != nil {
		t.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	// check books authors table should have two entries for john doe
	books, err := getBooksFromAuthor(tx, want.Author[0])
	checkErr(t, err)

	if len(books) != 2 {
		t.Errorf("got %d books, want %d books for author %q", len(books), 2, want.Author[0])
	}
}

func TestCreateBookNewAndExistingAuthor(t *testing.T) {
	want := &teal.Book{
		Title:      "Tiamat's Wrath",
		ISBN:       "1005",
		Author:     []string{"S.A. Corey", "Daniel Abrahams"},
		NumOfPages: 100,
		Rating:     10,
		State:      "unread",
	}

	err := db.CreateBook(want)
	checkErr(t, err)

	assertAuthorsExist(t, want)

	tx, err := db.db.Beginx()
	if err != nil {
		t.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	num := []int{2, 1}
	for i, v := range want.Author {
		books, err := getBooksFromAuthor(tx, v)
		checkErr(t, err)

		if len(books) != num[i] {
			t.Errorf("got %d books, want %d books for author %q", len(books), num[i], v)
		}
	}
}

func TestUpdateBookNoAuthorChange(t *testing.T) {
	want := testBook1
	want.NumOfPages = 999
	want.Rating = 1
	want.State = "unread"

	err := db.UpdateBook(want.ID, want)
	checkErr(t, err)

	got, err := db.RetrieveBookWithID(want.ID)
	checkErr(t, err)

	if !compareBook(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestUpdateBookAddNewAuthor(t *testing.T) {
	initTestDB(db) // reset DB

	want := testBook1
	want.Author = []string{"S.A. Corey", "Ty Franck"}

	err := db.UpdateBook(want.ID, want)
	checkErr(t, err)

	got, err := db.RetrieveBookWithID(want.ID)
	checkErr(t, err)

	if !compareBook(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	assertAuthorsExist(t, want)
	assertBookAuthorRelationship(t, want)
}

func TestUpdateBookAddExistingAuthor(t *testing.T) {
	initTestDB(db) // reset DB
	want := testBook1
	want.Author = []string{"S.A. Corey", "John Doe"}

	err := db.UpdateBook(want.ID, want)
	checkErr(t, err)

	got, err := db.RetrieveBookWithID(want.ID)
	checkErr(t, err)

	if !compareBook(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	assertAuthorsExist(t, want)
	assertBookAuthorRelationship(t, want)
}

func TestUpdateBookRemoveAuthor(t *testing.T) {
	initTestDB(db) // reset DB
	want := testBook3
	want.Author = []string{"Regina Phallange", "Ken Adams"}

	err := db.UpdateBook(want.ID, want)
	checkErr(t, err)

	got, err := db.RetrieveBookWithID(want.ID)
	checkErr(t, err)

	if !compareBook(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	// check john doe still exists in authors table
	_, err = db.RetrieveAuthorWithName(testAuthor3.Name)
	checkErr(t, err)

	// relationship with john doe dropped
	assertBookAuthorRelationship(t, want)
}

func TestUpdateBookRemoveAuthorCompletely(t *testing.T) {
	initTestDB(db) // reset DB
	want := testBook3
	want.Author = []string{"John Doe", "Regina Phallange"}

	err := db.UpdateBook(want.ID, want)
	checkErr(t, err)

	got, err := db.RetrieveBookWithID(want.ID)
	checkErr(t, err)

	if !compareBook(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	// check ken adams dropped from authors table completely
	_, err = db.RetrieveAuthorWithName(testAuthor5.Name)
	if err == nil {
		t.Errorf("expected error: author does not exist")
	}

	// relationship with ken adams dropped
	assertBookAuthorRelationship(t, want)
}

func TestUpdateBookRenameAuthor(t *testing.T) {
	initTestDB(db) // reset DB
	want := testBook4
	want.Author = []string{"John Adams"}

	err := db.UpdateBook(want.ID, want)
	checkErr(t, err)

	got, err := db.RetrieveBookWithID(want.ID)
	checkErr(t, err)

	if !compareBook(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}

	assertAuthorsExist(t, want)

	// check author still exists
	_, err = db.RetrieveAuthorWithName(testAuthor3.Name)
	checkErr(t, err)

	// relationship with john doe dropped
	// new relationship formed
	assertBookAuthorRelationship(t, want)
}

func TestUpdateBookNotExists(t *testing.T) {
	b := &teal.Book{}
	err := db.UpdateBook(-1, b)
	if err == nil {
		t.Fatalf("expected error: book not exists")
	}
}

func TestUpdateBookISBNConstraint(t *testing.T) {
	want := testBook1
	want.ISBN = testBook2.ISBN
	err := db.UpdateBook(want.ID, want)
	if err == nil {
		t.Errorf("expected error: unique constraint ISBN")
	}
}

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
		t.Errorf("unexpected err: %v", err)
	}

	if len(dest) != 0 {
		t.Errorf("no rows deleted from books_authors for book %d", testBook1.ID)
	}

	// check author entry completely deleted from authors
	_, err = db.RetrieveAuthorWithName(testBook1.Author[0])
	if err == nil {
		t.Errorf("expected error, author %q not deleted", testBook1.Author[0])
	}
}

func TestDeleteBookEnsureAuthorRemainsForExistingBooks(t *testing.T) {
	initTestDB(db) // reset db
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

func assertAuthorsExist(t *testing.T, want *teal.Book) {
	t.Helper()
	for _, author := range want.Author {
		got, err := db.RetrieveAuthorWithName(author)
		checkErr(t, err)

		if got.Name != author {
			t.Errorf("got %v, want %v", got.Name, author)
		}
	}
}

func assertBookAuthorRelationship(t *testing.T, book *teal.Book) {
	t.Helper()
	tx, err := db.db.Beginx()
	if err != nil {
		t.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	// get book's related authors
	authors, err := getAuthorsFromBook(tx, book.ISBN)
	checkErr(t, err)

	if len(authors) != len(book.Author) {
		t.Errorf("book has wrong number of authors in books_authors table")
	}

	// author must exist in book's related authors
	if !reflect.DeepEqual(authors, book.Author) {
		t.Errorf("got %v, want %v", authors, book.Author)

	}
}

func compareBook(a, b *teal.Book) bool {
	author := reflect.DeepEqual(a.Author, b.Author)
	return (a.Title == b.Title &&
		a.ISBN == b.ISBN &&
		a.NumOfPages == b.NumOfPages &&
		a.State == b.State &&
		a.Rating == b.Rating && author)
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
