package service

import (
	"reflect"
	"testing"

	"github.com/kencx/teal/pkg"
)

type mockBookRepository struct {
	bookList []*pkg.Book
}

func (m *mockBookRepository) GetAllBooks() ([]*pkg.Book, error) {
	return m.bookList, nil
}

func (m *mockBookRepository) GetBook(id int) (*pkg.Book, error) {
	return m.bookList[id-1], nil
}

func (m *mockBookRepository) GetBookByTitle(title string) (*pkg.Book, error) {
	for _, b := range m.bookList {
		if b.Title == title {
			return b, nil
		}
	}
	return nil, nil
}

func (m *mockBookRepository) CreateBook(b *pkg.Book) (int, error) {
	m.bookList = append(m.bookList, b)
	return len(m.bookList), nil
}

func (m *mockBookRepository) UpdateBook(id int, b *pkg.Book) error {
	m.bookList[id-1] = b
	return nil
}

func (m *mockBookRepository) DeleteBook(id int) error {
	m.bookList[id-1] = nil
	return nil
}

func TestGetBook(t *testing.T) {

	t.Run("Get book", func(t *testing.T) {

		expected := &pkg.Book{
			Title:  "FooBar",
			Author: "John Doe",
			ISBN:   "12345",
		}
		s := NewService(&mockBookRepository{
			bookList: []*pkg.Book{expected},
		})

		result, err := s.GetBook(1)
		checkErr(t, err)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, want %v", result, expected)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		s := NewService(&mockBookRepository{})

		result, err := s.GetBook(-1)
		if err == nil {
			t.Errorf("got %v, want %v", nil, err)
		}
		if result != nil {
			t.Errorf("got %v, want %v", result, nil)
		}
	})
}

func TestGetAllBooks(t *testing.T) {
	expected := []*pkg.Book{
		{
			Title:  "FooBar",
			Author: "John Doe",
			ISBN:   "12345",
		},
		{
			Title:  "BarBaz",
			Author: "John Doe",
			ISBN:   "45678",
		},
	}
	s := NewService(&mockBookRepository{
		bookList: expected,
	})

	result, err := s.GetAllBooks()
	checkErr(t, err)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestCreateBook(t *testing.T) {
	s := NewService(&mockBookRepository{})

	expected := &pkg.Book{
		Title:  "BarBaz",
		Author: "John Doe",
		ISBN:   "45678",
	}

	id, err := s.CreateBook(expected)
	checkErr(t, err)

	if id != 1 {
		t.Errorf("got %d, want %d", id, 1)
	}

	result, err := s.GetBook(id)
	checkErr(t, err)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestUpdateBook(t *testing.T) {
	s := NewService(&mockBookRepository{
		bookList: []*pkg.Book{
			{
				Title:  "FooBar",
				Author: "Ben Adams",
				ISBN:   "45678",
			},
		},
	})

	expected := &pkg.Book{
		Title:  "BarBaz",
		Author: "John Doe",
		ISBN:   "45678",
	}

	err := s.UpdateBook(1, expected)
	checkErr(t, err)

	result, err := s.GetBook(1)
	checkErr(t, err)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestDeleteBook(t *testing.T) {
	t.Run("delete book", func(t *testing.T) {
		s := NewService(&mockBookRepository{
			bookList: []*pkg.Book{
				{
					Title:  "FooBar",
					Author: "Ben Adams",
					ISBN:   "45678",
				},
			},
		})

		err := s.DeleteBook(1)
		checkErr(t, err)

		result, err := s.GetBook(1)
		checkErr(t, err)

		if result != nil {
			t.Errorf("got %v, want %v", result, nil)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		s := NewService(&mockBookRepository{})

		err := s.DeleteBook(-1)
		if err == nil {
			t.Errorf("got %v, want %v", nil, err)
		}
	})
}

func checkErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
