package context

import (
	"context"
	"fmt"

	"github.com/kencx/teal"
)

type baseKey string

const (
	bookKey     = baseKey("book")
	authorIdKey = baseKey("author")
	userIdKey   = baseKey("userId")
	userKey     = baseKey("user")
)

func WithBook(ctx context.Context, value *teal.Book) context.Context {
	return context.WithValue(ctx, bookKey, value)
}

func GetBook(ctx context.Context) (*teal.Book, error) {
	value, ok := ctx.Value(bookKey).(*teal.Book)
	if !ok {
		return nil, fmt.Errorf("ctx: failed to get Book from context")
	}
	return value, nil
}

func WithAuthorID(ctx context.Context, value int64) context.Context {
	return context.WithValue(ctx, authorIdKey, value)
}

func GetAuthorID(ctx context.Context) (int64, error) {
	value, ok := ctx.Value(authorIdKey).(int64)
	if !ok {
		return -1, fmt.Errorf("ctx: failed to get AuthorID from context")
	}
	return value, nil
}

func WithUserID(ctx context.Context, value int64) context.Context {
	return context.WithValue(ctx, userIdKey, value)
}

func GetUserID(ctx context.Context) (int64, error) {
	value, ok := ctx.Value(userIdKey).(int64)
	if !ok {
		return -1, fmt.Errorf("ctx: failed to get User from context")
	}
	return value, nil
}

func WithUser(ctx context.Context, value *teal.User) context.Context {
	return context.WithValue(ctx, userKey, value)
}

func GetUser(ctx context.Context) (*teal.User, error) {
	value, ok := ctx.Value(userKey).(*teal.User)
	if !ok {
		return nil, fmt.Errorf("ctx: failed to get User from context")
	}
	return value, nil
}
