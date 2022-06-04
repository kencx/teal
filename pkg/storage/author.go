package storage

import (
	"database/sql"
	"fmt"

	"github.com/kencx/teal/pkg"
)

func (db *DB) GetAuthor(id int) (*pkg.Author, error) {
	var a pkg.Author

	stmt, err := db.db.Prepare("SELECT id, name FROM book WHERE id = ?")
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

func (db *DB) GetAllAuthors() ([]*pkg.Author, error) {

	return nil, nil
}

func (db *DB) CreateAuthor(a *pkg.Author) (int, error) {
	return 0, nil
}

func (db *DB) UpdateAuthor(id int) error {
	return nil
}
func (db *DB) DeleteAuthor(id int) error {
	return nil
}
