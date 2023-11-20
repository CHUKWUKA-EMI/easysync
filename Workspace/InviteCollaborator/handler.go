package invitecollaborator

// import (
// 	"fmt"
// 	"net/http"

// 	data "github.com/chukwuka-emi/easysync/Data"
// 	user "github.com/chukwuka-emi/easysync/User"
// 	"github.com/chukwuka-emi/easysync/middlewares"
// 	"github.com/chukwuka-emi/easysync/utils"
// 	"github.com/gin-gonic/gin"
// )

// type userWorkspace struct {
// 	WorkspaceID uint `json:"workspace_id"`
// 	UserID      uint `json:"user_id"`
// }

// // Handler sends an invite to a workspace collaborator
// func Handler(ctx *gin.Context) {
// 	currentUserClaims, err := middlewares.GetUserClaimsFromRequestContext(ctx)

// 	if err != nil {
// 		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
// 		return
// 	}

// 	hasRequiredPermissions := utils.HasRequiredPermissions(currentUserClaims.Roles, []string{"owner"})

// 	if !hasRequiredPermissions {
// 		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have the required permissions to perform this operation."})
// 		return
// 	}

// 	var input request
// 	err = ctx.ShouldBindJSON(&input)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// var workspaceData workspace.Workspace
// 	// var userData user.User
// 	var result userWorkspace
// 	queryResult := data.DB.Raw(`SELECT * FROM user_workspaces WHERE
// 	workspace_id=? AND user_id=(SELECT id FROM users WHERE email=?)`,
// 		input.WorkspaceID, input.CollaboratorEmail).Scan(&result)

// 	if queryResult.RowsAffected > 0 {
// 		ctx.JSON(http.StatusBadRequest,
// 			gin.H{"error": fmt.Sprintf("User with email %s is already in this workspace.", input.CollaboratorEmail)})
// 		return
// 	}

// 	collaborator := user.User{
// 		Email: input.CollaboratorEmail,
// 		IsEmailConfirmed: false,
// 		OnboardingStep: user.UpdateUserRealName,
// 		ProfileState: user.ProfileInvited,
// 		Roles: []user.Role{{Name: user.MEMBER}},
// 	}
// 	data.DB.Create()
// 	ctx.JSON(http.StatusOK, gin.H{"data": result})
// }
