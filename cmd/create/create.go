package create

import (
	"github.com/codegangsta/cli"
	. "github.com/fionera/TeamdriveManager/cmd"

	createTeamdrive "github.com/fionera/TeamdriveManager/cmd/create/teamdrive"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "create",
			Usage: "All commands that create something, like a Teamdrive",
			Subcommands: []cli.Command{
				createTeamdrive.NewCommand(),
			},
		})
}
