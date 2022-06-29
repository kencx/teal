package storage

import (
	teal "github.com/kencx/teal"
)

const (
	DROP_ALL = `DROP TABLE IF EXISTS books;
				DROP TABLE IF EXISTS authors;
				DROP TABLE IF EXISTS books_authors;
				DROP TABLE IF EXISTS users;
				DROP TABLE IF EXISTS apikeys;`

	CREATE_TABLES = `CREATE TABLE IF NOT EXISTS books (
			id            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			title         TEXT NOT NULL,
			description   TEXT,
			isbn          TEXT NOT NULL UNIQUE,
			numOfPages    INTEGER DEFAULT 0,
			rating        INTEGER DEFAULT 0,
			state         TEXT NOT NULL DEFAULT "unread",
			dateAdded     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			dateUpdated   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			dateCompleted TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS authors (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS books_authors (
			book_id INTEGER REFERENCES books(id),
			author_id INTEGER REFERENCES authors(id),
			PRIMARY KEY(book_id, author_id)
		);

		CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			username VARCHAR(255) NOT NULL UNIQUE,
			hashed_password CHAR(60) NOT NULL,
			role TEXT NOT NULL DEFAULT "user",
			dateAdded TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS apikeys (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL,
			hashed_key TEXT NOT NULL UNIQUE,
			dateAdded TIMESTAMP NOT NULL,
			dateExpired TIMESTAMP NOT NULL
		);`
)

// structs here are in testdata.sql
var (
	allBooks = []*teal.Book{testBook1, testBook2, testBook3, testBook4}

	testBook1 = &teal.Book{
		ID:         1,
		Title:      "Leviathan Wakes",
		ISBN:       "1",
		NumOfPages: 250,
		Rating:     5,
		State:      "read",
		Author:     []string{"S.A. Corey"},
	}
	testBook2 = &teal.Book{
		ID:         2,
		Title:      "Red Rising",
		ISBN:       "2",
		NumOfPages: 900,
		Rating:     4,
		State:      "unread",
		Author:     []string{"Pierce Brown"},
	}
	testBook3 = &teal.Book{
		ID:     3,
		Title:  "Many Authors",
		ISBN:   "3",
		State:  "unread",
		Author: []string{"John Doe", "Regina Phallange", "Ken Adams"},
	}
	testBook4 = &teal.Book{
		ID:     4,
		Title:  "New Book",
		ISBN:   "4",
		State:  "unread",
		Author: []string{"John Doe"},
	}

	testAuthor1 = &teal.Author{
		ID:   1,
		Name: "S.A. Corey",
	}
	testAuthor2 = &teal.Author{
		ID:   2,
		Name: "Pierce Brown",
	}
	testAuthor3 = &teal.Author{
		ID:   3,
		Name: "John Doe",
	}
	testAuthor4 = &teal.Author{
		ID:   4,
		Name: "Regina Phallange",
	}
	testAuthor5 = &teal.Author{
		ID:   5,
		Name: "Ken Adams",
	}

	testUser1 = &teal.User{
		ID:             1,
		Name:           "John Doe",
		Username:       "johndoe",
		HashedPassword: []byte("abc123456789"),
		Role:           "user",
	}

	testUser2 = &teal.User{
		ID:             2,
		Name:           "Ben Adams",
		Username:       "benadams",
		HashedPassword: []byte("abc123456789"),
		Role:           "user",
	}
)
