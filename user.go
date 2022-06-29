package teal

import (
	"database/sql"
	"errors"
	"time"

	"github.com/kencx/teal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int64        `json:"id"`
	Name           string       `json:"name"`
	Username       string       `json:"username"`
	HashedPassword []byte       `json:"-"`
	Role           string       `json:"role"`
	LastLogin      sql.NullTime `json:"-"`
	DateAdded      time.Time    `json:"-"`
}

// Destination struct for POST user requests
type InputUser struct {
	Name     string
	Username string
	Password string
	Role     string
}

func (u *User) SetPassword(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), 12)
	if err != nil {
		return err
	}

	u.HashedPassword = hash
	return nil
}

func (u *User) PasswordMatches(text string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(text))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (u *InputUser) Validate(v *validator.Validator) {
	v.Check(u.Name != "", "name", "value is missing")
	v.Check(u.Username != "", "username", "value is missing")
	v.Check(u.Password != "", "password", "value is missing")
	v.Check(len(u.Password) >= 8, "password", "must be at least 8 chars long")
}
