package interfaces

import (
	"banking-system-backend/internal/dto"
	"context"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
}
