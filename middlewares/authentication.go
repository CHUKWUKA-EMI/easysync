package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/chukwuka-emi/easysync/Services/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Authenticate ...
func Authenticate(ctx *gin.Context) {
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
		return
	}

	authHeaderEntriesArray := strings.Fields(authHeader)
	if len(authHeaderEntriesArray) != 2 {

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid/Malformed authorization header"})
		return
	}

	token := authHeaderEntriesArray[1]

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing auth token or Invalid auth scheme"})
		return
	}

	authClaims, err := auth.VerifyJWT(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if authClaims.ExpiresAt.Unix() < jwt.NewNumericDate(time.Now()).Unix() {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
		return
	}
	ctx.Set("user", authClaims.User)

	ctx.Next()
}

// GetUserClaimsFromRequestContext extracts user claims from http context
func GetUserClaimsFromRequestContext(ctx *gin.Context) (auth.UserClaims, error) {
	currentUser, ok := ctx.Get("user")
	if !ok {
		return auth.UserClaims{}, errors.New("user authentication data missing in request context")
	}

	currentUserClaims, ok := currentUser.(auth.UserClaims)
	if !ok {
		return auth.UserClaims{}, errors.New("user auth claims has invalid format")
	}

	return currentUserClaims, nil
}
