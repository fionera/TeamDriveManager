package delete

import (
	"github.com/urfave/cli"

	. "github.com/fionera/TeamDriveManager/cmd"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "delete",
			Usage: "All commands that delete something, like selected Projects",
			Subcommands: []cli.Command{
				NewDeleteAddressCommand(),
				NewDeleteProjectCommand(),
				NewDeleteServiceAccountCommand(),
				NewDeleteTeamDriveCommand(),
			},
		})
}
