package assign

import (
	"github.com/urfave/cli"

	. "github.com/fionera/TeamDriveManager/cmd"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "assign",
			Usage: "All commands that either assign stuff, like users to TeamDrives ",
			Subcommands: []cli.Command{
				NewAssignGroupCmd(),
				NewAssignAddressCmd(),
				NewAssignServiceAccountsCmd(),
			},
		})
}
