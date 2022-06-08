package storage

import (
	"reflect"
	"testing"

	// . "github.com/go-jet/jet/v2/sqlite"
	"github.com/kencx/teal"
	// . "github.com/kencx/teal/storage/sqlite/table"
)

func TestRetrieveAuthorWithID(t *testing.T) {
	db := setup(t)

	got, err := db.RetrieveAuthorWithID(testAuthor1.ID)
	checkErr(t, err)

	want := testAuthor1
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestRetrieveAuthorWithName(t *testing.T) {
	db := setup(t)

	got, err := db.RetrieveAuthorWithName(testAuthor2.Name)
	checkErr(t, err)

	want := testAuthor2
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestRetrieveAuthorNotExists(t *testing.T) {
	db := setup(t)

	result, err := db.RetrieveAuthorWithID(-1)
	if err == nil {
		t.Fatalf("expected error")
	}

	if result != nil {
		t.Fatalf("got %v, want nil", result)
	}
}

func TestRetrieveAllAuthors(t *testing.T) {
	db := setup(t)

	got, err := db.RetrieveAllAuthors()
	checkErr(t, err)

	want := []*teal.Author{testAuthor1, testAuthor2, testAuthor3}

	if len(got) != len(want) {
		t.Fatalf("got %d books, want %d books", len(got), len(want))
	}

	// for i := 0; i < len(got); i++ {
	// 	if !compareAuthor(got[i], want[i]) {
	// 		t.Errorf("got %v, want %v", prettyPrint(got[i]), prettyPrint(want[i]))
	// 	}
	// 	if !compareAuthors(got[i].Author, want[i].Author) {
	// 		t.Errorf("got %v, want %v", prettyPrint(got[i].Author), prettyPrint(want[i].Author))
	// 	}
	// }
}
