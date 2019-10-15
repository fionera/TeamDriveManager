package teamdrive

import (
	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/cloudresourcemanager"
	. "github.com/fionera/TeamDriveManager/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"strings"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "project",
		Usage:  "List all projects",
		Action: CmdListTeamDrive,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "organization",
			},
		},
	}
}

func CmdListTeamDrive(c *cli.Context) {
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

	client, err := api.CreateClient(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate, []string{cloudresourcemanager.CloudPlatformScope})
	if err != nil {
		logrus.Error(err)
		return
	}

	crmApi, err := cloudresourcemanager.NewApi(client)
	if err != nil {
		logrus.Error(err)
		return
	}

	projects, err := crmApi.ListProjects(App.AppConfig.Organization)
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
