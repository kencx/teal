package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
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

	testInfoLog = log.New(io.Discard, "", log.LstdFlags)
	testErrLog  = log.New(io.Discard, "", log.LstdFlags)
)

func TestGetBook(t *testing.T) {
	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Books: &mock.BookStore{
			GetBookFn: func(id int64) (*teal.Book, error) {
				return testBook1, nil
			},
		},
	}

	w, err := getResponse("/api/books/1", s.GetBook)
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
	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Books: &mock.BookStore{
			GetBookFn: func(id int64) (*teal.Book, error) {
				return nil, teal.ErrDoesNotExist
			},
		},
	}

	w, err := getResponse("/api/books/1", s.GetBook)
	checkErr(t, err)

	var env map[string]string
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["error"]
	assertEqual(t, w.Code, http.StatusNotFound)
	assertEqual(t, got, "the item does not exist")
}

func TestGetAllBooks(t *testing.T) {
	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Books: &mock.BookStore{
			GetAllBooksFn: func() ([]*teal.Book, error) {
				return testBooks, nil
			},
		},
	}

	w, err := getResponse("/api/books", s.GetAllBooks)
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
	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Books: &mock.BookStore{
			GetAllBooksFn: func() ([]*teal.Book, error) {
				return nil, teal.ErrNoRows
			},
		},
	}

	w, err := getResponse("/api/books", s.GetAllBooks)
	checkErr(t, err)

	assertEqual(t, w.Code, http.StatusNoContent)
}

func TestQueryBooksFromAuthor(t *testing.T) {
	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Books: &mock.BookStore{
			GetByAuthorFn: func(name string) ([]*teal.Book, error) {
				return testBooks, nil
			},
		}}

	w, err := getResponse("/api/books/?author=John+Doe", s.GetAllBooks)
	checkErr(t, err)

	var env map[string][]*teal.Book
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["books"]
	assertEqual(t, w.Code, http.StatusOK)
	assertObjectEqual(t, got, testBooks)
}

func TestNilQueryBooksFromAuthor(t *testing.T) {
	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Books: &mock.BookStore{
			GetByAuthorFn: func(name string) ([]*teal.Book, error) {
				return nil, teal.ErrNoRows
			},
		}}

	w, err := getResponse("/api/books/?author=John+Doe", s.GetAllBooks)
	checkErr(t, err)

	assertEqual(t, w.Code, http.StatusNoContent)
}

func TestAddBook(t *testing.T) {
	want, err := util.ToJSON(testBook1)
	checkErr(t, err)

	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Books: &mock.BookStore{
			CreateBookFn: func(b *teal.Book) (*teal.Book, error) {
				return testBook1, nil
			},
		}}

	w, err := postResponse("/api/books", bytes.NewReader(want), s.AddBook)
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

	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Books: &mock.BookStore{
			CreateBookFn: func(b *teal.Book) (*teal.Book, error) {
				return failBook, nil
			},
		}}

	w, err := postResponse("/api/books", bytes.NewBuffer(want), s.AddBook)
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

	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Books: &mock.BookStore{
			UpdateBookFn: func(id int64, b *teal.Book) (*teal.Book, error) {
				return testBook2, nil
			},
		}}
	w, err := putResponse("/api/books/1", bytes.NewBuffer(want), s.UpdateBook)
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

	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Books: &mock.BookStore{
			UpdateBookFn: func(id int64, b *teal.Book) (*teal.Book, error) {
				return nil, teal.ErrDoesNotExist
			},
		},
	}

	w, err := putResponse("/api/books/10", bytes.NewBuffer(want), s.UpdateBook)
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

	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Books: &mock.BookStore{
			UpdateBookFn: func(id int64, b *teal.Book) (*teal.Book, error) {
				return failBook, nil
			},
		}}

	w, err := putResponse("/api/books/1", bytes.NewBuffer(want), s.UpdateBook)
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

	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Books: &mock.BookStore{
			DeleteBookFn: func(id int64) error {
				return nil
			},
		}}

	w, err := deleteResponse("/api/books/1", s.DeleteBook)
	checkErr(t, err)

	assertEqual(t, w.Code, http.StatusOK)
}

func TestDeleteBookNil(t *testing.T) {

	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
		Books: &mock.BookStore{
			DeleteBookFn: func(id int64) error {
				return teal.ErrDoesNotExist
			},
		},
	}

	w, err := deleteResponse("/api/books/10", s.DeleteBook)
	checkErr(t, err)

	var env map[string]string
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["error"]
	assertEqual(t, w.Code, http.StatusNotFound)
	assertEqual(t, got, "the item does not exist")
}

func getResponse(url string, f func(http.ResponseWriter, *http.Request)) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	http.HandlerFunc(f).ServeHTTP(w, req)
	return w, nil
}

func deleteResponse(url string, f func(http.ResponseWriter, *http.Request)) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	http.HandlerFunc(f).ServeHTTP(w, req)
	return w, nil
}

func postResponse(url string, data io.Reader, f func(http.ResponseWriter, *http.Request)) (*httptest.ResponseRecorder, error) {

	req, err := http.NewRequest(http.MethodPost, url, data)
	if err != nil {
		return nil, err
	}
	w := httptest.NewRecorder()

	http.HandlerFunc(f).ServeHTTP(w, req)
	return w, nil
}

func putResponse(url string, data io.Reader, f func(http.ResponseWriter, *http.Request)) (*httptest.ResponseRecorder, error) {

	req, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		return nil, err
	}
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	http.HandlerFunc(f).ServeHTTP(w, req)
	return w, nil
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertObjectEqual(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func checkErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
