package api

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

func NewCloudResourceManagerService(jwtConfig *jwt.Config) (*cloudresourcemanager.Service, error) {
	jwtConfig.Scopes = append(jwtConfig.Scopes, cloudresourcemanager.CloudPlatformScope)
	return cloudresourcemanager.NewService(context.Background(), option.WithTokenSource(jwtConfig.TokenSource(context.TODO())))
}

func CreateProject(crm *cloudresourcemanager.Service, projectId, organization string) error {
	projectCreation, err := crm.Projects.Create(&cloudresourcemanager.Project{
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
		operation, err := crm.Operations.Get(projectCreation.Name).Do()
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

func ListProjects(crm *cloudresourcemanager.Service, organization string) ([]*cloudresourcemanager.Project, error) {
	var projects []*cloudresourcemanager.Project
	err := crm.Projects.List().Pages(context.Background(), func(list *cloudresourcemanager.ListProjectsResponse) error {
		projects = append(projects, list.Projects...)

		return nil
	})
	if err != nil {
		return nil, errors.Errorf("Error listing projects: %s", err)
	}

	return projects, nil
}

func DeleteProject(crm *cloudresourcemanager.Service, projectId string) error {
	_, err := crm.Projects.Delete(projectId).Do()

	if err != nil {
		return errors.Errorf("Error listing projects: %s", err)
	}

	return nil
}
