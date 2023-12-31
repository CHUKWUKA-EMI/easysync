package createuser

import (
	"context"
	"net/http"
	"os"

	data "github.com/chukwuka-emi/easysync/Data"
	channel "github.com/chukwuka-emi/easysync/Features/Channel"
	user "github.com/chukwuka-emi/easysync/Features/User"
	services "github.com/chukwuka-emi/easysync/Services"
	"github.com/chukwuka-emi/easysync/Services/auth"

	workspace "github.com/chukwuka-emi/easysync/Features/Workspace"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var ctx = context.Background()

// Handler ...
func Handler(httpContext *gin.Context) {
	var input createUserRequest

	err := httpContext.ShouldBindJSON(&input)

	if err != nil {
		httpContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = services.RedisClient.Get(ctx, input.Email).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			httpContext.JSON(http.StatusBadRequest, gin.H{"error": "Confirmation code has expired"})
			return
		}

		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userObj := user.User{
		Email:            input.Email,
		IsEmailConfirmed: true,
		OnboardingStep:   user.SetWorkspaceName,
		ProfileState:     user.ProfileActive,
		Roles: []user.Role{
			{Name: user.OWNER},
			{Name: user.ADMIN},
		},
		Workspaces: []workspace.Workspace{
			{
				Channels: []channel.Channel{
					{
						Name:        "company-wide",
						Description: "General channel",
						Type:        channel.Public,
						OwnerEmail:  input.Email,
					},
				},
			},
		},
	}

	err = data.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&userObj).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	authClaims := user.BuildAuthClaims(&userObj)

	authTokens, err := auth.SignJWT(authClaims)
	if err != nil {
		httpContext.JSON(http.StatusInternalServerError, gin.H{"error": "The system encountered an error. Please try again"})
		return
	}
	// save refresh token
	go data.DB.Create(&user.Token{RefreshToken: authTokens.RefreshToken, UserID: userObj.ID})

	httpContext.SetCookie("accessToken", authTokens.AccessToken, 3600, "/", os.Getenv("DOMAIN"), false, true)
	httpContext.SetCookie("refreshToken", authTokens.RefreshToken, 365*24*3600, "/", os.Getenv("DOMAIN"), false, true)

	httpContext.JSON(http.StatusOK, gin.H{"data": &userObj})
}
