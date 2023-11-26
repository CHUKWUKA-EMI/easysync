package user

import (
	// "fmt"

	"os"
	"time"

	data "github.com/chukwuka-emi/easysync/Data"
	channel "github.com/chukwuka-emi/easysync/Features/Channel"
	workspace "github.com/chukwuka-emi/easysync/Features/Workspace"
	"github.com/chukwuka-emi/easysync/Services/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// User ...
type User struct {
	data.BaseModel
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
	Channels         []channel.Channel     `gorm:"many2many:user_channels;"`
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
	data.BaseModel
	Name string `gorm:"not null;" json:"name"`
}

// Token ...
type Token struct {
	data.BaseModel
	RefreshToken string    `json:"refreshToken"`
	UserID       uuid.UUID `gorm:"type:char(36);" json:"userId"`
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
			ID:    userObj.ID,
			Email: userObj.Email,
			Roles: roleClaims,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	return authClaims
}
