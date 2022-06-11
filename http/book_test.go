package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/kencx/teal"
	"github.com/kencx/teal/http/response"
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
	s := Server{
		// TODO mock loggers
		InfoLog: log.New(os.Stdout, "", log.LstdFlags),
		ErrLog:  log.New(os.Stdout, "", log.LstdFlags),
		Books: &mockBookService{
			getBookFn: func(id int) (*teal.Book, error) {
				return testBook1, nil
			},
		},
	}

	w, err := getResponse("/1", s.GetBook)
	checkErr(t, err)

	var got teal.Book
	err = json.NewDecoder(w.Body).Decode(&got)
	checkErr(t, err)

	assertEqual(t, got.Title, testBook1.Title)
	assertEqual(t, got.Author[0], testBook1.Author[0])
	assertEqual(t, got.ISBN, testBook1.ISBN)
	assertEqual(t, w.Code, http.StatusOK)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestGetAllBooks(t *testing.T) {
	s := Server{
		InfoLog: log.New(os.Stdout, "", log.LstdFlags),
		ErrLog:  log.New(os.Stdout, "", log.LstdFlags),
		Books: &mockBookService{
			getAllBooksFn: func() ([]*teal.Book, error) {
				return testBooks, nil
			},
		},
	}

	w, err := getResponse("/", s.GetAllBooks)
	checkErr(t, err)

	var got []teal.Book
	err = json.NewDecoder(w.Body).Decode(&got)
	checkErr(t, err)

	for i, v := range got {
		assertEqual(t, v.Title, testBooks[i].Title)
		assertEqual(t, v.Author[0], testBooks[i].Author[0])
		assertEqual(t, v.ISBN, testBooks[i].ISBN)
	}
	assertEqual(t, w.Code, http.StatusOK)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestAddBook(t *testing.T) {
	want, err := util.ToJSON(testBook1)
	checkErr(t, err)

	s := Server{
		InfoLog: log.New(os.Stdout, "", log.LstdFlags),
		ErrLog:  log.New(os.Stdout, "", log.LstdFlags),
		Books: &mockBookService{
			createBookFn: func(ctx context.Context, b *teal.Book) (*teal.Book, error) {
				return testBook1, nil
			},
		}}

	w, err := postResponse("/", bytes.NewReader(want), s.AddBook)
	checkErr(t, err)

	var got teal.Book
	err = json.NewDecoder(w.Body).Decode(&got)
	checkErr(t, err)

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
		InfoLog: log.New(os.Stdout, "", log.LstdFlags),
		ErrLog:  log.New(os.Stdout, "", log.LstdFlags),
		Books: &mockBookService{
			createBookFn: func(ctx context.Context, b *teal.Book) (*teal.Book, error) {
				return failBook, nil
			},
		}}

	w, err := postResponse("/", bytes.NewBuffer(want), s.AddBook)
	checkErr(t, err)

	// get response
	var body response.ValidationErrResponse
	err = json.NewDecoder(w.Body).Decode(&body)
	checkErr(t, err)

	assertEqual(t, w.Code, http.StatusBadRequest)
	for _, v := range body.Err {
		log.Println(v.Message)
		strings.Contains(v.Message, "title")
	}
}

func getResponse(url string, f func(http.ResponseWriter, *http.Request)) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	w := httptest.NewRecorder()

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
