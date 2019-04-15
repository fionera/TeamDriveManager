package create

import (
	"github.com/codegangsta/cli"
	. "github.com/fionera/TeamDriveManager/cmd"
	deleteProject "github.com/fionera/TeamDriveManager/cmd/delete/project"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "delete",
			Usage: "All commands that delete something, like selected Projects",
			Subcommands: []cli.Command{
				deleteProject.NewCommand(),
			},
		})
}
