package delete

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewDeleteAddressCommand() cli.Command {
	return cli.Command{
		Name:      "address",
		Usage:     "Delete teamdrive assignmend for given address",
		Action:    CmdDeleteAddress,
		Flags:     []cli.Flag{},
		UsageText: "<TEAMDRIVE-NAME> <ADDRESS-TO-REMOVE>",
	}
}

func CmdDeleteAddress(c *cli.Context) {

	teamDriveName := c.Args().Get(0)
	address := c.Args().Get(1)

	if teamDriveName == "" {
		logrus.Error("Please supply a teamdrive name")
		return
	}

	if address == "" {
		logrus.Error("Please supply an address")
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

	teamDrives, err := api.ListTeamDrives(driveApi)
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, teamDrive := range teamDrives {
		if teamDrive.Name == teamDriveName {

			permissions, err := api.ListPermissions(driveApi, teamDrive.Id)
			for _, permission := range permissions {
				if permission.EmailAddress == address {
					err := api.DeletePermission(driveApi, teamDrive.Id, permission.Id)
					if err != nil {
						logrus.Error(err)
						return
					}
				}
			}
			if err != nil {
				logrus.Error(err)
				return
			}

			logrus.Info("Deleted Permission!")

			break
		}
	}

}
