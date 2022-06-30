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
	GetAuthorFn       func(id int64) (*teal.Author, error)
	GetAuthorByNameFn func(name string) (*teal.Author, error)
	GetAllAuthorsFn   func() ([]*teal.Author, error)
	CreateAuthorFn    func(a *teal.Author) (*teal.Author, error)
	UpdateAuthorFn    func(id int64, a *teal.Author) (*teal.Author, error)
	DeleteAuthorFn    func(id int64) error
}

type UserStore struct {
	GetUserFn           func(id int64) (*teal.User, error)
	GetUserByUsernameFn func(username string) (*teal.User, error)
	GetAllUsersFn       func() ([]*teal.User, error)
	CreateUserFn        func(a *teal.User) (*teal.User, error)
	UpdateUserFn        func(id int64, a *teal.User) (*teal.User, error)
	DeleteUserFn        func(id int64) error
}

func (s *BookStore) Get(id int64) (*teal.Book, error) {
	return s.GetBookFn(id)
}

func (s *BookStore) GetByISBN(isbn string) (*teal.Book, error) {
	return s.GetBookByISBNFn(isbn)
}

func (s *BookStore) GetByTitle(title string) (*teal.Book, error) {
	return s.GetBookByTitleFn(title)
}

func (s *BookStore) GetAll() ([]*teal.Book, error) {
	return s.GetAllBooksFn()
}

func (s *BookStore) Create(b *teal.Book) (*teal.Book, error) {
	return s.CreateBookFn(b)
}

func (s *BookStore) Update(id int64, b *teal.Book) (*teal.Book, error) {
	return s.UpdateBookFn(id, b)
}

func (s *BookStore) Delete(id int64) error {
	return s.DeleteBookFn(id)
}

func (s *BookStore) GetByAuthor(name string) ([]*teal.Book, error) {
	return s.GetByAuthorFn(name)
}

func (s *AuthorStore) Get(id int64) (*teal.Author, error) {
	return s.GetAuthorFn(id)
}

func (s *AuthorStore) GetByName(name string) (*teal.Author, error) {
	return s.GetAuthorByNameFn(name)
}

func (s *AuthorStore) GetAll() ([]*teal.Author, error) {
	return s.GetAllAuthorsFn()
}

func (s *AuthorStore) Create(a *teal.Author) (*teal.Author, error) {
	return s.CreateAuthorFn(a)
}

func (s *AuthorStore) Update(id int64, a *teal.Author) (*teal.Author, error) {
	return s.UpdateAuthorFn(id, a)
}

func (s *AuthorStore) Delete(id int64) error {
	return s.DeleteAuthorFn(id)
}

func (s *UserStore) Get(id int64) (*teal.User, error) {
	return s.GetUserFn(id)
}

func (s *UserStore) GetByUsername(username string) (*teal.User, error) {
	return s.GetUserByUsernameFn(username)
}

func (s *UserStore) Create(a *teal.User) (*teal.User, error) {
	return s.CreateUserFn(a)

}

func (s *UserStore) Update(id int64, a *teal.User) (*teal.User, error) {
	return s.UpdateUserFn(id, a)

}

func (s *UserStore) Delete(id int64) error {
	return s.DeleteUserFn(id)
}
