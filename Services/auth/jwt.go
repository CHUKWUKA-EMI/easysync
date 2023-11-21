package auth

import (
	"errors"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UserClaims ...
type UserClaims struct {
	ID    uuid.UUID     `json:"id" binding:"required"`
	Email string        `json:"email" binding:"required"`
	Roles []interface{} `json:"roles" binding:"required"`
}

// Claims ...
type Claims struct {
	User UserClaims `json:"user" binding:"required"`
	jwt.RegisteredClaims
}

// WorkspaceInviteClaims ...
type WorkspaceInviteClaims struct {
	WorkspaceID uuid.UUID `json:"workspaceID"`
	UserID      uuid.UUID `json:"userID"`
	jwt.RegisteredClaims
}

// Tokens is an object containing access token and refresh token
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// SignJWT generates JWT with standard and custom claims
func SignJWT(claims Claims) (Tokens, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	accessTokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return Tokens{}, err
	}

	refreshTokenString, err := generateRefreshToken()
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func generateRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    os.Getenv("JWT_ISSUER"),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(365 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	})

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyJWT verifies and decodes JWT
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

// RefreshToken regenerates an access token for a user
func RefreshToken(refreshTokenString string, claims Claims) (Tokens, error) {
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !refreshToken.Valid {
		return Tokens{}, err
	}

	tokens, err := SignJWT(claims)

	if err != nil {
		return Tokens{}, err
	}

	return tokens, nil
}

// GenerateWorkspaceInvitationToken generates token for a workspace invite
func GenerateWorkspaceInvitationToken(workspaceID uuid.UUID, userID uuid.UUID) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &WorkspaceInviteClaims{
		WorkspaceID: workspaceID,
		UserID:      userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(365 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
