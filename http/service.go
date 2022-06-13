package http

import (
	"context"

	"github.com/kencx/teal"
)

type BookService interface {
	Get(id int) (*teal.Book, error)
	GetByTitle(title string) (*teal.Book, error)
	GetAll() ([]*teal.Book, error)
	Create(ctx context.Context, b *teal.Book) (*teal.Book, error)
	Update(ctx context.Context, id int, b *teal.Book) (*teal.Book, error)
	Delete(ctx context.Context, id int) error

	GetByAuthor(name string) ([]*teal.Book, error)
}

type AuthorService interface {
	Get(id int) (*teal.Author, error)
	GetAll() ([]*teal.Author, error)
	Create(ctx context.Context, b *teal.Author) (*teal.Author, error)
	Update(ctx context.Context, id int, b *teal.Author) (*teal.Author, error)
	Delete(ctx context.Context, id int) error
}
