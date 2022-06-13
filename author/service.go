package author

import (
	"context"
	"fmt"

	"github.com/kencx/teal"
)

type Store interface {
	RetrieveAuthorWithID(id int) (*teal.Author, error)
	RetrieveAllAuthors() ([]*teal.Author, error)
	CreateAuthor(ctx context.Context, b *teal.Author) (*teal.Author, error)
	UpdateAuthor(ctx context.Context, id int, b *teal.Author) (*teal.Author, error)
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

func (s *Service) Create(ctx context.Context, a *teal.Author) (*teal.Author, error) {

	res, err := s.db.CreateAuthor(ctx, a)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) Update(ctx context.Context, id int, a *teal.Author) (*teal.Author, error) {
	res, err := s.db.UpdateAuthor(ctx, id, a)
	if err != nil {
		return nil, err
	}
	return res, nil
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
