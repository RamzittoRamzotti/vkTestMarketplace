package utils

import (
	"context"
)

type ctxKeyUserID struct{}

func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, ctxKeyUserID{}, userID)
}

func UserIDFromContext(ctx context.Context) int {
	v := ctx.Value(ctxKeyUserID{})
	if id, ok := v.(int); ok {
		return id
	}
	return 0
}
