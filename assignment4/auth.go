package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("ezeze1337")

// UserClaims FOR CONTAINING
type UserClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT GENERATES JWT 
func GenerateJWT(username, role string) (string, error) {
	claims := &UserClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)), // TOKEN WORKS EXACTLY 1 MIN
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateJWT VALIDATES JWT AND RETURN CLAIMS
func ValidateJWT(tokenStr string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("INVALID TOKEN")
	}
	return claims, nil
}
