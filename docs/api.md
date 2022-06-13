# API

## Resources

### Books

#### List

```
GET /api/books/
```

List all books

Parameters:
- author - Filters books by a specific author

Example response:
```json
[
  {
    "id": 1,
    "title": "Foobar",
    "description": "",
    "author": [
      "John Doe",
      "Jane Doe"
    ],
    "isbn": "",
    "num_of_pages": 100,
    "rating": 5,
    "state": "unread"
  }
]
```

```
GET /api/books/[id]/
```

Retrieves a single book by ID.

#### Create

```
POST /api/books/
```

Create a single book with given payload.

Example payload:
```json
{
  "title": "Foobar",
  "description": "",
  "author": [
  	"John Doe",
  	"Jane Doe"
  ],
  "isbn": "1",
  "num_of_pages": 100,
  "rating": 5,
  "state": "unread"
}
```

#### Update

```
PUT /api/books/[id]/
```

Update a single book by ID.

#### Delete

```
DELETE /api/books/[id]/
```

Delete a single book by ID.

### Authors

#### List

```
GET /api/authors/
```

List all authors

```
GET /api/authors/[id]/
```

Retrieve a single author by ID.

#### Create

```
POST /api/authors/
```

Create an author.

Example payload:
```json
{
  "name": "John Doe"
}
```

#### Update

```
PUT /api/authors/[id]/
```

Update a single author by ID.

#### Delete

```
DELETE /api/authors/[id]/
```

Delete a single author by ID.
