package services

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	"banking-system-backend/internal/repositories/mongorepo"
	"banking-system-backend/internal/requestctx"
	"banking-system-backend/pkg/utils"
	"context"
	"strings"

	"go.uber.org/zap"
)

type AuthService struct {
	userRepo *mongorepo.UserRepository
}

func NewAuthService(repo *mongorepo.UserRepository) *AuthService {
	return &AuthService{userRepo: repo}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {

	email := utils.NormalizeEmail(req.Email)
	log := requestctx.GetLogger(ctx).With(
		zap.String("email", email),
	)
	log.Info("register started")

	rawPassword := strings.TrimSpace(req.Password)

	hashedPassword, err := utils.HashPassword(rawPassword)
	if err != nil {
		log.Error("failed to hash password",
			zap.Error(err),
		)
		return err
	}

	user := &models.User{
		Email:    email,
		Password: hashedPassword,
		Role:     constants.RoleCustomer,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		log.Error("failed to create user",
			zap.Error(err),
		)
		return err
	}

	log.Info("registration successful")

	return nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {

	email := utils.NormalizeEmail(req.Email)

	log := requestctx.GetLogger(ctx).With(
		zap.String("email", email),
	)

	log.Info("login started")

	var res dto.LoginResponse

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Warn("invalid email or password")
		return res, constants.ErrInvalidEmailORPass
	}
	if user == nil {
		log.Warn("invalid email or password")
		return res, constants.ErrInvalidEmailORPass
	}

	if !utils.ComparePassword(user.Password, req.Password) {
		log.Warn("invalid email or password")
		return res, constants.ErrInvalidEmailORPass
	}

	token, err := utils.GenerateToken(user.ID.Hex(), user.Role)
	if err != nil {
		log.Error("failed to generate token",
			zap.String("user_id", user.ID.Hex()),
			zap.Error(err),
		)
		return res, constants.ErrUserIDRequired
	}

	res.Email = user.Email
	res.Role = user.Role
	res.Token = token

	log.Info("login successful",
		zap.String("user_id", user.ID.Hex()),
		zap.String("role", user.Role),
	)
	return res, nil
}
