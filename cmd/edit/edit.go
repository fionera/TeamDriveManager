package list

import (
	"github.com/urfave/cli"

	. "github.com/fionera/TeamDriveManager/cmd"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "edit",
			Usage: "All commands that edit something, like unhide TeamDrives",
			Subcommands: []cli.Command{
				NewUnhideTeamDriveCommand(),
				NewHideTeamDriveCommand(),
			},
		})
}