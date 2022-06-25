package mock

import (
	"context"

	"github.com/kencx/teal"
)

type BookStore struct {
	GetAllBooksFn    func() ([]*teal.Book, error)
	GetBookFn        func(id int) (*teal.Book, error)
	GetBookByTitleFn func(title string) (*teal.Book, error)
	CreateBookFn     func(ctx context.Context, b *teal.Book) (*teal.Book, error)
	UpdateBookFn     func(ctx context.Context, id int, b *teal.Book) (*teal.Book, error)
	DeleteBookFn     func(ctx context.Context, id int) error
	GetByAuthorFn    func(name string) ([]*teal.Book, error)
}

type AuthorStore struct {
	GetAuthorFn     func(id int) (*teal.Author, error)
	GetAllAuthorsFn func() ([]*teal.Author, error)
	CreateAuthorFn  func(ctx context.Context, a *teal.Author) (*teal.Author, error)
	UpdateAuthorFn  func(ctx context.Context, id int, a *teal.Author) (*teal.Author, error)
	DeleteAuthorFn  func(ctx context.Context, id int) error
}

func (m *BookStore) Get(id int) (*teal.Book, error) {
	return m.GetBookFn(id)
}

func (m *BookStore) GetByTitle(title string) (*teal.Book, error) {
	return m.GetBookByTitleFn(title)
}

func (m *BookStore) GetAll() ([]*teal.Book, error) {
	return m.GetAllBooksFn()
}

func (m *BookStore) Create(ctx context.Context, b *teal.Book) (*teal.Book, error) {
	return m.CreateBookFn(ctx, b)
}

func (m *BookStore) Update(ctx context.Context, id int, b *teal.Book) (*teal.Book, error) {
	return m.UpdateBookFn(ctx, id, b)
}

func (m *BookStore) Delete(ctx context.Context, id int) error {
	return m.DeleteBookFn(ctx, id)
}

func (m *BookStore) GetByAuthor(name string) ([]*teal.Book, error) {
	return m.GetByAuthorFn(name)
}

func (m *AuthorStore) Get(id int) (*teal.Author, error) {
	return m.GetAuthorFn(id)

}

func (m *AuthorStore) GetAll() ([]*teal.Author, error) {
	return m.GetAllAuthorsFn()

}

func (m *AuthorStore) Create(ctx context.Context, a *teal.Author) (*teal.Author, error) {
	return m.CreateAuthorFn(ctx, a)

}

func (m *AuthorStore) Update(ctx context.Context, id int, a *teal.Author) (*teal.Author, error) {
	return m.UpdateAuthorFn(ctx, id, a)

}

func (m *AuthorStore) Delete(ctx context.Context, id int) error {
	return m.DeleteAuthorFn(ctx, id)
}
