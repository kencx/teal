package pkg

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title" validate:"required"`
	Description   string `json:"description,omitempty"`
	Author        string `json:"author" validate:"required"`
	ISBN          string `json:"isbn" validate:"required,isbn"`
	NumOfPages    int    `json:"num_of_pages"`
	Rating        int    `json:"rating"`
	State         string `json:"state"` // default empty
	Read          string `json:"read"`  // default unread
	DateAdded     string `json:"-"`
	DateUpdated   string `json:"-"`
	DateCompleted string `json:"-"`
}

func (b Book) String() string {
	return fmt.Sprintf("id=%d title=%s author=%s isbn=%s", b.ID, b.Title, b.Author, b.ISBN)
}

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
