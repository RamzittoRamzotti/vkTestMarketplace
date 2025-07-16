package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"unicode/utf8"
)

func ValidateUsername(username string) error {
	if utf8.RuneCountInString(username) < 3 || utf8.RuneCountInString(username) > 32 {
		return errors.New("username must be 3-32 characters")
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	if !matched {
		return errors.New("username must contain only letters, numbers, and underscores")
	}
	return nil
}

func ValidatePassword(password string) error {
	if utf8.RuneCountInString(password) < 6 || utf8.RuneCountInString(password) > 64 {
		return errors.New("password must be 6-64 characters")
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
