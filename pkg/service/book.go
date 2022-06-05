package service

import (
	"fmt"

	"github.com/kencx/teal/pkg"
)

func (s *BookService) GetBook(id int) (*pkg.Book, error) {
	if id <= 0 {
		return nil, fmt.Errorf("svc: invalid id %d", id)
	}

	b, err := s.db.GetBook(id)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *BookService) GetBookByTitle(title string) (*pkg.Book, error) {
	// input validate title

	b, err := s.db.GetBookByTitle(title)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *BookService) GetAllBooks() ([]*pkg.Book, error) {
	b, err := s.db.GetAllBooks()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *BookService) CreateBook(b *pkg.Book) (int, error) {
	// validate b

	id, err := s.db.CreateBook(b)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *BookService) UpdateBook(id int, b *pkg.Book) error {
	// validate b
	if id <= 0 {
		return fmt.Errorf("svc: invalid id %d", id)
	}

	err := s.db.UpdateBook(id, b)
	if err != nil {
		return err
	}

	return nil
}

func (s *BookService) DeleteBook(id int) error {
	if id <= 0 {
		return fmt.Errorf("svc: invalid id %d", id)
	}

	err := s.db.DeleteBook(id)
	if err != nil {
		return err
	}

	return nil
}
