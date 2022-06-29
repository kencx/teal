package mock

import (
	"github.com/kencx/teal"
)

type BookStore struct {
	GetAllBooksFn    func() ([]*teal.Book, error)
	GetBookFn        func(id int64) (*teal.Book, error)
	GetBookByISBNFn  func(isbn string) (*teal.Book, error)
	GetBookByTitleFn func(title string) (*teal.Book, error)
	CreateBookFn     func(b *teal.Book) (*teal.Book, error)
	UpdateBookFn     func(id int64, b *teal.Book) (*teal.Book, error)
	DeleteBookFn     func(id int64) error
	GetByAuthorFn    func(name string) ([]*teal.Book, error)
}

type AuthorStore struct {
	GetAuthorFn     func(id int64) (*teal.Author, error)
	GetAllAuthorsFn func() ([]*teal.Author, error)
	CreateAuthorFn  func(a *teal.Author) (*teal.Author, error)
	UpdateAuthorFn  func(id int64, a *teal.Author) (*teal.Author, error)
	DeleteAuthorFn  func(id int64) error
}

func (m *BookStore) Get(id int64) (*teal.Book, error) {
	return m.GetBookFn(id)
}

func (m *BookStore) GetByISBN(isbn string) (*teal.Book, error) {
	return m.GetBookByISBNFn(isbn)
}

func (m *BookStore) GetByTitle(title string) (*teal.Book, error) {
	return m.GetBookByTitleFn(title)
}

func (m *BookStore) GetAll() ([]*teal.Book, error) {
	return m.GetAllBooksFn()
}

func (m *BookStore) Create(b *teal.Book) (*teal.Book, error) {
	return m.CreateBookFn(b)
}

func (m *BookStore) Update(id int64, b *teal.Book) (*teal.Book, error) {
	return m.UpdateBookFn(id, b)
}

func (m *BookStore) Delete(id int64) error {
	return m.DeleteBookFn(id)
}

func (m *BookStore) GetByAuthor(name string) ([]*teal.Book, error) {
	return m.GetByAuthorFn(name)
}

func (m *AuthorStore) Get(id int64) (*teal.Author, error) {
	return m.GetAuthorFn(id)

}

func (m *AuthorStore) GetAll() ([]*teal.Author, error) {
	return m.GetAllAuthorsFn()

}

func (m *AuthorStore) Create(a *teal.Author) (*teal.Author, error) {
	return m.CreateAuthorFn(a)

}

func (m *AuthorStore) Update(id int64, a *teal.Author) (*teal.Author, error) {
	return m.UpdateAuthorFn(id, a)

}

func (m *AuthorStore) Delete(id int64) error {
	return m.DeleteAuthorFn(id)
}
