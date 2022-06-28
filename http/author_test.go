package http

import (
	"bytes"
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
	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Authors: &mock.AuthorStore{
			GetAuthorFn: func(id int) (*teal.Author, error) {
				return testAuthor1, nil
			},
		},
	}

	w, err := getResponse("/api/authors/1", s.GetAuthor)
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
	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Authors: &mock.AuthorStore{
			GetAllAuthorsFn: func() ([]*teal.Author, error) {
				return testAuthors, nil
			},
		},
	}

	w, err := getResponse("/api/authors/", s.GetAllAuthors)
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

	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Authors: &mock.AuthorStore{
			CreateAuthorFn: func(a *teal.Author) (*teal.Author, error) {
				return testAuthor1, nil
			},
		},
	}

	w, err := postResponse("/api/authors/", bytes.NewBuffer(want), s.AddAuthor)
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

	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Authors: &mock.AuthorStore{
			UpdateAuthorFn: func(id int, a *teal.Author) (*teal.Author, error) {
				return testAuthor2, nil
			},
		},
	}

	w, err := putResponse("/api/authors/1", bytes.NewBuffer(want), s.UpdateAuthor)
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
	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Authors: &mock.AuthorStore{
			DeleteAuthorFn: func(id int) error {
				return nil
			},
		},
	}

	w, err := deleteResponse("/api/authors/1", s.DeleteAuthor)
	checkErr(t, err)

	assertEqual(t, w.Code, http.StatusOK)
}
