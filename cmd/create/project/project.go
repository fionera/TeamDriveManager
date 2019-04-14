package project

import (
	"github.com/codegangsta/cli"
	"github.com/fionera/TeamdriveManager/api"
	"github.com/fionera/TeamdriveManager/api/cloudresourcemanager"
	"github.com/fionera/TeamdriveManager/api/servicemanagement"
	"github.com/sirupsen/logrus"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "project",
		Usage:  "Create a Project",
		Action: CmdCreateProject,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "service-account-file",
			},
			cli.StringFlag{
				Name: "impersonate",
			},
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
	serviceAccountFile := c.String("service-account-file")
	impersonate := c.String("impersonate")
	projectId := c.String("project-id")
	organization := c.String("organization")

	client, err := api.CreateClient(serviceAccountFile, impersonate, []string{cloudresourcemanager.CloudPlatformScope, servicemanagement.ServiceManagementScope})
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
