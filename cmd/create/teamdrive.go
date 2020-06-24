package create

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewCreateTeamDriveCommand() cli.Command {
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

	tokenSource, err := api.NewTokenSource(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate)
	if err != nil {
		logrus.Error(err)
		return
	}

	driveApi, err := api.NewDriveService(tokenSource)
	if err != nil {
		logrus.Error(err)
		return
	}

	teamDrive, err := api.CreateTeamDrive(driveApi, c.Args().First())
	if err != nil {
		logrus.Panic(err)
		return
	}

	logrus.Infof("Successfully created TeamDrive %s", teamDrive.Name)
}
