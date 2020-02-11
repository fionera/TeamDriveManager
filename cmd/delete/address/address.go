package teamdrive

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/admin"
	"github.com/fionera/TeamDriveManager/api/drive"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewCommand() cli.Command {
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

	client, err := api.CreateClient(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate, []string{admin.AdminDirectoryGroupScope, drive.DriveScope})
	if err != nil {
		logrus.Error(err)
		return
	}

	driveApi, err := drive.NewApi(client)
	if err != nil {
		logrus.Error(err)
		return
	}

	teamDrives, err := driveApi.ListTeamDrives()
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, teamDrive := range teamDrives {
		if teamDrive.Name == teamDriveName {

			permissions, err := driveApi.ListPermissions(teamDrive.Id)
			for _, permission := range permissions {
				if permission.EmailAddress == address {
					err := driveApi.DeletePermission(teamDrive.Id, permission.Id)
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
