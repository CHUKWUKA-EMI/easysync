package createchannel

import (
	"fmt"
	"net/http"
	"strings"

	channel "github.com/chukwuka-emi/easysync/Channel"
	data "github.com/chukwuka-emi/easysync/Data"
	user "github.com/chukwuka-emi/easysync/User"
	"github.com/chukwuka-emi/easysync/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler creates a workspace channel
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

	if input.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "channel name cannot be empty."})
		return
	}

	if len(strings.Fields(input.Name)) > 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "channel name cannot have whitespace."})
		return
	}

	if input.WorkspaceID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "workspaceId is missing in request body."})
		return
	}

	if input.Type == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "channel type was not spacified."})
		return
	}

	if input.Type != channel.Private && input.Type != channel.Public {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "channel type is not valid."})
		return
	}

	var userData user.User
	queryResult := data.DB.Where("id=?", currentUserClaims.ID).First(&userData)
	if queryResult.RowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "current user record not found."})
		return
	}
	if queryResult.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying user data"})
		return
	}

	queryResult = data.DB.Where("name=? AND workspace_id=?", input.Name, input.WorkspaceID).First(&channel.Channel{})

	if queryResult.RowsAffected > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Channel with name: %s already exists.", input.Name)})
		return
	}

	channelData := channel.Channel{
		Name:        strings.ToLower(input.Name),
		OwnerEmail:  currentUserClaims.Email,
		Type:        input.Type,
		WorkspaceID: input.WorkspaceID,
	}
	err = data.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&channelData).Error; err != nil {
			return err
		}

		if userData.OnboardingStep == user.CreateChannel {
			if err := tx.Model(&userData).Updates(&user.User{OnboardingStep: user.OnboardingComplete}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "The system encountered an error while creating a channel. PLease try again."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": channelData})
}
