package mock

import (
	"context"

	"github.com/kencx/teal"
)

type Store struct {
	GetAllBooksFn    func() ([]*teal.Book, error)
	GetBookFn        func(id int) (*teal.Book, error)
	GetBookByTitleFn func(title string) (*teal.Book, error)
	CreateBookFn     func(ctx context.Context, b *teal.Book) (*teal.Book, error)
	UpdateBookFn     func(ctx context.Context, id int, b *teal.Book) (*teal.Book, error)
	DeleteBookFn     func(ctx context.Context, id int) error

	GetAuthorFn     func(id int) (*teal.Author, error)
	GetAllAuthorsFn func() ([]*teal.Author, error)
	CreateAuthorFn  func(ctx context.Context, a *teal.Author) (*teal.Author, error)
	UpdateAuthorFn  func(ctx context.Context, id int, a *teal.Author) (*teal.Author, error)
	DeleteAuthorFn  func(ctx context.Context, id int) error
}

func (m *Store) RetrieveAllBooks() ([]*teal.Book, error) {
	return m.GetAllBooksFn()
}

func (m *Store) RetrieveBookWithID(id int) (*teal.Book, error) {
	return m.GetBookFn(id)
}

func (m *Store) RetrieveBookWithTitle(title string) (*teal.Book, error) {
	return m.GetBookByTitleFn(title)
}

func (m *Store) CreateBook(ctx context.Context, b *teal.Book) (*teal.Book, error) {
	return m.CreateBookFn(ctx, b)
}

func (m *Store) UpdateBook(ctx context.Context, id int, b *teal.Book) (*teal.Book, error) {
	return m.UpdateBookFn(ctx, id, b)
}

func (m *Store) DeleteBook(ctx context.Context, id int) error {
	return m.DeleteBookFn(ctx, id)
}

func (m *Store) RetrieveAuthorWithID(id int) (*teal.Author, error) {
	return m.GetAuthorFn(id)

}

func (m *Store) RetrieveAllAuthors() ([]*teal.Author, error) {
	return m.GetAllAuthorsFn()

}

func (m *Store) CreateAuthor(ctx context.Context, a *teal.Author) (*teal.Author, error) {
	return m.CreateAuthorFn(ctx, a)

}

func (m *Store) UpdateAuthor(ctx context.Context, id int, a *teal.Author) (*teal.Author, error) {
	return m.UpdateAuthorFn(ctx, id, a)

}

func (m *Store) DeleteAuthor(ctx context.Context, id int) error {
	return m.DeleteAuthorFn(ctx, id)
}
