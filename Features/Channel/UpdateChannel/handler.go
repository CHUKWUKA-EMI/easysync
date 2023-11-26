package updatechannel

import (
	"fmt"
	"net/http"
	"strings"

	data "github.com/chukwuka-emi/easysync/Data"
	channel "github.com/chukwuka-emi/easysync/Features/Channel"
	"github.com/chukwuka-emi/easysync/middlewares"
	"github.com/gin-gonic/gin"
)

// Handler updates a channel
func Handler(ctx *gin.Context) {
	_, err := middlewares.GetUserClaimsFromRequestContext(ctx)

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

	if input.Name != "" && len(strings.Fields(input.Name)) > 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "channel name cannot have whitespace."})
		return
	}

	if input.Type != "" && input.Type != channel.Private && input.Type != channel.Public {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "channel type is not valid."})
		return
	}

	var channelData channel.Channel
	channelID := ctx.Param("id")

	queryResult := data.DB.Where("id=?", channelID).First(&channelData)
	if queryResult.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Channel not found."})
		return
	}
	if queryResult.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving channel information."})
		return
	}

	queryResult = data.DB.Not("id=?", channelID).Where("name=?", strings.ToLower(input.Name)).First(&channel.Channel{})
	fmt.Println("ROWS", queryResult.RowsAffected)
	if queryResult.RowsAffected > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "The provided channel name already belongs to another channel."})
		return
	}

	input.Name = strings.ToLower(input.Name)
	update := data.DB.Model(&channelData).Updates(&input)

	if update.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "The system encountered an error while updating channel information. PLease try again."})
		return
	}

	fmt.Println("UPDATED", channelData)

	ctx.JSON(http.StatusOK, gin.H{"data": channelData})
}
