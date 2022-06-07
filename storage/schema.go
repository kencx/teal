package storage

import teal "github.com/kencx/teal"

const (
	DROP_ALL = `DROP TABLE IF EXISTS books;
				DROP TABLE IF EXISTS authors;
				DROP TABLE IF EXISTS books_authors;`

	CREATE_TABLES = `CREATE TABLE IF NOT EXISTS books (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		title         TEXT NOT NULL,
		description   TEXT NULL,
		isbn          TEXT NOT NULL UNIQUE,
		numOfPages    INTEGER,
		rating        TEXT,
		state         TEXT,
		dateAdded     INTEGER,
		dateUpdated   INTEGER,
		dateCompleted INTEGER
	);

	CREATE TABLE IF NOT EXISTS authors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS books_authors (
		book_id INTEGER REFERENCES books(id),
		author_id INTEGER REFERENCES authors(id),
		PRIMARY KEY(book_id, author_id)
	);`
)

var (
	testBook1 = &teal.Book{
		ID:         1,
		Title:      "Leviathan Wakes",
		ISBN:       "9999",
		NumOfPages: 250,
		Rating:     5,
		State:      "read",
	}
)

// func (r *Store) NewCreateBook(b *teal.Book) error {
//
// 	book_id, err := createBook(r, b)
// 	if err != nil {
// 		return err
// 	}
//
// 	var author_ids []int
// 	for _, a := range b.Author {
// 		id, err := createAuthor(r, &a)
// 		if err != nil {
// 			return err
// 		}
// 		author_ids = append(author_ids, id)
// 	}
//
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer endTx(tx, err)
//
// 	for _, id := range author_ids {
// 		_, err = tx.Exec(`INSERT INTO books_authors (book_id, author_id) VALUES ($1, $2)`, book_id, id)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
