package storage

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/kencx/teal"
	"github.com/kencx/teal/storage/sqlite/model"
	. "github.com/kencx/teal/storage/sqlite/table"
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestRetrieveBookWithTitle(t *testing.T) {
	got, err := db.RetrieveBookWithTitle(testBook2.Title)
	checkErr(t, err)

	want := testBook2
	if !reflect.DeepEqual(got, want) {
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

	want := []*teal.Book{testBook1, testBook2, testBook3}

	if len(got) != len(want) {
		t.Fatalf("got %d books, want %d books", len(got), len(want))
	}

	for i := 0; i < len(got); i++ {
		if !compareBook(got[i], want[i]) {
			t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
		}
		if !compareAuthors(got[i].Author, want[i].Author) {
			t.Errorf("got %v, want %v", prettyPrint(got[i].Author), prettyPrint(want[i].Author))
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
			Author: teal.Authors{{Name: "George Orwell"}},
		},
	}, {
		name: "book with all data",
		want: &teal.Book{
			Title:      "World War Z",
			ISBN:       "1002",
			Author:     teal.Authors{{Name: "Max Brooks"}},
			NumOfPages: 100,
			Rating:     10,
			State:      "read",
		},
	}, {
		name: "book with two authors",
		want: &teal.Book{
			Title:      "Pro Git",
			ISBN:       "1003",
			Author:     teal.Authors{{Name: "Scott Chacon"}, {Name: "Ben Straub"}},
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
			if !compareAuthors(got.Author, tt.want.Author) {
				t.Errorf("got %v, want %v", prettyPrint(got.Author), prettyPrint(tt.want.Author))
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
		Author:     teal.Authors{{Name: "John Doe"}},
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

	var dest []struct {
		model.Authors
	}

	// check for number of entries in authors
	if err := SELECT(Authors.Name).
		FROM(Authors).
		WHERE(Authors.Name.EQ(String(want.Author[0].Name))).
		Query(tx, &dest); err != nil {
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
		Author:     teal.Authors{{Name: "John Doe"}, {Name: "Newest Author"}},
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

	var dest []struct {
		model.Authors
	}

	// check for number of entries in authors
	if err := SELECT(Authors.Name).
		FROM(Authors).
		WHERE(Authors.Name.EQ(String(want.Author[0].Name)).
			OR(Authors.Name.EQ(String(want.Author[1].Name)))).
		Query(tx, &dest); err != nil {
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

// func TestDelete(t *testing.T) {
// 	db := setup()
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
// 	db := setup()
//
// 	err := deleteBook(db, 999)
// 	if err == nil {
// 		t.Fatalf("expected error: book not exists")
// 	}
// }

func compareBook(a, b *teal.Book) bool {
	return (a.Title == b.Title &&
		a.ISBN == b.ISBN &&
		a.NumOfPages == b.NumOfPages &&
		a.State == b.State &&
		a.Rating == b.Rating &&
		a.Read == b.Read)
}

func compareAuthors(a, b teal.Authors) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i].Name != b[i].Name {
			return false
		}
	}
	return true
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
