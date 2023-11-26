package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	channel "github.com/chukwuka-emi/easysync/Features/Channel"
	user "github.com/chukwuka-emi/easysync/Features/User"
	workspace "github.com/chukwuka-emi/easysync/Features/Workspace"
)

func main() {
	stmts, err := gormschema.New("mysql").Load(
		&user.User{},
		&user.Role{},
		&workspace.Workspace{},
		&channel.Channel{},
		&workspace.Invite{},
		&user.Token{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
