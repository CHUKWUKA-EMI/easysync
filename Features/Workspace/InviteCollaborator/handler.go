package invitecollaborator

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	data "github.com/chukwuka-emi/easysync/Data"
	user "github.com/chukwuka-emi/easysync/Features/User"
	workspace "github.com/chukwuka-emi/easysync/Features/Workspace"
	"github.com/chukwuka-emi/easysync/Services/auth"
	emailservice "github.com/chukwuka-emi/easysync/Services/email"
	"github.com/chukwuka-emi/easysync/middlewares"
	"github.com/chukwuka-emi/easysync/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userWorkspace struct {
	WorkspaceID uuid.UUID `json:"workspace_id"`
	UserID      uuid.UUID `json:"user_id"`
}

// Handler sends an invite to a workspace collaborator
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

	var workspaceData workspace.Workspace
	queryResult := data.DB.Where("id=?", input.WorkspaceID).First(&workspaceData)

	if queryResult.Error != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": queryResult.Error.Error()})
		return
	}

	// var userData user.User
	var result userWorkspace
	queryResult = data.DB.Raw(`SELECT * FROM user_workspaces WHERE
	workspace_id=? AND user_id=(SELECT id FROM users WHERE email=? LIMIT 1) LIMIT 1`,
		input.WorkspaceID, input.CollaboratorEmail).Scan(&result)

	if queryResult.Error != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()})
		return
	}
	if queryResult.RowsAffected > 0 {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("User with email %s is already on this workspace.", input.CollaboratorEmail)})
		return
	}

	collaborator := user.User{
		Email:            input.CollaboratorEmail,
		IsEmailConfirmed: false,
		OnboardingStep:   user.UpdateUserRealName,
		ProfileState:     user.ProfileInvited,
		Roles:            []user.Role{{Name: user.MEMBER}},
	}

	err = data.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&collaborator).Error
		if err != nil {
			return err
		}

		err = tx.Exec("INSERT INTO user_workspaces (user_id,workspace_id) VALUES (?,?)", collaborator.ID, workspaceData.ID).Error
		if err != nil {
			return err
		}

		invitationToken, err := auth.GenerateWorkspaceInvitationToken(workspaceData.ID, collaborator.ID)
		if err != nil {
			return err
		}
		invite := workspace.Invite{
			WorkspaceID: workspaceData.ID,
			UserID:      collaborator.ID,
			Status:      workspace.InvitePending,
			Token:       invitationToken,
		}
		err = tx.Create(&invite).Error

		if err != nil {
			return err
		}

		invitationURL := fmt.Sprintf("%s/workspace/invite/%s", os.Getenv("BACKEND_BASE_URL"), invitationToken)
		emailTemplate := &emailservice.WorkspaceInviteTemplate{
			WorkspaceName: workspaceData.Name,
			EmailSubject:  fmt.Sprintf("%s has invited you to collaborate on EasySync", workspaceData.Name),
			InvitationURL: invitationURL,
		}
		emailHTMLContent := emailTemplate.Build()
		email := &emailservice.EmailService{
			Sender:    emailservice.EmailUser{Name: "EasySync", Email: os.Getenv("SENDER_EMAIL")},
			Recipient: emailservice.EmailUser{Name: "", Email: input.CollaboratorEmail},
			Subject:   fmt.Sprintf("%s has invited you to collaborate on EasySync", workspaceData.Name),
			Content:   emailHTMLContent,
		}
		response, err := email.SendEmail()
		if err != nil {
			return err
		}

		if response.StatusCode != 200 && response.StatusCode != 201 {
			return errors.New(response.Status)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Our system encountered an error while processing your request. Please retry.", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Our system encountered an error while processing your request. Please retry."})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": collaborator})
}
