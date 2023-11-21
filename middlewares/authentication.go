package middlewares

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	data "github.com/chukwuka-emi/easysync/Data"
	"github.com/chukwuka-emi/easysync/Services/auth"
	user "github.com/chukwuka-emi/easysync/User"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func handleTokenRefresh(refreshToken string) (auth.Tokens, error) {
	var tokenData user.Token
	err := data.DB.Where("refresh_token=?", refreshToken).First(&tokenData).Error
	if err != nil {
		var errorMessage error
		if err.Error() == "record not found" {
			errorMessage = errors.New("authentication refresh token not found. Please login again")
		} else {
			errorMessage = err
		}
		return auth.Tokens{}, errorMessage
	}
	var userData user.User
	err = data.DB.Preload("Roles").Where("id=?", tokenData.UserID).First(&userData).Error
	if err != nil {
		return auth.Tokens{}, err
	}

	claims := user.BuildAuthClaims(&userData)
	tokens, err := auth.RefreshToken(refreshToken, claims)
	if err != nil {
		return auth.Tokens{}, err
	}
	return tokens, nil
}

// Authenticate ...
func Authenticate(ctx *gin.Context) {
	token, err := ctx.Cookie("accessToken")
	if err != nil {
		if strings.Contains(err.Error(), "named cookie not present") {
			refreshToken, err := ctx.Cookie("refreshToken")
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}

			authTokens, err := handleTokenRefresh(refreshToken)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			token = authTokens.AccessToken
			ctx.SetCookie("accessToken", authTokens.AccessToken, 3600, "/", os.Getenv("DOMAIN"), false, true)
			ctx.SetCookie("refreshToken", authTokens.RefreshToken, 365*24*3600, "/", os.Getenv("DOMAIN"), false, true)
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
	}

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing auth token"})
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
