package storage

import (
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
		t.Fatalf("expected error")
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

func TestCreateAuthor(t *testing.T) {

	want := &teal.Author{Name: "FooBar"}

	err := db.CreateAuthor(want)
	checkErr(t, err)

	got, err := db.RetrieveAuthorWithName(want.Name)
	checkErr(t, err)

	if got.Name != want.Name {
		t.Errorf("got %v, want %v", got.Name, want.Name)
	}
}
