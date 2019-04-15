package teamdrive

import (
	"github.com/codegangsta/cli"
	"github.com/fionera/TeamdriveManager/api"
	"github.com/fionera/TeamdriveManager/api/drive"
	. "github.com/fionera/TeamdriveManager/config"
	"github.com/sirupsen/logrus"
	"strings"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "teamdrive",
		Usage:  "List all TeamDrives",
		Action: CmdListTeamDrive,
		Flags:  []cli.Flag{},
	}
}

func CmdListTeamDrive(c *cli.Context) {
	filter := strings.Join(c.Args(), " ")

	if filter != "" {
		logrus.Infof("Using filter `%s`", filter)
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

	teamDrives, err := driveApi.ListTeamDrives()
	if err != nil {
		logrus.Panic(err)
		return
	}

	var i int
	for _, teamDrive := range teamDrives {
		if !strings.HasPrefix(teamDrive.Name, filter) {
			continue
		}

		logrus.Infof("%s", teamDrive.Name)
		i++
	}

	logrus.Infof("Found %d TeamDrives", i)
}
