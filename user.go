package teal

import (
	"time"

	"github.com/kencx/teal/validator"
)

type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	HashedPassword []byte `json:"-"`
	Email          string `json:"email"`

	LastLogin time.Time `json:"last_login,omitempty"`
	Token     string    `json:"-"`
	Role      string    `json:"role,omitempty"`
}

// func (u *User) ChangePassword(hash string) {
// 	u.HashedPassword = hash
// }

func (u *User) UpdateLastLogin() {
	u.LastLogin = time.Now()
}

func (u *User) Validate(v *validator.Validator) {
	v.Check(u.Name != "", "name", "value is missing")
	v.Check(u.Username != "", "username", "value is missing")
	// TODO check password before hashing
	v.Check(u.Email != "", "email", "value is missing")
}
