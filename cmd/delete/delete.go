package create

import (
	. "github.com/fionera/TeamDriveManager/cmd"
	deleteProject "github.com/fionera/TeamDriveManager/cmd/delete/project"
	deleteServiceaccount "github.com/fionera/TeamDriveManager/cmd/delete/serviceaccount"
	"github.com/urfave/cli"
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
