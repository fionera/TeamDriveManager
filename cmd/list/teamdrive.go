package list

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewListTeamDriveCommand() cli.Command {
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

	var boolResponse bool
	confirm := &survey.Confirm{
		Message: "Use Domain Admin access?",
		Default: false,
	}

	err = survey.AskOne(confirm, &boolResponse, nil)
	if err != nil {
		logrus.Panic(err)
		return
	}

	var list = api.ListTeamDrives
	if boolResponse {
		list = api.ListAllTeamDrives
	}

	teamDrives, err := list(driveApi)
	if err != nil {
		logrus.Panic(err)
		return
	}

	var i int
	for _, teamDrive := range teamDrives {
		if !strings.HasPrefix(teamDrive.Name, filter) {
			continue
		}

		logrus.Infof("`%s``%s`", teamDrive.Name, teamDrive.Id)
		i++
	}

	logrus.Infof("Found %d TeamDrives", i)
}
