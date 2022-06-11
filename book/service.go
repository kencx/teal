package book

import (
	"context"
	"fmt"

	teal "github.com/kencx/teal"
)

type Store interface {
	RetrieveBookWithID(id int) (*teal.Book, error)
	RetrieveBookWithTitle(title string) (*teal.Book, error)
	RetrieveAllBooks() ([]*teal.Book, error)
	CreateBook(ctx context.Context, b *teal.Book) (*teal.Book, error)
	UpdateBook(ctx context.Context, id int, b *teal.Book) error
	DeleteBook(ctx context.Context, id int) error
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

	b, err := s.db.RetrieveBookWithID(id)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) GetByTitle(title string) (*teal.Book, error) {
	// input validate title

	b, err := s.db.RetrieveBookWithTitle(title)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) GetAll() ([]*teal.Book, error) {
	b, err := s.db.RetrieveAllBooks()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) Create(ctx context.Context, b *teal.Book) (*teal.Book, error) {
	// validate b

	book, err := s.db.CreateBook(ctx, b)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("svc: invalid id %d", id)
	}

	err := s.db.DeleteBook(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
