package http

import (
	"context"

	"github.com/kencx/teal"
)

type Store interface {
	RetrieveBookWithID(id int) (*teal.Book, error)
	RetrieveBookWithTitle(title string) (*teal.Book, error)
	RetrieveAllBooks() ([]*teal.Book, error)
	CreateBook(ctx context.Context, b *teal.Book) (*teal.Book, error)
	UpdateBook(ctx context.Context, id int, b *teal.Book) (*teal.Book, error)
	DeleteBook(ctx context.Context, id int) error

	RetrieveAuthorWithID(id int) (*teal.Author, error)
	RetrieveAllAuthors() ([]*teal.Author, error)
	CreateAuthor(ctx context.Context, b *teal.Author) (*teal.Author, error)
	UpdateAuthor(ctx context.Context, id int, b *teal.Author) (*teal.Author, error)
	DeleteAuthor(ctx context.Context, id int) error
}
