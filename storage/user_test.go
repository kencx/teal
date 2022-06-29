package storage

import (
	"reflect"
	"testing"

	"github.com/kencx/teal"
)

func TestGetUser(t *testing.T) {
	got, err := ts.Users.Get(testUser1.ID)
	checkErr(t, err)

	want := testUser1
	if !assertUsersEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestGetUserWithUsername(t *testing.T) {
	got, err := ts.Users.GetByUsername(testUser1.Username)
	checkErr(t, err)

	want := testUser1
	if !assertUsersEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestGetUserNotExists(t *testing.T) {
	result, err := ts.Users.Get(-1)
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

func TestCreateUser(t *testing.T) {
	want := &teal.User{
		Name:           "Bar Baz",
		Username:       "barbaz",
		HashedPassword: []byte("abc123456789"),
	}

	got, err := ts.Users.Create(want)
	checkErr(t, err)

	if !assertUsersEqual(got, want) {
		t.Errorf("got %v, want %v", prettyPrint(got), prettyPrint(want))
	}
}

func TestCreateUserDuplicateUsername(t *testing.T) {
	want := testUser1
	result, err := ts.Users.Create(want)

	if err == nil {
		t.Errorf("expected error: ErrDuplicateUsername")
	}
	if err != teal.ErrDuplicateUsername {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatalf("got %v, want nil", result)
	}
}

func TestUpdateUser(t *testing.T) {
	defer resetDB(testdb)

	want := testUser1
	want.Username = "Foo Bar"

	got, err := ts.Users.Update(want.ID, want)
	checkErr(t, err)

	if !assertUsersEqual(got, want) {
		t.Errorf("got %v, want %v", got.Name, want.Name)
	}
}

func TestUpdateUserExistingUsername(t *testing.T) {
	want := testUser1
	want.Username = testUser2.Username

	result, err := ts.Users.Update(want.ID, want)
	if err == nil {
		t.Errorf("expected error: ErrDuplicateUsername")
	}
	if err != teal.ErrDuplicateUsername {
		t.Errorf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("got %v, want nil", result)
	}
}

func TestDeleteUser(t *testing.T) {
	err := ts.Users.Delete(testUser1.ID)
	checkErr(t, err)

	_, err = ts.Users.Get(testUser1.ID)
	if err == nil {
		t.Errorf("expected error, user %d not deleted", testUser1.ID)
	}
}

func TestDeleteUserNotExists(t *testing.T) {
	err := ts.Users.Delete(testUser1.ID)
	if err == nil {
		t.Errorf("expected error: user does not exists")
	}
}

func assertUsersEqual(a, b *teal.User) bool {
	passwordHashEqual := reflect.DeepEqual(a.HashedPassword, b.HashedPassword)

	return (a.ID == b.ID &&
		a.Name == b.Name &&
		a.Username == b.Username &&
		a.Role == b.Role &&
		passwordHashEqual)
}
