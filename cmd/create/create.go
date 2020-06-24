package create

import (
	"github.com/urfave/cli"

	. "github.com/fionera/TeamDriveManager/cmd"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "create",
			Usage: "All commands that create something, like a TeamDrive",
			Subcommands: []cli.Command{
				NewCreateTeamDriveCommand(),
				NewCreateProjectCommand(),
				NewCreateServiceAccountCommand(),
				NewCreateGroupCommand(),
			},
		})
}
