package create

import (
	"github.com/codegangsta/cli"
	. "github.com/fionera/TeamdriveManager/cmd"

	createProject "github.com/fionera/TeamdriveManager/cmd/create/project"
	createServiceAccount "github.com/fionera/TeamdriveManager/cmd/create/serviceaccount"
	createTeamDrive "github.com/fionera/TeamdriveManager/cmd/create/teamdrive"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "create",
			Usage: "All commands that create something, like a Teamdrive",
			Subcommands: []cli.Command{
				createTeamDrive.NewCommand(),
				createProject.NewCommand(),
				createServiceAccount.NewCommand(),
			},
		})
}
