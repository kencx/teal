package book

import (
	"fmt"

	teal "github.com/kencx/teal"
)

type Store interface {
	GetBook(id int) (*teal.Book, error)
	GetBookByTitle(title string) (*teal.Book, error)
	GetAllBooks() ([]*teal.Book, error)
	CreateBook(b *teal.Book) (int, error)
	UpdateBook(id int, b *teal.Book) error
	DeleteBook(id int) error
}

type Service struct {
	db Store
}

func NewService(db Store) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Get(id int) (*teal.Book, error) {
	if id <= 0 {
		return nil, fmt.Errorf("svc: invalid id %d", id)
	}

	b, err := s.db.GetBook(id)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) GetByTitle(title string) (*teal.Book, error) {
	// input validate title

	b, err := s.db.GetBookByTitle(title)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) GetAll() ([]*teal.Book, error) {
	b, err := s.db.GetAllBooks()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) Create(b *teal.Book) (int, error) {
	// validate b

	id, err := s.db.CreateBook(b)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *Service) Delete(id int) error {
	if id <= 0 {
		return fmt.Errorf("svc: invalid id %d", id)
	}

	err := s.db.DeleteBook(id)
	if err != nil {
		return err
	}

	return nil
}
