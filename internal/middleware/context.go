package middleware

import (
	"context"
	"errors"
)

func SetUserToContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextKey, userID)
}

func GetUserFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(contextKey).(string)
	if !ok {
		return "", errors.New("userID not found in context")
	}
	return userID, nil
}
