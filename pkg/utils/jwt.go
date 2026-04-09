package utils

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/config"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtKey = []byte(config.AppConfig.JWTSecret)
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, role string) (string, error) {
	if len(jwtKey) == 0 {
		return "", constants.ErrMissingJWTSecret
	}

	if userID == "" {
		return "", constants.ErrUserIDRequired
	}

	if role == "" {
		return "", constants.ErrRoleRequired
	}

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token == nil {
		return "", constants.ErrTokenCreationFailed
	}

	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("error in signed token")
		return "", constants.ErrTokenSigningFailed
	}
	return signedToken, nil
}

func ValidateToken(tokenStr string) (string, string, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", "", constants.ErrInvalidToken
	}

	return claims.UserID, claims.Role, nil
}
