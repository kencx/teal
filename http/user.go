package http

import (
	"github.com/kencx/teal"
)

type UserStore interface {
	Get(id int64) (*teal.User, error)
	GetByUsername(username string) (*teal.User, error)
	GetAll() ([]*teal.User, error)
	Create(b *teal.User) (*teal.User, error)
	Update(id int64, b *teal.User) (*teal.User, error)
	Delete(id int64) error
}
