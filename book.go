package teal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/kencx/teal/validator"
)

type Book struct {
	ID            int64        `json:"id" db:"id"`
	Title         string       `json:"title" db:"title"`
	Description   NullString   `json:"description,omitempty" db:"description"`
	Author        []string     `json:"author"`
	ISBN          string       `json:"isbn" db:"isbn"`
	NumOfPages    int          `json:"num_of_pages" db:"numOfPages"`
	Rating        int          `json:"rating" db:"rating"`
	State         string       `json:"state" db:"state"`
	DateAdded     sql.NullTime `json:"-" db:"dateAdded"`
	DateUpdated   sql.NullTime `json:"-" db:"dateUpdated"`
	DateCompleted sql.NullTime `json:"-" db:"dateCompleted"`
}

func (b Book) String() string {
	return fmt.Sprintf(`[title=%s description=%s author=%v isbn=%s npages=%d rating=%d state=%s]`,
		b.Title, b.Description.String, b.Author, b.ISBN, b.NumOfPages, b.Rating, b.State)
}

type NullString struct {
	sql.NullString
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

func (n NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s != nil {
		n.Valid = true
		n.String = *s
	} else {
		n.Valid = false
	}
	return nil
}

var isbnRgx = regexp.MustCompile(`[0-9]+`)

func (b *Book) Validate(v *validator.Validator) {
	v.Check(b.Title != "", "title", "value is missing")

	v.Check(len(b.Author) != 0, "author", "value is missing")

	v.Check(b.ISBN != "", "isbn", "value is missing")
	v.Check(validator.Matches(b.ISBN, isbnRgx), "isbn", "incorrect format")

	v.Check(b.NumOfPages >= 0, "numOfPages", "must be >= 0")

	v.Check(b.Rating >= 0, "rating", "must be >= 0")
	v.Check(b.Rating <= 10, "rating", "must be <= 10")
}
