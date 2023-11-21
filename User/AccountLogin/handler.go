package accountlogin

import (
	"context"
	"net/http"
	"os"
	"strings"

	data "github.com/chukwuka-emi/easysync/Data"
	services "github.com/chukwuka-emi/easysync/Services"
	"github.com/chukwuka-emi/easysync/Services/auth"
	user "github.com/chukwuka-emi/easysync/User"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

// Handler ...
func Handler(httpContext *gin.Context) {
	var input accountLoginRequest

	err := httpContext.ShouldBindJSON(&input)

	if err != nil {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = services.RedisClient.Get(ctx, input.Email).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			httpContext.JSON(http.StatusBadRequest, gin.H{"error": "login code has expired"})
			return
		}

		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var userData user.User
	err = data.DB.Preload("Roles").Where("email=?", input.Email).Omit("password").First(&userData).Error

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			httpContext.JSON(http.StatusNotFound, gin.H{"error": "User's record not found"})
			return
		}
		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	authClaims := user.BuildAuthClaims(&userData)

	authTokens, err := auth.SignJWT(authClaims)
	if err != nil {
		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": "The system encountered an error. Please try again"})
		return
	}

	var existingTokenData user.Token
	queryResult := data.DB.Where("user_id=?", userData.ID).First(&existingTokenData)

	if queryResult.Error != nil {
		if strings.Contains(err.Error(), "record not found") {

			data.DB.Create(&user.Token{RefreshToken: authTokens.RefreshToken, UserID: userData.ID})

		} else {
			httpContext.JSON(http.StatusInternalServerError,
				gin.H{"error": err.Error()})
			return
		}
	} else {
		err = data.DB.Model(&existingTokenData).Updates(&user.Token{RefreshToken: authTokens.RefreshToken}).Error
		if err != nil {
			httpContext.JSON(http.StatusInternalServerError, gin.H{"error": "The system encountered an error. Please try again"})
			return
		}
	}

	httpContext.SetCookie("accessToken", authTokens.AccessToken, 3600, "/", os.Getenv("DOMAIN"), false, true)
	httpContext.SetCookie("refreshToken", authTokens.RefreshToken, 365*24*3600, "/", os.Getenv("DOMAIN"), false, true)

	httpContext.JSON(http.StatusOK, gin.H{"message": "Authentication successful."})
}
