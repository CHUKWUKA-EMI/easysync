package auth

import (
	"errors"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

// UserClaims ...
type UserClaims struct {
	ID    int           `json:"id" binding:"required"`
	Email string        `json:"email" binding:"required"`
	Roles []interface{} `json:"roles" binding:"required"`
}

// Claims ...
type Claims struct {
	User UserClaims `json:"user" binding:"required"`
	jwt.RegisteredClaims
}

// SignJWT ...
func SignJWT(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyJWT ...
func VerifyJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("unknown claims type")
	}

	return claims, nil
}
