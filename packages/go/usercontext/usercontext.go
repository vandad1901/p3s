package usercontext

import (
	"context"
	"errors"
)

type UserCtxKey struct{}

var ErrNotFound = errors.New("user not found in context")

func CtxWithUser(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, UserCtxKey{}, userID)
}

func CtxUser(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(UserCtxKey{}).(int64)
	if !ok {
		return 0, ErrNotFound
	}

	return userID, nil
}
