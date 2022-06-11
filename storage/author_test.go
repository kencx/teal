package storage

import (
	"context"
	"reflect"
	"testing"

	"github.com/kencx/teal"
)

func TestRetrieveAuthorWithID(t *testing.T) {
	got, err := db.RetrieveAuthorWithID(testAuthor1.ID)
	checkErr(t, err)

	want := testAuthor1
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestRetrieveAuthorWithName(t *testing.T) {
	got, err := db.RetrieveAuthorWithName(testAuthor2.Name)
	checkErr(t, err)

	want := testAuthor2
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestRetrieveAuthorNotExists(t *testing.T) {
	result, err := db.RetrieveAuthorWithID(-1)
	if err == nil {
		t.Fatalf("expected error: ErrDoesNotExist")
	}

	if err != teal.ErrDoesNotExist {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != nil {
		t.Fatalf("got %v, want nil", result)
	}
}

func TestRetrieveAllAuthors(t *testing.T) {
	got, err := db.RetrieveAllAuthors()
	checkErr(t, err)

	want := []*teal.Author{testAuthor1, testAuthor2, testAuthor3, testAuthor4, testAuthor5}

	if len(got) != len(want) {
		t.Fatalf("got %d books, want %d books", len(got), len(want))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

// TODO
// func TestRetrieveAllAuthorEmpty(t *testing.T) {
// 	// delete all entries
// 	got, err := db.RetrieveAllAuthors()
//
// 	if err == nil {
// 		t.Fatalf("expected error: ErrNoRows")
// 	}
//
// 	if err != teal.ErrNoRows {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
//
// 	if got != nil {
// 		t.Fatalf("got %v, want nil", got)
// 	}
// }

func TestCreateAuthor(t *testing.T) {

	want := &teal.Author{Name: "FooBar"}

	err := db.CreateAuthor(context.Background(), want)
	checkErr(t, err)

	got, err := db.RetrieveAuthorWithName(want.Name)
	checkErr(t, err)

	if got.Name != want.Name {
		t.Errorf("got %v, want %v", got.Name, want.Name)
	}
}

func TestInsertAuthorDuplicates(t *testing.T) {

	tx, err := db.db.Beginx()
	if err != nil {
		t.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	want := testAuthor3

	_, err = insertOrGetAuthor(tx, want)
	checkErr(t, err)

	// check for number of entries in authors
	var dest []string

	stmt := `SELECT name FROM authors WHERE name=$1`
	if err := tx.Select(&dest, stmt, want.Name); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(dest) != 1 {
		t.Error("more than one author inserted")
	}
}

func TestUpdateAuthor(t *testing.T) {

	want := testAuthor1
	want.Name = "John Watson"

	err := db.UpdateAuthor(context.Background(), want.ID, want)
	checkErr(t, err)

	got, err := db.RetrieveAuthorWithID(want.ID)
	checkErr(t, err)

	if got.Name != want.Name {
		t.Errorf("got %v, want %v", got.Name, want.Name)
	}
}

func TestUpdateAuthorExisting(t *testing.T) {

	want := testAuthor1
	want.Name = "John Doe"

	err := db.UpdateAuthor(context.Background(), want.ID, want)
	if err == nil {
		t.Errorf("expected error: unique constraint Name")
	}
}

func TestDeleteAuthor(t *testing.T) {

	err := db.DeleteAuthor(context.Background(), testAuthor1.ID)
	checkErr(t, err)

	_, err = db.RetrieveAuthorWithID(testAuthor1.ID)
	if err == nil {
		t.Errorf("expected error, author %d not deleted", testAuthor1.ID)
	}

	tx, err := db.db.Beginx()
	if err != nil {
		t.Errorf("db: failed to start transaction: %v", err)
	}
	defer endTx(tx, err)

	// check entries deleted from books_authors
	var dest []int
	stmt := `SELECT book_id FROM books_authors WHERE books_authors.author_id=$1`
	if err := tx.Select(&dest, stmt, testAuthor1.ID); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(dest) != 0 {
		t.Errorf("no rows deleted from books_authors for author %d", testAuthor1.ID)
	}
}

func TestDeleteAuthorNotExists(t *testing.T) {
	err := db.DeleteAuthor(context.Background(), testAuthor1.ID)
	if err == nil {
		t.Errorf("expected error: author not exists")
	}
}
