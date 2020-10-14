package delete

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/api/googleapi"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewDeleteTeamDriveCommand() cli.Command {
	return cli.Command{
		Name:   "teamdrive",
		Usage:  "Delete a Teamdrive",
		Action: CmdDeleteTeamDrive,
		Flags:  []cli.Flag{},
	}
}

func CmdDeleteTeamDrive(c *cli.Context) {
	if !c.Args().Present() {
		logrus.Error("Please supply a teamdrive id")
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

	err = api.DeleteTeamDrive(driveApi, c.Args().First())
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok {
			switch gerr.Code {
			case 403:
				logrus.Error("Teamdrive contains objects and therefore cannot be deleted.")
				return
			default:
				logrus.Error("An error occurred when deleting account.", err)
				return
			}
		} else {
			logrus.Fatal("An unknown error occurred: ", err)
		}
	}

	logrus.Infof("Successfully deleted TeamDrive %s", c.Args().First())
}
