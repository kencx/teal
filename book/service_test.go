package book

import (
	"reflect"
	"testing"

	teal "github.com/kencx/teal"
)

type mockBookRepository struct {
	bookList []*teal.Book
}

func (m *mockBookRepository) GetAllBooks() ([]*teal.Book, error) {
	return m.bookList, nil
}

func (m *mockBookRepository) GetBook(id int) (*teal.Book, error) {
	return m.bookList[id-1], nil
}

func (m *mockBookRepository) GetBookByTitle(title string) (*teal.Book, error) {
	for _, b := range m.bookList {
		if b.Title == title {
			return b, nil
		}
	}
	return nil, nil
}

func (m *mockBookRepository) CreateBook(b *teal.Book) (int, error) {
	m.bookList = append(m.bookList, b)
	return len(m.bookList), nil
}

func (m *mockBookRepository) UpdateBook(id int, b *teal.Book) error {
	m.bookList[id-1] = b
	return nil
}

func (m *mockBookRepository) DeleteBook(id int) error {
	m.bookList[id-1] = nil
	return nil
}

func TestGet(t *testing.T) {

	t.Run("Get book", func(t *testing.T) {

		expected := &teal.Book{
			Title:  "FooBar",
			Author: "John Doe",
			ISBN:   "12345",
		}
		s := NewService(&mockBookRepository{
			bookList: []*teal.Book{expected},
		})

		result, err := s.Get(1)
		checkErr(t, err)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, want %v", result, expected)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		s := NewService(&mockBookRepository{})

		result, err := s.Get(-1)
		if err == nil {
			t.Errorf("got %v, want %v", nil, err)
		}
		if result != nil {
			t.Errorf("got %v, want %v", result, nil)
		}
	})
}

func TestGetAll(t *testing.T) {
	expected := []*teal.Book{
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

	result, err := s.GetAll()
	checkErr(t, err)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestCreate(t *testing.T) {
	s := NewService(&mockBookRepository{})

	expected := &teal.Book{
		Title:  "BarBaz",
		Author: "John Doe",
		ISBN:   "45678",
	}

	id, err := s.Create(expected)
	checkErr(t, err)

	if id != 1 {
		t.Errorf("got %d, want %d", id, 1)
	}

	result, err := s.Get(id)
	checkErr(t, err)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

// func TestUpdate(t *testing.T) {
// 	s := NewService(&mockBookRepository{
// 		bookList: []*teal.Book{
// 			{
// 				Title:  "FooBar",
// 				Author: "Ben Adams",
// 				ISBN:   "45678",
// 			},
// 		},
// 	})
//
// 	expected := &teal.Book{
// 		Title:  "BarBaz",
// 		Author: "John Doe",
// 		ISBN:   "45678",
// 	}
//
// 	err := s.Update(1, expected)
// 	checkErr(t, err)
//
// 	result, err := s.Get(1)
// 	checkErr(t, err)
//
// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("got %v, want %v", result, expected)
// 	}
// }

func TestDelete(t *testing.T) {
	t.Run("delete book", func(t *testing.T) {
		s := NewService(&mockBookRepository{
			bookList: []*teal.Book{
				{
					Title:  "FooBar",
					Author: "Ben Adams",
					ISBN:   "45678",
				},
			},
		})

		err := s.Delete(1)
		checkErr(t, err)

		result, err := s.Get(1)
		checkErr(t, err)

		if result != nil {
			t.Errorf("got %v, want %v", result, nil)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		s := NewService(&mockBookRepository{})

		err := s.Delete(-1)
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