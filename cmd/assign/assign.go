package group

import (
	"github.com/codegangsta/cli"
	. "github.com/fionera/TeamDriveManager/cmd"
	assignGroup "github.com/fionera/TeamDriveManager/cmd/assign/group"
	assignServiceAccount "github.com/fionera/TeamDriveManager/cmd/assign/serviceaccount"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "assign",
			Usage: "All commands that either assign stuff, like users to TeamDrives ",
			Subcommands: []cli.Command{
				assignGroup.NewAssignGroupCmd(),
				assignServiceAccount.NewAssignServiceAccountsCmd(),
			},
		})
}
