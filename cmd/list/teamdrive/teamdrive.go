package teamdrive

import (
	"github.com/codegangsta/cli"
	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/drive"
	. "github.com/fionera/TeamDriveManager/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
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

	boolResponse := false
	confirm := &survey.Confirm{
		Message: "Use Domain Admin access?",
		Default: false,
	}

	err = survey.AskOne(confirm, &boolResponse, nil)
	if err != nil {
		logrus.Panic(err)
		return
	}

	var list = driveApi.ListTeamDrives
	if boolResponse {
		list = driveApi.ListAllTeamDrives
	}

	teamDrives, err := list()
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
