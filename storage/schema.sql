DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS books_authors;
DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS books (
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

-- create state enum

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
	email VARCHAR(255) NOT NULL UNIQUE,
	role TEXT NOT NULL DEFAULT "user",
	dateAdded TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
);

-- create role enum
