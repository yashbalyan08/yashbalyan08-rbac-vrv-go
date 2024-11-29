package utils

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("superkey"))

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token with username and role
func GenerateJWT(username, role string) (string, error) {
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ParseJWT parses and validates a JWT token
func ParseJWT(tokenString string) (*Claims, error) {
	log.Println("JWT: ", string(jwtKey))
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
