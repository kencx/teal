package teal

import "time"

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`

	LastLogin time.Time `json:"last_login,omitempty"`
	Token     string    `json:"-"`
	Role      string    `json:"role,omitempty"`
}

func (u *User) ChangePassword(hash string) {
	u.Password = hash
}

func (u *User) UpdateLastLogin() {
	u.LastLogin = time.Now()
}
