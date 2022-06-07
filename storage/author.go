package storage

import (
	"database/sql"
	"fmt"

	teal "github.com/kencx/teal"
)

func (r *Store) GetAuthor(id int) (*teal.Author, error) {
	var a teal.Author

	stmt, err := r.db.Prepare("SELECT id, name FROM book WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("prepare stmt failed: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&a.ID, &a.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("query row failed: %w", err)
	}

	return &a, nil
}

func (r *Store) GetAllAuthors() ([]*teal.Author, error) {

	return nil, nil
}

func (r *Store) CreateAuthor(a *teal.Author) (int, error) {
	return 0, nil
}

func (r *Store) UpdateAuthor(id int) error {
	return nil
}
func (r *Store) DeleteAuthor(id int) error {
	return nil
}
