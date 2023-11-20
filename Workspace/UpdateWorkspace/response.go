package updateworkspace

import (
	workspace "github.com/chukwuka-emi/easysync/Workspace"
)

type response struct {
	OnboardingStep string              `json:"onboardingStep"`
	Workspace      workspace.Workspace `json:"workspace"`
}
