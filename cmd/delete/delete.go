package create

import (
	"github.com/urfave/cli"

	. "github.com/fionera/TeamDriveManager/cmd"
	deleteProject "github.com/fionera/TeamDriveManager/cmd/delete/project"
	deleteServiceaccount "github.com/fionera/TeamDriveManager/cmd/delete/serviceaccount"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "delete",
			Usage: "All commands that delete something, like selected Projects",
			Subcommands: []cli.Command{
				deleteProject.NewCommand(),
				deleteServiceaccount.NewCommand(),
			},
		})
}
