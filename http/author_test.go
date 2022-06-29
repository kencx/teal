package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kencx/teal"
	"github.com/kencx/teal/mock"
	"github.com/kencx/teal/util"
)

var (
	testAuthor1 = &teal.Author{
		Name: "Author 1",
	}
	testAuthor2 = &teal.Author{
		Name: "Author 2",
	}
	testAuthors = []*teal.Author{testAuthor1, testAuthor2}
)

func TestGetAuthor(t *testing.T) {
	testServer.Authors = &mock.AuthorStore{
		GetAuthorFn: func(id int64) (*teal.Author, error) {
			return testAuthor1, nil
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/authors/1",
		data:   nil,
		params: map[string]string{"id": "1"},
		fn:     testServer.GetAuthor,
	}
	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string]*teal.Author
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["authors"]
	assertEqual(t, got.Name, testAuthor1.Name)
	assertEqual(t, w.Code, http.StatusOK)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestGetAllAuthors(t *testing.T) {

	testServer.Authors = &mock.AuthorStore{
		GetAllAuthorsFn: func() ([]*teal.Author, error) {
			return testAuthors, nil
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/authors/",
		data:   nil,
		params: nil,
		fn:     testServer.GetAllAuthors,
	}
	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string][]*teal.Author
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["authors"]
	for i, v := range got {
		assertEqual(t, v.Name, testAuthors[i].Name)
	}
	assertEqual(t, w.Code, http.StatusOK)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestAddAuthor(t *testing.T) {
	want, err := util.ToJSON(testAuthor1)
	checkErr(t, err)

	testServer.Authors = &mock.AuthorStore{
		CreateAuthorFn: func(a *teal.Author) (*teal.Author, error) {
			return testAuthor1, nil
		},
	}

	tc := &testCase{
		method: http.MethodPost,
		url:    "/api/authors/",
		data:   want,
		params: nil,
		fn:     testServer.AddAuthor,
	}
	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string]*teal.Author
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["authors"]
	assertEqual(t, got.Name, testAuthor1.Name)
	assertEqual(t, w.Code, http.StatusCreated)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestUpdateAuthor(t *testing.T) {
	want, err := util.ToJSON(testAuthor2)
	checkErr(t, err)

	testServer.Authors = &mock.AuthorStore{
		UpdateAuthorFn: func(id int64, a *teal.Author) (*teal.Author, error) {
			return testAuthor2, nil
		},
	}

	tc := &testCase{
		method: http.MethodPut,
		url:    "/api/authors/1",
		data:   want,
		params: map[string]string{"id": "1"},
		fn:     testServer.UpdateAuthor,
	}
	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string]*teal.Author
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["authors"]
	assertEqual(t, got.Name, testAuthor2.Name)
	assertEqual(t, w.Code, http.StatusOK)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestDeleteAuthor(t *testing.T) {

	testServer.Authors = &mock.AuthorStore{
		DeleteAuthorFn: func(id int64) error {
			return nil
		},
	}

	tc := &testCase{
		method: http.MethodDelete,
		url:    "/api/authors/1",
		data:   nil,
		params: map[string]string{"id": "1"},
		fn:     testServer.DeleteAuthor,
	}
	w, err := testResponse(t, tc)
	checkErr(t, err)
	assertEqual(t, w.Code, http.StatusOK)
}
