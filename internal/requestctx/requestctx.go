package requestctx

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type contextKey string

const (
	LoggerKey    contextKey = "logger"
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id"
	RoleKey      contextKey = "role"
)

// =========================
// LOGGER
// =========================

func WithLogger(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, log)
}

func GetLogger(ctx context.Context) *zap.Logger {
	if log, ok := ctx.Value(LoggerKey).(*zap.Logger); ok {
		return log
	}
	return zap.L() // safe fallback
}

// =========================
// REQUEST ID
// =========================

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// =========================
// USER ID
// =========================

func WithUserID(ctx context.Context, userID primitive.ObjectID) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func GetUserID(ctx context.Context) (primitive.ObjectID, bool) {
	id, ok := ctx.Value(UserIDKey).(primitive.ObjectID)
	return id, ok
}

// =========================
// Role
// =========================

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, RoleKey, role)
}

func GetRole(ctx context.Context) string {
	if role, ok := ctx.Value(RoleKey).(string); ok {
		return role
	}
	return ""
}
