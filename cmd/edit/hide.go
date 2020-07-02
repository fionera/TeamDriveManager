package list

import (
	"github.com/AlecAivazis/survey/v2"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewHideTeamDriveCommand() cli.Command {
	return cli.Command{
		Name:   "hide",
		Usage:  "Hide one or all TeamDrives",
		Action: CmdHideTeamDrive,
		Flags:  []cli.Flag{},
	}
}

func CmdHideTeamDrive(c *cli.Context) {
	driveId := c.Args().First()

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

	if driveId != "" {
		response, err := api.HideTeamDrive(driveApi, driveId)

		if err != nil {
			logrus.Panic(err)
			return
		}

		logrus.Infof("`%s``%t`", response.Name, response.Hidden)
	} else {
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

		for _, teamDrive := range teamDrives {
			if teamDrive.Hidden {
				continue
			}
			response, err := api.HideTeamDrive(driveApi, teamDrive.Id)

			if err != nil {
				logrus.Error(err)
				continue
			}

			logrus.Infof("`%s``%t`", response.Name, response.Hidden)
		}
	}
}