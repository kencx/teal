package teal

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Book struct {
	ID            int            `json:"id" db:"id"`
	Title         string         `json:"title" validate:"required" db:"title"`
	Description   sql.NullString `json:"description,omitempty" db:"description"`
	Author        []string       `json:"author" validate:"required"`
	ISBN          string         `json:"isbn" validate:"required,isbn" db:"isbn"`
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

var ErrBookNotFound = fmt.Errorf("Book not found")

func (b *Book) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("isbn", validateISBN)
	return validate.Struct(b)
}

func validateISBN(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[0-9]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	return len(matches) == 1
}
