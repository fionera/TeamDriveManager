package list

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewUnhideTeamDriveCommand() cli.Command {
	return cli.Command{
		Name:   "unhide",
		Usage:  "Unhide one or all TeamDrives",
		Action: CmdUnhideTeamDrive,
		Flags:  []cli.Flag{},
	}
}

func CmdUnhideTeamDrive(c *cli.Context) {
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

	var drives []string
	if driveId != "" {
		drives = []string{driveId}
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
			if !teamDrive.Hidden {
				continue
			}

			drives = append(drives, teamDrive.Id)
		}
	}

	for _, teamDrive := range drives {
		response, err := api.UnHideTeamDrive(driveApi, teamDrive)

		if err != nil {
			logrus.Error(err)
			continue
		}

		logrus.Infof("`%s``%t`", response.Name)
	}
}
