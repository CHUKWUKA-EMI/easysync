package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	channel "github.com/chukwuka-emi/easysync/Channel"
	user "github.com/chukwuka-emi/easysync/User"
	workspace "github.com/chukwuka-emi/easysync/Workspace"
)

func main() {
	stmts, err := gormschema.New("mysql").Load(
		&user.User{},
		&user.Role{},
		&workspace.Workspace{},
		&channel.Channel{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
