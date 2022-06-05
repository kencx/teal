package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kencx/teal/pkg"
)

type mockBookService struct {
	getAllBooksFn    func() ([]*pkg.Book, error)
	getBookFn        func(id int) (*pkg.Book, error)
	getBookByTitleFn func(title string) (*pkg.Book, error)
	createBookFn     func(b *pkg.Book) (int, error)
	updateBookFn     func(id int, b *pkg.Book) error
	deleteBookFn     func(id int) error
}

func (m *mockBookService) GetAllBooks() ([]*pkg.Book, error) {
	return m.getAllBooksFn()
}

func (m *mockBookService) GetBook(id int) (*pkg.Book, error) {
	return m.getBookFn(id)
}

func (m *mockBookService) GetBookByTitle(title string) (*pkg.Book, error) {
	return m.getBookByTitleFn(title)
}

func (m *mockBookService) CreateBook(b *pkg.Book) (int, error) {
	return m.createBookFn(b)
}

func (m *mockBookService) UpdateBook(id int, b *pkg.Book) error {
	return m.updateBookFn(id, b)
}

func (m *mockBookService) DeleteBook(id int) error {
	return m.deleteBookFn(id)
}

func TestGetAllBooks(t *testing.T) {
	expected := []*pkg.Book{{Title: "FooBar", Author: "John Doe", ISBN: "52634"}}
	s := Server{books: &mockBookService{
		getAllBooksFn: func() ([]*pkg.Book, error) {
			return expected, nil
		},
	}}

	body, code, header := getResponse(t, http.MethodGet, "/", s.GetAllBooks)
	result := []pkg.Book{}
	err := FromJSON(body, &result)
	checkErr(t, err)

	assertEqual(t, result[0].Title, expected[0].Title)
	assertEqual(t, code, http.StatusOK)
	assertEqual(t, header.Get("Content-Type"), "application/json")
}

func TestGetBook(t *testing.T) {
	expected := &pkg.Book{Title: "FooBar", Author: "John Doe", ISBN: "52634"}
	s := Server{books: &mockBookService{
		getBookFn: func(id int) (*pkg.Book, error) {
			return expected, nil
		},
	}}

	body, code, header := getResponse(t, http.MethodGet, "/1", s.GetBook)
	want, err := ToJSON(expected)
	checkErr(t, err)

	assertEqual(t, string(body), string(want))
	assertEqual(t, code, http.StatusOK)
	assertEqual(t, header.Get("Content-Type"), "application/json")
}

func getResponse(t *testing.T, method, url string, f func(rw http.ResponseWriter, r *http.Request)) ([]byte, int, http.Header) {
	t.Helper()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, url, nil)
	checkErr(t, err)

	http.HandlerFunc(f).ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	checkErr(t, err)

	return body, res.StatusCode, res.Header
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func checkErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
