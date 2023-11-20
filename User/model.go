package user

import (
	// "fmt"

	"os"
	"time"

	"github.com/chukwuka-emi/easysync/Services/auth"
	workspace "github.com/chukwuka-emi/easysync/Workspace"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// User ...
type User struct {
	gorm.Model
	Email            string                `gorm:"index;not null;" json:"email"`
	IsEmailConfirmed bool                  `gorm:"default:false;" json:"isEmailConfirmed"`
	Password         string                `json:"password"`
	FirstName        string                `json:"firstName"`
	LastName         string                `json:"lastName"`
	RealName         string                `json:"realName"`
	DisplayName      string                `json:"displayName"`
	Occupation       string                `json:"occupation"`
	ProfileImageURL  string                `json:"profileImageUrl"`
	OnboardingStep   string                `json:"onboardingStep"`
	ProfileState     string                `gorm:"type:enum('profile_active', 'profile_invited', 'profile_deactivated');" json:"profileState"`
	Workspaces       []workspace.Workspace `gorm:"many2many:user_workspaces;"`
	Roles            []Role                `gorm:"many2many:user_roles;"`
}

// RoleName ...
type RoleName string

const (
	// OWNER ...
	OWNER string = "owner"
	// ADMIN ...
	ADMIN string = "admin"
	// MEMBER ...
	MEMBER string = "member"
	// GUEST ...
	GUEST string = "guest"
)

// Role ...
type Role struct {
	gorm.Model
	Name string `gorm:"not null;" json:"name"`
}

const (
	// SetWorkspaceName ...
	SetWorkspaceName string = "set_workspace_name"
	// UpdateUserRealName ...
	UpdateUserRealName string = "update_user_real_name"
	// CreateChannel ...
	CreateChannel string = "create_channel"
	// OnboardingComplete ...
	OnboardingComplete string = "onboarding_complete"
)

const (
	// ProfileActive ...
	ProfileActive string = "profile_active"
	// ProfileInvited ...
	ProfileInvited string = "profile_invited"
	// ProfileDeactivated ...
	ProfileDeactivated string = "profile_deactivated"
)

// BuildAuthClaims ...
func BuildAuthClaims(userObj *User) auth.Claims {
	roleClaims := make([]interface{}, len(userObj.Roles))
	for i, role := range userObj.Roles {
		roleClaims[i] = role
	}

	authClaims := auth.Claims{
		User: auth.UserClaims{
			ID:    int(userObj.ID),
			Email: userObj.Email,
			Roles: roleClaims,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	return authClaims
}
