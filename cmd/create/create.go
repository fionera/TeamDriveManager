package create

import (
	. "github.com/fionera/TeamDriveManager/cmd"
	"github.com/urfave/cli"

	createGroup "github.com/fionera/TeamDriveManager/cmd/create/group"
	createProject "github.com/fionera/TeamDriveManager/cmd/create/project"
	createServiceAccount "github.com/fionera/TeamDriveManager/cmd/create/serviceaccount"
	createTeamDrive "github.com/fionera/TeamDriveManager/cmd/create/teamdrive"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "create",
			Usage: "All commands that create something, like a TeamDrive",
			Subcommands: []cli.Command{
				createTeamDrive.NewCommand(),
				createProject.NewCommand(),
				createServiceAccount.NewCommand(),
				createGroup.NewCommand(),
			},
		})
}
