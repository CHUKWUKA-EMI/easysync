package main

import (
	"context"
	"net/http"

	createchannel "github.com/chukwuka-emi/easysync/Channel/CreateChannel"
	updatechannel "github.com/chukwuka-emi/easysync/Channel/UpdateChannel"
	accountlogin "github.com/chukwuka-emi/easysync/User/AccountLogin"
	createuser "github.com/chukwuka-emi/easysync/User/CreateUser"
	initiateemailverification "github.com/chukwuka-emi/easysync/User/InitiateEmailVerification"
	updateuser "github.com/chukwuka-emi/easysync/User/UpdateUser"
	acceptinvite "github.com/chukwuka-emi/easysync/Workspace/AcceptInvite"
	invitecollaborator "github.com/chukwuka-emi/easysync/Workspace/InviteCollaborator"
	updateworkspace "github.com/chukwuka-emi/easysync/Workspace/UpdateWorkspace"
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

	r.Use(middlewares.Authenticate)

	r.PATCH("/workspace", updateworkspace.Handler)
	r.POST("/workspace/invite", invitecollaborator.Handler)
	r.PATCH("/user", updateuser.Handler)
	r.POST("/channel", createchannel.Handler)
	r.PATCH("/channel/:id", updatechannel.Handler)
}
