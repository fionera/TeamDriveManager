package group

import (
	. "github.com/fionera/TeamDriveManager/cmd"
	assignAddress "github.com/fionera/TeamDriveManager/cmd/assign/address"
	assignGroup "github.com/fionera/TeamDriveManager/cmd/assign/group"
	assignServiceAccount "github.com/fionera/TeamDriveManager/cmd/assign/serviceaccount"
	"github.com/urfave/cli"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "assign",
			Usage: "All commands that either assign stuff, like users to TeamDrives ",
			Subcommands: []cli.Command{
				assignGroup.NewAssignGroupCmd(),
				assignAddress.NewAssignAddressCmd(),
				assignServiceAccount.NewAssignServiceAccountsCmd(),
			},
		})
}
