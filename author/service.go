package author

import (
	"context"
	"fmt"

	"github.com/kencx/teal"
)

type Store interface {
	RetrieveAuthorWithID(id int) (*teal.Author, error)
	RetrieveAllAuthors() ([]*teal.Author, error)
	CreateAuthor(ctx context.Context, b *teal.Author) error
	// UpdateAuthor(id int, b *teal.Author) error
	DeleteAuthor(ctx context.Context, id int) error
}

type Service struct {
	db Store
}

func NewService(db Store) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Get(id int) (*teal.Author, error) {
	if id <= 0 {
		return nil, fmt.Errorf("svc: invalid id %d", id)
	}

	b, err := s.db.RetrieveAuthorWithID(id)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) GetAll() ([]*teal.Author, error) {
	b, err := s.db.RetrieveAllAuthors()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) Create(ctx context.Context, b *teal.Author) (int, error) {
	// validate b

	err := s.db.CreateAuthor(ctx, b)
	if err != nil {
		return -1, err
	}

	return 0, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("svc: invalid id %d", id)
	}

	err := s.db.DeleteAuthor(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
