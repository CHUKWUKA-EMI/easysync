package createconversation

import (
	"net/http"
	"time"

	data "github.com/chukwuka-emi/easysync/Data"
	chat "github.com/chukwuka-emi/easysync/Features/Chat"
	"github.com/chukwuka-emi/easysync/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

// Handler creates a conversation
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

	if input.TargetUserID == currentUserClaims.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "You cannot create a conversation between yourself and yourself"})
		return
	}

	iter := data.AstraDBSession.Query(`SELECT members FROM easysynk.conversations 
	WHERE members CONTAINS ? AND members CONTAINS ? LIMIT 1;`,
		gocql.UUID(currentUserClaims.ID), gocql.UUID(input.TargetUserID)).Iter()

	if iter.NumRows() > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Conversation already exists."})
		return
	}

	err = iter.Close()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	conversation := chat.Conversation{
		ID:        gocql.TimeUUID(),
		Members:   []gocql.UUID{gocql.UUID(currentUserClaims.ID), gocql.UUID(input.TargetUserID)},
		IsOpen:    true,
		CreatedAt: time.Now(),
	}

	err = data.AstraDBSession.Query(`INSERT into easysynk.conversations 
	                         (id,members,is_open,created_at) 
							 VALUES (?,?,?,?);`,
		conversation.ID,
		conversation.Members,
		conversation.IsOpen,
		conversation.CreatedAt).Exec()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": conversation})
}
