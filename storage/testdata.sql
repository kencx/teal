DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS books_authors;

CREATE TABLE books (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	title         text NOT NULL,
	description   text NULL,
	isbn          text NOT NULL UNIQUE,
	numOfPages    INTEGER,
	rating        text,
	state         text,
	dateAdded     text,
	dateUpdated   text,
	dateCompleted text
);

CREATE TABLE authors (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name text NOT NULL
);

CREATE TABLE books_authors (
	book_id INT REFERENCES books(id),
	author_id INT REFERENCES authors(id),
	PRIMARY KEY(book_id, author_id)
);

-- ORDER MATTERS, APPEND NEW DATA TO END
-- book 1
INSERT INTO books (
	title, isbn, numOfPages, rating, state
) VALUES ("Leviathan Wakes", "9999", 250, 5, "read");

INSERT INTO authors (
	name
) VALUES ("S.A. Corey");

INSERT INTO books_authors (
	book_id, author_id
	) VALUES (
	(SELECT id FROM books WHERE title = "Leviathan Wakes"),
	(SELECT id FROM authors WHERE name = "S.A. Corey")
);

-- book 2
INSERT INTO books (
	title, isbn, numOfPages, rating, state
) VALUES ("Red Rising", "1234", 900, 4, "unread");

INSERT INTO authors (
	name
) VALUES ("Pierce Brown");

INSERT INTO books_authors (
	book_id, author_id
	) VALUES (
	(SELECT id FROM books WHERE title = "Red Rising"),
	(SELECT id FROM authors WHERE name = "Pierce Brown")
);

-- book 3 multiple authors
INSERT INTO books (
	title, isbn
) VALUES ("Many Authors", "56789");

INSERT INTO authors (
	name
) VALUES
	("John Doe"),
	("Regina Phallange"),
	("Ken Adams");

INSERT INTO books_authors (
	book_id, author_id
	) VALUES
	((SELECT id FROM books WHERE title = "Many Authors"), (SELECT id FROM authors WHERE name = "John Doe")),
	((SELECT id FROM books WHERE title = "Many Authors"), (SELECT id FROM authors WHERE name = "Regina Pallange")),
	((SELECT id FROM books WHERE title = "Many Authors"), (SELECT id FROM authors WHERE name = "Ken Adams"));

