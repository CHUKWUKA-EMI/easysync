package updateuser

import (
	"net/http"

	data "github.com/chukwuka-emi/easysync/Data"
	user "github.com/chukwuka-emi/easysync/User"
	"github.com/chukwuka-emi/easysync/middlewares"
	"github.com/gin-gonic/gin"
)

// Handler updates a user profile
func Handler(ctx *gin.Context) {
	currentUserClaims, err := middlewares.GetUserClaimsFromRequestContext(ctx)

	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var input request

	err = ctx.ShouldBindJSON(&input)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userData user.User
	result := data.DB.Where("id=?", currentUserClaims.ID).First(&userData)

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User record not found."})
		return
	}
	if userData.OnboardingStep == user.UpdateUserRealName {
		input.OnboardingStep = user.CreateChannel
	}
	err = data.DB.Model(&userData).Omit("email").Updates(&input).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "The system encountered an error while updating user profile"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": userData})
}
