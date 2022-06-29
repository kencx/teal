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
	testBook1 = &teal.Book{
		Title:  "FooBar",
		Author: []string{"John Doe"},
		ISBN:   "100",
	}
	testBook2 = &teal.Book{
		Title:  "FooBar",
		Author: []string{"John Doe"},
		ISBN:   "101",
	}
	testBook3 = &teal.Book{
		Title:  "FooBar",
		Author: []string{"John Doe"},
		ISBN:   "102",
	}
	testBooks = []*teal.Book{testBook1, testBook2, testBook3}
)

func TestGetBook(t *testing.T) {
	testServer.Books = &mock.BookStore{
		GetBookFn: func(id int64) (*teal.Book, error) {
			return testBook1, nil
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/books/1",
		data:   nil,
		params: map[string]string{"id": "1"},
		fn:     testServer.GetBook,
	}
	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string]*teal.Book
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["books"]
	assertEqual(t, got.Title, testBook1.Title)
	assertEqual(t, got.Author[0], testBook1.Author[0])
	assertEqual(t, got.ISBN, testBook1.ISBN)
	assertEqual(t, w.Code, http.StatusOK)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestGetBookNil(t *testing.T) {
	testServer.Books = &mock.BookStore{
		GetBookFn: func(id int64) (*teal.Book, error) {
			return nil, teal.ErrDoesNotExist
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/books/1",
		data:   nil,
		params: map[string]string{"id": "1"},
		fn:     testServer.GetBook,
	}

	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string]string
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["error"]
	assertEqual(t, w.Code, http.StatusNotFound)
	assertEqual(t, got, "the item does not exist")
}

func TestGetAllBooks(t *testing.T) {
	testServer.Books = &mock.BookStore{
		GetAllBooksFn: func() ([]*teal.Book, error) {
			return testBooks, nil
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/books/",
		data:   nil,
		params: nil,
		fn:     testServer.GetAllBooks,
	}

	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string][]*teal.Book
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["books"]
	for i, v := range got {
		assertEqual(t, v.Title, testBooks[i].Title)
		assertEqual(t, v.Author[0], testBooks[i].Author[0])
		assertEqual(t, v.ISBN, testBooks[i].ISBN)
	}
	assertEqual(t, w.Code, http.StatusOK)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestGetAllBooksNil(t *testing.T) {
	testServer.Books = &mock.BookStore{
		GetAllBooksFn: func() ([]*teal.Book, error) {
			return nil, teal.ErrNoRows
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/books/",
		data:   nil,
		params: nil,
		fn:     testServer.GetAllBooks,
	}

	w, err := testResponse(t, tc)
	checkErr(t, err)
	assertEqual(t, w.Code, http.StatusNoContent)
}

func TestQueryBooksFromAuthor(t *testing.T) {
	testServer.Books = &mock.BookStore{
		GetByAuthorFn: func(name string) ([]*teal.Book, error) {
			return testBooks, nil
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/books/?author=John+Doe",
		data:   nil,
		params: nil,
		fn:     testServer.GetAllBooks,
	}

	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string][]*teal.Book
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["books"]
	assertEqual(t, w.Code, http.StatusOK)
	assertObjectEqual(t, got, testBooks)
}

func TestNilQueryBooksFromAuthor(t *testing.T) {
	testServer.Books = &mock.BookStore{
		GetByAuthorFn: func(name string) ([]*teal.Book, error) {
			return nil, teal.ErrNoRows
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/books/?author=John+Doe",
		data:   nil,
		params: nil,
		fn:     testServer.GetAllBooks,
	}

	w, err := testResponse(t, tc)
	checkErr(t, err)
	assertEqual(t, w.Code, http.StatusNoContent)
}

func TestAddBook(t *testing.T) {
	want, err := util.ToJSON(testBook1)
	checkErr(t, err)

	testServer.Books = &mock.BookStore{
		CreateBookFn: func(b *teal.Book) (*teal.Book, error) {
			return testBook1, nil
		},
	}

	tc := &testCase{
		method: http.MethodPost,
		url:    "/api/books",
		data:   want,
		params: nil,
		fn:     testServer.AddBook,
	}

	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string]*teal.Book
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["books"]
	assertEqual(t, got.Title, testBook1.Title)
	assertEqual(t, got.Author[0], testBook1.Author[0])
	assertEqual(t, got.ISBN, testBook1.ISBN)
	assertEqual(t, w.Code, http.StatusCreated)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestAddBookFailValidation(t *testing.T) {
	failBook := &teal.Book{
		Title:  "",
		Author: []string{"John Doe"},
		ISBN:   "12345",
	}
	want, err := util.ToJSON(failBook)
	checkErr(t, err)

	testServer.Books = &mock.BookStore{
		CreateBookFn: func(b *teal.Book) (*teal.Book, error) {
			return failBook, nil
		},
	}

	tc := &testCase{
		method: http.MethodPost,
		url:    "/api/books",
		data:   want,
		params: nil,
		fn:     testServer.AddBook,
	}

	w, err := testResponse(t, tc)
	checkErr(t, err)

	// check validation error
	var body map[string]map[string]string
	err = json.NewDecoder(w.Body).Decode(&body)
	checkErr(t, err)

	assertEqual(t, w.Code, http.StatusUnprocessableEntity)
	got := body["error"]

	val, ok := got["title"]
	if !ok {
		t.Errorf("validation error field %q not present", "title")
	}
	assertEqual(t, val, "value is missing")
}

func TestUpdateBook(t *testing.T) {
	want, err := util.ToJSON(testBook2)
	checkErr(t, err)

	testServer.Books = &mock.BookStore{
		UpdateBookFn: func(id int64, b *teal.Book) (*teal.Book, error) {
			return testBook2, nil
		},
	}

	tc := &testCase{
		method: http.MethodPut,
		url:    "/api/books/1",
		data:   want,
		params: map[string]string{"id": "1"},
		fn:     testServer.UpdateBook,
	}

	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string]*teal.Book
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["books"]
	assertEqual(t, got.Title, testBook2.Title)
	assertEqual(t, got.Author[0], testBook2.Author[0])
	assertEqual(t, got.ISBN, testBook2.ISBN)
	assertEqual(t, w.Code, http.StatusOK)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestUpdateBookNil(t *testing.T) {
	want, err := util.ToJSON(testBook2)
	checkErr(t, err)

	testServer.Books = &mock.BookStore{
		UpdateBookFn: func(id int64, b *teal.Book) (*teal.Book, error) {
			return nil, teal.ErrDoesNotExist
		},
	}

	tc := &testCase{
		method: http.MethodPut,
		url:    "/api/books/10",
		data:   want,
		params: map[string]string{"id": "10"},
		fn:     testServer.UpdateBook,
	}

	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string]string
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["error"]
	assertEqual(t, w.Code, http.StatusNotFound)
	assertEqual(t, got, "the item does not exist")
}

func TestUpdateBookFailValidation(t *testing.T) {
	failBook := &teal.Book{
		Title:  "",
		Author: []string{"John Doe"},
		ISBN:   "12345",
	}
	want, err := util.ToJSON(failBook)
	checkErr(t, err)

	testServer.Books = &mock.BookStore{
		UpdateBookFn: func(id int64, b *teal.Book) (*teal.Book, error) {
			return failBook, nil
		},
	}

	tc := &testCase{
		method: http.MethodPut,
		url:    "/api/books/1",
		data:   want,
		params: map[string]string{"id": "1"},
		fn:     testServer.UpdateBook,
	}

	w, err := testResponse(t, tc)
	checkErr(t, err)

	// check validation error
	var body map[string]map[string]string
	err = json.NewDecoder(w.Body).Decode(&body)
	checkErr(t, err)

	assertEqual(t, w.Code, http.StatusUnprocessableEntity)
	got := body["error"]

	val, ok := got["title"]
	if !ok {
		t.Errorf("validation error field %q not present", "title")
	}
	assertEqual(t, val, "value is missing")
}

func TestDeleteBook(t *testing.T) {

	testServer.Books = &mock.BookStore{
		DeleteBookFn: func(id int64) error {
			return nil
		},
	}

	tc := &testCase{
		method: http.MethodDelete,
		url:    "/api/books/1",
		data:   nil,
		params: map[string]string{"id": "1"},
		fn:     testServer.DeleteBook,
	}

	w, err := testResponse(t, tc)
	checkErr(t, err)

	assertEqual(t, w.Code, http.StatusOK)
}

func TestDeleteBookNil(t *testing.T) {

	testServer.Books = &mock.BookStore{
		DeleteBookFn: func(id int64) error {
			return teal.ErrDoesNotExist
		},
	}

	tc := &testCase{
		method: http.MethodDelete,
		url:    "/api/books/10",
		data:   nil,
		params: map[string]string{"id": "10"},
		fn:     testServer.DeleteBook,
	}

	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string]string
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["error"]
	assertEqual(t, w.Code, http.StatusNotFound)
	assertEqual(t, got, "the item does not exist")
}
