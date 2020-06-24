package list

import (
	"github.com/urfave/cli"

	. "github.com/fionera/TeamDriveManager/cmd"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "list",
			Usage: "All commands that list something, like all TeamDrives",
			Subcommands: []cli.Command{
				NewListTeamDriveCommand(),
				NewListGroupCommand(),
				NewListProjectCommand(),
				NewListServiceAccountCommand(),
			},
		})
}
