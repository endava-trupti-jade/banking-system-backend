package requestctx

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type contextKey string

const (
	loggerKey    contextKey = "logger"
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
	roleKey      contextKey = "role"
)

// =========================
// LOGGER
// =========================

func WithLogger(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func GetLogger(ctx context.Context) *zap.Logger {
	if log, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return log
	}

	// If logger missing from context: app still logs, no panic, no nil pointer crash
	return zap.L() // safe fallback
}

// =========================
// REQUEST ID
// =========================

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// =========================
// USER ID
// =========================

func WithUserID(ctx context.Context, userID primitive.ObjectID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (primitive.ObjectID, bool) {
	id, ok := ctx.Value(userIDKey).(primitive.ObjectID)
	return id, ok
}

// =========================
// Role
// =========================

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

func GetRole(ctx context.Context) (string, bool) {
	if role, ok := ctx.Value(roleKey).(string); ok {
		return role, true
	}
	return "", false
}

// func GetAuthCtx(ctx context.Context) (*AuthContext, error) {
// 	userID, ok1 := GetUserID(ctx)
// 	if !ok1 || userID == primitive.NilObjectID {
// 		return nil, constants.ErrInvalidUserIDCtx
// 	}

// 	role, ok2 := GetRole(ctx)
// 	if !ok2 || role == "" {
// 		return nil, constants.ErrInvalidRoleCtx
// 	}
// 	return &AuthContext{
// 		UserID: userID,
// 		Role:   role,
// 	}, nil
// }
