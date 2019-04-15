package cloudresourcemanager

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/cloudresourcemanager/v1"
	"time"
)

func (a *Api) CreateProject(projectId, organization string) error {
	projectCreation, err := a.crm.Projects.Create(&cloudresourcemanager.Project{
		ProjectId: projectId,
		Name:      projectId,
		Parent: &cloudresourcemanager.ResourceId{
			Type: "organization",
			Id:   organization,
		},
	}).Do()
	if err != nil {
		return errors.Errorf("Error creating Project: %s", err)
	}

	logrus.Infof("Creating Project %s", projectId)

	for {
		operation, err := a.crm.Operations.Get(projectCreation.Name).Do()
		if err != nil {
			logrus.Panic(err)
			return err
		}

		if operation.Done {
			logrus.Infof("Creation finished")
			break
		} else {
			logrus.Infof("Creation still running. Polling again in 2 Seconds")
			time.Sleep(2 * time.Second)
		}
	}

	return nil
}

func (a *Api) ListProjects(organization string) ([]*cloudresourcemanager.Project, error) {
	var projects []*cloudresourcemanager.Project
	err := a.crm.Projects.List().Pages(context.Background(), func(list *cloudresourcemanager.ListProjectsResponse) error {
		projects = append(projects, list.Projects...)

		return nil
	})
	if err != nil {
		return nil, errors.Errorf("Error listing projects: %s", err)
	}

	return projects, nil
}

func (a *Api) DeleteProject(projectId string) error {
	_, err := a.crm.Projects.Delete(projectId).Do()

	if err != nil {
		return errors.Errorf("Error listing projects: %s", err)
	}

	return nil
}
