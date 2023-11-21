package acceptinvite

import (
	"net/http"
	"os"

	data "github.com/chukwuka-emi/easysync/Data"
	"github.com/chukwuka-emi/easysync/Services/auth"
	user "github.com/chukwuka-emi/easysync/User"
	workspace "github.com/chukwuka-emi/easysync/Workspace"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Handler handles workspace invitation acceptance
func Handler(ctx *gin.Context) {
	tokenString := ctx.Param("token")

	if tokenString == "" {
		ctx.String(http.StatusBadRequest, "Token is not present in URL")
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &auth.WorkspaceInviteClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	claims, ok := token.Claims.(*auth.WorkspaceInviteClaims)

	if !ok {
		ctx.String(http.StatusBadRequest, "Malformed token")
		return
	}

	var inviteData workspace.Invite
	err = data.DB.Where("user_id=? AND workspace_id=?", claims.UserID, claims.WorkspaceID).First(&inviteData).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userData user.User
	err = data.DB.Where("id=?", claims.UserID).First(&userData).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var workspaceData workspace.Workspace
	err = data.DB.Where("id=?", claims.WorkspaceID).First(&workspaceData).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = data.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&userData).Updates(&user.User{
			IsEmailConfirmed: true,
			OnboardingStep:   user.UpdateUserRealName,
			ProfileState:     user.ProfileActive,
		}).Error

		if err != nil {
			return err
		}

		err = tx.Delete(&inviteData).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(301, os.Getenv("FRONTEND_URL"))
}
