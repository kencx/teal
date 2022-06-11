package teal

import (
	"database/sql"
	"regexp"
)

type Book struct {
	ID            int            `json:"id" db:"id"`
	Title         string         `json:"title" db:"title"`
	Description   sql.NullString `json:"description,omitempty" db:"description"`
	Author        []string       `json:"author"`
	ISBN          string         `json:"isbn" db:"isbn"`
	NumOfPages    int            `json:"num_of_pages" db:"numOfPages"`
	Rating        int            `json:"rating" db:"rating"`
	State         string         `json:"state" db:"state"` // default empty
	DateAdded     sql.NullTime   `json:"-" db:"dateAdded"`
	DateUpdated   sql.NullTime   `json:"-" db:"dateUpdated"`
	DateCompleted sql.NullTime   `json:"-" db:"dateCompleted"`
}

// func (b Book) String() string {
// 	return fmt.Sprintf(`id=%d title=%s desc=%s author=%v isbn=%s dateAdded=%s dateUpdated=%s dateCompleted=%s`,
// 		b.ID, b.Title, b.Description, b.Author, b.ISBN, b.DateAdded, b.DateUpdated, b.DateCompleted)
// }

func (b *Book) Validate() []*ValidationError {
	var verrs []*ValidationError

	if b.Title == "" {
		verrs = append(verrs, NewValidationError("title", "value is missing"))
	}
	if verr := validateAuthor(b.Author); verr != nil {
		verrs = append(verrs, verr)
	}
	if verr := validateISBN(b.ISBN); verr != nil {
		verrs = append(verrs, verr)
	}
	return verrs
}

func validateAuthor(authors []string) *ValidationError {
	if len(authors) == 0 {
		return NewValidationError("author", "value is missing")
	}

	if len(authors) == 1 && authors[0] == "" {
		return NewValidationError("author", "value is missing")
	}

	return nil
}

func validateISBN(isbn string) *ValidationError {
	if isbn == "" {
		return NewValidationError("isbn", "value is missing")
	}

	// TODO isbn regex
	re := regexp.MustCompile(`[0-9]+`)
	matches := re.FindAllString(isbn, -1)
	if len(matches) != 1 {
		return NewValidationError("isbn", "incorrect format")
	}
	return nil
}
