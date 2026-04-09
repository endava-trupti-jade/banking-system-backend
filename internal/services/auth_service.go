package services

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	"banking-system-backend/internal/repositories/mongorepo"
	"banking-system-backend/pkg/utils"
	"context"
	"log"
	"strings"
)

type AuthService struct {
	userRepo *mongorepo.UserRepository
}

func NewAuthService(repo *mongorepo.UserRepository) *AuthService {
	return &AuthService{userRepo: repo}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
	log.Println("AuthService Register() started")
	rawPassword := strings.TrimSpace(req.Password)

	hashedPassword, err := utils.HashPassword(rawPassword)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     constants.RoleCustomer,
	}

	user.Email = utils.NormalizeEmail(user.Email)

	log.Println("AuthService Register() end")
	return s.userRepo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	log.Println("AuthService Login() started")
	var res dto.LoginResponse

	email := utils.NormalizeEmail(req.Email)
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Println("invalid email or password")
		return res, constants.ErrInvalidEmailORPass
	}
	if user == nil {
		return res, constants.ErrInvalidEmailORPass
	}

	if !utils.ComparePassword(user.Password, req.Password) {
		return res, constants.ErrInvalidEmailORPass
	}

	token, err := utils.GenerateToken(user.ID.Hex(), user.Role)
	if err != nil {
		log.Println("invalid email or password")
		return res, constants.ErrUserIDRequired
	}

	res.Email = user.Email
	res.Role = user.Role
	res.Token = token

	log.Println("AuthService Login() end")
	return res, nil
}
