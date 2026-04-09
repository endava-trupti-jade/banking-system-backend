package utils

import (
	_ "log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(hash, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte((p))) == nil
}
