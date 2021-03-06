-- ORDER MATTERS, APPEND NEW DATA TO END
-- remember to add structs to testdata.go as well

-- book 1
INSERT INTO books (
	title, isbn, numOfPages, rating, state
) VALUES ("Leviathan Wakes", "1", 250, 5, "read");

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
) VALUES ("Red Rising", "2", 900, 4, "unread");

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
) VALUES ("Many Authors", "3");

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
	((SELECT id FROM books WHERE title = "Many Authors"), (SELECT id FROM authors WHERE name = "Regina Phallange")),
	((SELECT id FROM books WHERE title = "Many Authors"), (SELECT id FROM authors WHERE name = "Ken Adams"));

-- book 4 existing author
INSERT INTO books (
	title, isbn
) VALUES ("New Book", "4");

INSERT INTO books_authors (
	book_id, author_id
	) VALUES
	((SELECT id FROM books WHERE title = "New Book"), (SELECT id FROM authors WHERE name = "John Doe"));

-- user 1, 2
INSERT INTO users (
	name, username, hashed_password
) VALUES
	("John Doe", "johndoe", "abc123456789"),
	("Ben Adams", "benadams", "abc123456789");

