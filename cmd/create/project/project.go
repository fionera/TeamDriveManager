package project

import (
	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/cloudresourcemanager"
	"github.com/fionera/TeamDriveManager/api/servicemanagement"
	. "github.com/fionera/TeamDriveManager/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func NewCommand() cli.Command {
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

	client, err := api.CreateClient(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate, []string{cloudresourcemanager.CloudPlatformScope, servicemanagement.ServiceManagementScope})
	if err != nil {
		logrus.Panic(err)
		return
	}

	crmApi, err := cloudresourcemanager.NewApi(client)
	if err != nil {
		logrus.Panic(err)
		return
	}

	smApi, err := servicemanagement.NewApi(client)
	if err != nil {
		logrus.Panic(err)
		return
	}

	err = crmApi.CreateProject(projectId, organization)
	if err != nil {
		logrus.Panic(err)
		return
	}

	err = smApi.EnableApi("project:"+projectId, servicemanagement.DriveApi)
	if err != nil {
		logrus.Panic(err)
		return
	}
}
