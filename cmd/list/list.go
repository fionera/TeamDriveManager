package create

import (
	"github.com/codegangsta/cli"
	. "github.com/fionera/TeamdriveManager/cmd"
	listGroups "github.com/fionera/TeamdriveManager/cmd/list/group"
	listProjects "github.com/fionera/TeamdriveManager/cmd/list/project"
	listServiceAccounts "github.com/fionera/TeamdriveManager/cmd/list/serviceaccount"
	listTeamDrives "github.com/fionera/TeamdriveManager/cmd/list/teamdrive"
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
