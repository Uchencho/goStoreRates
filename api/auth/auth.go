package auth

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func GenerateToken(username string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(username), 5)
	return base64.StdEncoding.EncodeToString(bytes), err
}
