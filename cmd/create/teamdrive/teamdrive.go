package teamdrive

import (
	"github.com/codegangsta/cli"
	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/drive"
	. "github.com/fionera/TeamDriveManager/config"
	"github.com/sirupsen/logrus"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "teamdrive",
		Usage:  "Create a Teamdrive",
		Action: CmdCreateTeamDrive,
		Flags:  []cli.Flag{},
	}
}

func CmdCreateTeamDrive(c *cli.Context) {
	if !c.Args().Present() {
		logrus.Error("Please supply a teamdrive name")
		return
	}

	client, err := api.CreateClient(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate, []string{drive.DriveScope})
	if err != nil {
		logrus.Error(err)
		return
	}

	driveApi, err := drive.NewApi(client)
	if err != nil {
		logrus.Error(err)
		return
	}

	teamDrive, err := driveApi.CreateTeamDrive(c.Args().First())
	if err != nil {
		logrus.Panic(err)
		return
	}

	logrus.Infof("Successfully created TeamDrive %s", teamDrive.Name)
}
