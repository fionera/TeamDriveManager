package create

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewCreateProjectCommand() cli.Command {
	return cli.Command{
		Name:   "project",
		Usage:  "Create a Project",
		Action: CmdCreateProject,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "project-id",
			},
			cli.StringFlag{
				Name: "organization",
			},
		},
	}
}

func CmdCreateProject(c *cli.Context) {
	projectId := c.String("project-id")
	organization := c.String("organization")

	if projectId == "" {
		logrus.Error("Please supply the ProjectID to use")
		return
	}

	if organization == "" {
		if App.AppConfig.Organization != "" {
			organization = App.AppConfig.Organization
		} else {
			logrus.Error("Please supply the Organization to use")
			return
		}
	}

	tokenSource, err := api.NewTokenSource(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate)
	if err != nil {
		logrus.Panic(err)
		return
	}

	crmApi, err := api.NewCloudResourceManagerService(tokenSource)
	if err != nil {
		logrus.Panic(err)
		return
	}

	smApi, err := api.NewServiceManagementService(tokenSource)
	if err != nil {
		logrus.Panic(err)
		return
	}

	err = api.CreateProject(crmApi, projectId, organization)
	if err != nil {
		logrus.Panic(err)
		return
	}

	err = api.EnableApi(smApi, "project:"+projectId, api.DriveApi)
	if err != nil {
		logrus.Panic(err)
		return
	}
}
