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
	HashedPassword password     `json:"-"`
	Email          string       `json:"email"`
	LastLogin      sql.NullTime `json:"last_login,omitempty"`
	Role           string       `json:"role,omitempty"`
	DateAdded      time.Time    `json:"-"`
}

type password struct {
	Text *string
	Hash []byte `db:"hashed_password"`
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), 12)
	if err != nil {
		return err
	}

	p.Text = &text
	p.Hash = hash
	return nil
}

func (p *password) Matches(text string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(text))

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

func (u *User) Validate(v *validator.Validator) {
	v.Check(u.Name != "", "name", "value is missing")
	v.Check(u.Username != "", "username", "value is missing")
	v.Check(u.Email != "", "email", "value is missing")

	if u.HashedPassword.Text != nil {
		ValidatePasswordText(v, *u.HashedPassword.Text)
	}

	if u.HashedPassword.Hash == nil {
		panic("missing password hash for user")
	}
}

func ValidatePasswordText(v *validator.Validator, password string) {
	v.Check(password != "", "password", "value is missing")
	v.Check(len(password) >= 8, "password", "must be at least 8 chars long")
}
