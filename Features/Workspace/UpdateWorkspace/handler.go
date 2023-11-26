package updateworkspace

import (
	"net/http"

	data "github.com/chukwuka-emi/easysync/Data"
	userSlice "github.com/chukwuka-emi/easysync/Features/User"
	model "github.com/chukwuka-emi/easysync/Features/Workspace"
	"github.com/chukwuka-emi/easysync/middlewares"
	"github.com/chukwuka-emi/easysync/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler ...
func Handler(ctx *gin.Context) {
	currentUserClaims, err := middlewares.GetUserClaimsFromRequestContext(ctx)

	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	hasRequiredPermissions := utils.HasRequiredPermissions(currentUserClaims.Roles, []string{"owner"})

	if !hasRequiredPermissions {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have the required permissions to perform this operation."})
		return
	}
	var input request
	err = ctx.ShouldBindJSON(&input)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user userSlice.User
	queryResult := data.DB.Where("id=?", currentUserClaims.ID).First(&user)

	if queryResult.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User record not found."})
		return
	}

	if queryResult.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user information"})
		return
	}

	var workspaceData model.Workspace

	queryResult = data.DB.Where("id=?", input.ID).First(&workspaceData)
	if queryResult.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Workspace with the provided ID does not exist."})
		return
	}
	if queryResult.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving workspace data."})
		return
	}

	err = data.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&workspaceData).Update("name", input.Name).Error; err != nil {
			return err
		}

		if user.OnboardingStep == userSlice.SetWorkspaceName {
			if err := tx.Model(&user).Updates(&userSlice.User{OnboardingStep: userSlice.UpdateUserRealName}).Error; err != nil {
				return err
			}

		}
		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong."})
		return
	}

	ctx.JSON(http.StatusOK,
		gin.H{"data": &response{
			OnboardingStep: user.OnboardingStep,
			Workspace:      workspaceData}})
}
