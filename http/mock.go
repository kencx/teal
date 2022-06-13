package http

import (
	"context"

	"github.com/kencx/teal"
)

type mockStore struct {
	getAllBooksFn    func() ([]*teal.Book, error)
	getBookFn        func(id int) (*teal.Book, error)
	getBookByTitleFn func(title string) (*teal.Book, error)
	createBookFn     func(ctx context.Context, b *teal.Book) (*teal.Book, error)
	updateBookFn     func(ctx context.Context, id int, b *teal.Book) (*teal.Book, error)
	deleteBookFn     func(ctx context.Context, id int) error
}

func (m *mockStore) RetrieveAllBooks() ([]*teal.Book, error) {
	return m.getAllBooksFn()
}

func (m *mockStore) RetrieveBookWithID(id int) (*teal.Book, error) {
	return m.getBookFn(id)
}

func (m *mockStore) RetrieveBookWithTitle(title string) (*teal.Book, error) {
	return m.getBookByTitleFn(title)
}

func (m *mockStore) CreateBook(ctx context.Context, b *teal.Book) (*teal.Book, error) {
	return m.createBookFn(ctx, b)
}

func (m *mockStore) UpdateBook(ctx context.Context, id int, b *teal.Book) (*teal.Book, error) {
	return m.updateBookFn(ctx, id, b)
}

func (m *mockStore) DeleteBook(ctx context.Context, id int) error {
	return m.deleteBookFn(ctx, id)
}

func (m *mockStore) RetrieveAuthorWithID(id int) (*teal.Author, error) {
	return nil, nil

}

func (m *mockStore) RetrieveAllAuthors() ([]*teal.Author, error) {
	return nil, nil

}

func (m *mockStore) CreateAuthor(ctx context.Context, b *teal.Author) error {
	return nil

}

func (m *mockStore) UpdateAuthor(ctx context.Context, id int, b *teal.Author) (*teal.Author, error) {
	return nil, nil

}

func (m *mockStore) DeleteAuthor(ctx context.Context, id int) error {
	return nil
}
