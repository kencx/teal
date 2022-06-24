package http

import (
	"context"

	"github.com/kencx/teal"
)

type UserStore interface {
	Get(id int) (*teal.User, error)
	GetByUsername(username string) (*teal.User, error)
	GetAll() ([]*teal.User, error)
	Create(ctx context.Context, b *teal.User) (*teal.User, error)
	Update(ctx context.Context, id int, b *teal.User) (*teal.User, error)
	Delete(ctx context.Context, id int) error
}
