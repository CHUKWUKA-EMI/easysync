package addmember

import (
	"net/http"

	data "github.com/chukwuka-emi/easysync/Data"
	channel "github.com/chukwuka-emi/easysync/Features/Channel"
	user "github.com/chukwuka-emi/easysync/Features/User"
	"github.com/chukwuka-emi/easysync/middlewares"
	"github.com/gin-gonic/gin"
)

// Handler links a user to a channel
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

	var channelData channel.Channel
	err = data.DB.Where("id=?", input.ChannelID).First(&channelData).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if channelData.OwnerEmail != currentUserClaims.Email {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add a member to this channel. Please contact the channel creator."})
		return
	}

	var userData user.User
	err = data.DB.Where("id=?", input.UserID).First(&userData).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userData.OnboardingStep != user.OnboardingComplete {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You cannot add a user that has not completed the onboarding steps."})
		return
	}

	queryResult := data.DB.Raw(`SELECT * FROM user_channels WHERE
	channel_id=? AND user_id=? LIMIT 1`, input.ChannelID, input.UserID)

	if queryResult.Error != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()})
		return
	}
	if queryResult.RowsAffected > 0 {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "User is already a member of this channel."})
		return
	}

	err = data.DB.Exec("INSERT INTO user_channels (user_id,channel_id) VALUES (?,?)", input.UserID, input.ChannelID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success", "data": userData})
}
