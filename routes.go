package main

import (
	"context"
	"net/http"

	addmember "github.com/chukwuka-emi/easysync/Features/Channel/AddMember"
	createchannel "github.com/chukwuka-emi/easysync/Features/Channel/CreateChannel"
	updatechannel "github.com/chukwuka-emi/easysync/Features/Channel/UpdateChannel"
	chat "github.com/chukwuka-emi/easysync/Features/Chat"
	createconversation "github.com/chukwuka-emi/easysync/Features/Chat/CreateConversation"
	getchats "github.com/chukwuka-emi/easysync/Features/Chat/GetChats"
	accountlogin "github.com/chukwuka-emi/easysync/Features/User/AccountLogin"
	createuser "github.com/chukwuka-emi/easysync/Features/User/CreateUser"
	initiateemailverification "github.com/chukwuka-emi/easysync/Features/User/InitiateEmailVerification"
	updateuser "github.com/chukwuka-emi/easysync/Features/User/UpdateUser"
	acceptinvite "github.com/chukwuka-emi/easysync/Features/Workspace/AcceptInvite"
	invitecollaborator "github.com/chukwuka-emi/easysync/Features/Workspace/InviteCollaborator"
	updateworkspace "github.com/chukwuka-emi/easysync/Features/Workspace/UpdateWorkspace"
	"github.com/chukwuka-emi/easysync/middlewares"
	"github.com/gin-gonic/gin"
)

var ctx context.Context

func handleRoutes(r *gin.Engine) {

	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Server is Okay and running!")
	})

	r.POST("/email-check", initiateemailverification.Handler)
	r.POST("/user", createuser.Handler)
	r.POST("user/signin", accountlogin.Handler)
	r.GET("/workspace/invite/:token", acceptinvite.Handler)
	r.GET("/websocket/init", chat.ChatHub.HandleUpgrade)

	r.Use(middlewares.Authenticate)

	r.PATCH("/workspace", updateworkspace.Handler)
	r.POST("/workspace/invite", invitecollaborator.Handler)
	r.PATCH("/user", updateuser.Handler)
	r.POST("/channel", createchannel.Handler)
	r.PATCH("/channel/:id", updatechannel.Handler)
	r.POST("/channel/invite", addmember.Handler)
	r.POST("/conversations", createconversation.Handler)
	r.GET("/conversations/:id/chats", getchats.Handler)
}
