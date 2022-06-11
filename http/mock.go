package http

import (
	"context"

	"github.com/kencx/teal"
)

type mockBookService struct {
	getAllBooksFn    func() ([]*teal.Book, error)
	getBookFn        func(id int) (*teal.Book, error)
	getBookByTitleFn func(title string) (*teal.Book, error)
	createBookFn     func(ctx context.Context, b *teal.Book) (*teal.Book, error)
	updateBookFn     func(ctx context.Context, id int, b *teal.Book) (*teal.Book, error)
	deleteBookFn     func(ctx context.Context, id int) error
}

func (m *mockBookService) GetAll() ([]*teal.Book, error) {
	return m.getAllBooksFn()
}

func (m *mockBookService) Get(id int) (*teal.Book, error) {
	return m.getBookFn(id)
}

func (m *mockBookService) GetByTitle(title string) (*teal.Book, error) {
	return m.getBookByTitleFn(title)
}

func (m *mockBookService) Create(ctx context.Context, b *teal.Book) (*teal.Book, error) {
	return m.createBookFn(ctx, b)
}

func (m *mockBookService) Update(ctx context.Context, id int, b *teal.Book) (*teal.Book, error) {
	return m.updateBookFn(ctx, id, b)
}

func (m *mockBookService) Delete(ctx context.Context, id int) error {
	return m.deleteBookFn(ctx, id)
}
