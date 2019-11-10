package teamdrive

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/cloudresourcemanager"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "project",
		Usage:  "Delete selected Projects",
		Action: CmdDeleteProject,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "organization",
			},
		},
	}
}

func CmdDeleteProject(c *cli.Context) {
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

	var projectIds []string
	for _, project := range projects {
		if !strings.HasPrefix(project.Name, filter) || project.LifecycleState != "ACTIVE" {
			continue
		}

		projectIds = append(projectIds, project.ProjectId)
	}

	var toDelete []string
	prompt := &survey.MultiSelect{
		Message: "Which Projects to delete?",
		Options: projectIds,
	}
	err = survey.AskOne(prompt, &toDelete, nil)
	if err != nil {
		logrus.Panic(err)
		return
	}

	for _, projectId := range toDelete {
	deleteProject:
		err := crmApi.DeleteProject(projectId)
		logrus.Infof("Deleted Project `%s`", projectId)
		if err != nil {
			logrus.Error(err)
			goto deleteProject
		}
	}

	logrus.Infof("Done :3")
}
