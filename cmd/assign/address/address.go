package address

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/admin"
	"github.com/fionera/TeamDriveManager/api/drive"
	"github.com/fionera/TeamDriveManager/api/iam"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewAssignAddressCmd() cli.Command {
	return cli.Command{
		Name:      "address",
		Usage:     "Assign an address to a specified teamdrive",
		Action:    CmdAssignAddress,
		Flags:     []cli.Flag{},
		UsageText: "<TEAMDRIVE> <ADDRESS> <TYPE>",
	}
}

func CmdAssignAddress(c *cli.Context) {

	client, err := api.CreateClient(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate, []string{iam.CloudPlatformScope, admin.AdminDirectoryGroupScope})
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
		if teamDrive.Name == c.Args().Get(0) {
			_, err := driveApi.CreatePermission(teamDrive.Id, "organizer", c.Args().Get(1), c.Args().Get(2))
			if err != nil {
				logrus.Error(err)
				return
			}

			logrus.Info("Added Permission")

			break
		}
	}
}
