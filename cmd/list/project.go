package list

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewListProjectCommand() cli.Command {
	return cli.Command{
		Name:   "project",
		Usage:  "List all projects",
		Action: CmdListProject,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "organization",
			},
		},
	}
}

func CmdListProject(c *cli.Context) {
	organization := c.String("organization")

	if organization == "" {
		if App.AppConfig.Organization != "" {
			organization = App.AppConfig.Organization
		} else {
			logrus.Error("Please supply the Organization to use")
			return
		}
	}

	filter := strings.Join(c.Args(), " ")
	if filter != "" {
		logrus.Infof("Using filter `%s`", filter)
	}

	tokenSource, err := api.NewTokenSource(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate)
	if err != nil {
		logrus.Error(err)
		return
	}

	crmApi, err := api.NewCloudResourceManagerService(tokenSource)
	if err != nil {
		logrus.Error(err)
		return
	}

	projects, err := api.ListProjects(crmApi, App.AppConfig.Organization)
	if err != nil {
		logrus.Panic(err)
		return
	}

	var i int
	for _, project := range projects {
		if !strings.HasPrefix(project.Name, filter) {
			continue
		}

		logrus.Infof("Name: `%s` - ID: `%s` ", project.Name, project.ProjectId)
		i++
	}

	logrus.Infof("Found %d Projects", i)
}
