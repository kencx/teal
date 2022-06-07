package author

import (
	"fmt"

	"github.com/kencx/teal"
)

type Store interface {
	GetAuthor(id int) (*teal.Author, error)
	GetAllAuthors() ([]*teal.Author, error)
	CreateAuthor(b *teal.Author) (int, error)
	// UpdateAuthor(id int, b *teal.Author) error
	DeleteAuthor(id int) error
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

	b, err := s.db.GetAuthor(id)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) GetAll() ([]*teal.Author, error) {
	b, err := s.db.GetAllAuthors()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) Create(b *teal.Author) (int, error) {
	// validate b

	id, err := s.db.CreateAuthor(b)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *Service) Delete(id int) error {
	if id <= 0 {
		return fmt.Errorf("svc: invalid id %d", id)
	}

	err := s.db.DeleteAuthor(id)
	if err != nil {
		return err
	}

	return nil
}
