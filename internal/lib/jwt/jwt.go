package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

func GenerateToken(userID int, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(30 * time.Minute).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr, secret string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}
	uidFloat, ok := claims["uid"].(float64)
	if !ok {
		return 0, errors.New("uid not found")
	}
	return int(uidFloat), nil
}
