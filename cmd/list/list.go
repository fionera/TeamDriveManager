package create

import (
	. "github.com/fionera/TeamDriveManager/cmd"
	listGroups "github.com/fionera/TeamDriveManager/cmd/list/group"
	listProjects "github.com/fionera/TeamDriveManager/cmd/list/project"
	listServiceAccounts "github.com/fionera/TeamDriveManager/cmd/list/serviceaccount"
	listTeamDrives "github.com/fionera/TeamDriveManager/cmd/list/teamdrive"
	"github.com/urfave/cli"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "list",
			Usage: "All commands that list something, like all TeamDrives",
			Subcommands: []cli.Command{
				listTeamDrives.NewCommand(),
				listGroups.NewCommand(),
				listProjects.NewCommand(),
				listServiceAccounts.NewCommand(),
			},
		})
}
