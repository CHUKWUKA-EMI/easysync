package getchats

import (
	"net/http"
	"time"

	data "github.com/chukwuka-emi/easysync/Data"
	chatModel "github.com/chukwuka-emi/easysync/Features/Chat"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type responsePayload struct {
	Chats             []chatModel.Chat `json:"chats"`
	LastEvaluationKey uint64           `json:"lastEvaluationKey"`
	Size              uint             `json:"size"`
	Note              string           `json:"note"`
}

// Handler fetches chat messages for a specific channel
func Handler(ctx *gin.Context) {
	conversationIDStr := ctx.Param("id")

	if conversationIDStr == "" {
		ctx.String(http.StatusBadRequest, "missing channel id")
		return
	}

	conversationID, err := uuid.Parse(conversationIDStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input request
	err = ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := 15
	if input.Limit != 0 {
		limit = int(input.Limit)
	}

	lastEvaluationKey := uint64(time.Now().UnixMilli())
	if input.ExclusiveStartKey != 0 {
		lastEvaluationKey = input.ExclusiveStartKey
	}

	iter := data.AstraDBSession.Query(`SELECT
	                         id, 
							 conversation_id AS conversationId,
							 sender_id AS senderId,
							 message AS content,
							 created_at AS createdAt,
							 updated_at AS updatedAt
	                         FROM easysynk.chats 
	                         WHERE conversation_id=? AND created_at < ? LIMIT ?;`,
		gocql.UUID(conversationID),
		time.UnixMilli(int64(lastEvaluationKey)),
		limit,
	).Iter()

	var chats []chatModel.Chat
	var chat chatModel.Chat

	chats = []chatModel.Chat{}
	for iter.Scan(&chat.ID,
		&chat.ConversationID,
		&chat.SenderID,
		&chat.Content,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	) {
		chats = append(chats, chat)
	}

	err = iter.Close()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var lastItem chatModel.Chat
	var note string
	if len(chats) > 0 {
		lastItem = chats[len(chats)-1]
		note = "Pass the value of the 'lastEvaluationKey' to the 'exclusiveStartKey' to fetch the next page."
	} else {
		lastItem = chatModel.Chat{}
		note = ""
	}

	responseData := responsePayload{
		Chats:             chats,
		LastEvaluationKey: lastItem.ID,
		Size:              uint(len(chats)),
		Note:              note,
	}
	ctx.JSON(http.StatusOK, responseData)
}
