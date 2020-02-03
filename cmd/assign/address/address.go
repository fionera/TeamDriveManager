package address

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/admin"
	"github.com/fionera/TeamDriveManager/api/drive"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewAssignAddressCmd() cli.Command {
	return cli.Command{
		Name:      "address",
		Usage:     "Assign an address to a specified teamdrive",
		Action:    CmdAssignAddress,
		Flags:     []cli.Flag{},
		UsageText: "<TEAMDRIVE-NAME> <ADDRESS> <TYPE> <ROLE>",
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func CmdAssignAddress(c *cli.Context) {

	const addressTypeUser string = "user"
	const addressTypeGroup string = "group"

	const roleTypeOrganizer string = "organizer"
	const roleTypeFileOrganizer string = "fileOrganizer"
	const roleTypeWriter string = "writer"
	const roleTypeCommenter string = "commenter"
	const roleTypeReader string = "reader"

	supportedTypes := []string{
		addressTypeUser,
		addressTypeGroup,
	}

	supportedRoles := []string{
		roleTypeOrganizer,
		roleTypeFileOrganizer,
		roleTypeWriter,
		roleTypeCommenter,
		roleTypeReader,
	}

	teamDriveName := c.Args().Get(0)
	address := c.Args().Get(1)
	addressType := c.Args().Get(2)
	role := c.Args().Get(3)

	if teamDriveName == "" {
		logrus.Error("Please supply a teamdrive name")
		return
	}

	if address == "" {
		logrus.Error("Please supply an address")
		return
	}

	if addressType == "" {
		logrus.Error("Please supply an address type (allowed: 'user' or 'group')")
		return
	} else {
		if !contains(supportedTypes, addressType) {
			logrus.Error("Unsupported type: '" + addressType + "' (allowed: 'user' or 'group')")
			return
		}
	}

	if role == "" {
		logrus.Info("No role supplied. Setting 'reader' permission...")
		role = "reader"
	} else {
		if !contains(supportedRoles, role) {
			logrus.Error("Unsupported role: '" + role + "' (allowed: 'organizer', 'fileOrganizer', 'writer', 'commenter', 'reader')")
			return
		}
	}

	client, err := api.CreateClient(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate, []string{drive.DriveScope, admin.AdminDirectoryGroupScope})
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
			_, err := driveApi.CreatePermission(teamDrive.Id, role, address, addressType)
			if err != nil {
				logrus.Error(err)
				return
			}

			logrus.Info("Added Permission")

			break
		}
	}
}
