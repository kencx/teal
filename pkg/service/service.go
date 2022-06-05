package service

import (
	"github.com/kencx/teal/pkg"
)

type BookRepository interface {
	GetBook(id int) (*pkg.Book, error)
	GetBookByTitle(title string) (*pkg.Book, error)
	GetAllBooks() ([]*pkg.Book, error)
	CreateBook(b *pkg.Book) (int, error)
	UpdateBook(id int, b *pkg.Book) error
	DeleteBook(id int) error
}

type AuthorRepository interface {
	GetAuthor(id int) (*pkg.Author, error)
	GetAuthorByTitle(title string) (*pkg.Author, error)
	GetAllAuthors() ([]*pkg.Author, error)
	CreateAuthor(b *pkg.Author) (int, error)
	UpdateAuthor(id int, b *pkg.Author) error
	DeleteAuthor(id int) error
}

type BookService struct {
	db BookRepository
}

type AuthorService struct {
	db AuthorRepository
}

func NewService(b BookRepository) *BookService {
	return &BookService{
		db: b,
	}
}
