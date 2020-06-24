package api

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"google.golang.org/api/servicemanagement/v1"
)

const (
	DriveApi = "drive.googleapis.com"
)

func NewServiceManagementService(jwtConfig *jwt.Config) (*servicemanagement.APIService, error) {
	jwtConfig.Scopes = append(jwtConfig.Scopes, servicemanagement.CloudPlatformScope)
	return servicemanagement.NewService(context.TODO(), option.WithTokenSource(jwtConfig.TokenSource(context.TODO())))
}

func EnableApi(serviceManagementApi *servicemanagement.APIService, consumerId, serviceName string) error {
	logrus.Infof("Enabling %s API", serviceName)

	operation, err := serviceManagementApi.Services.Enable(serviceName, &servicemanagement.EnableServiceRequest{
		ConsumerId: consumerId,
	}).Do()
	if err != nil {
		return err
	}

	for {
		operation, err := serviceManagementApi.Operations.Get(operation.Name).Do()
		if err != nil {
			return err
		}

		if operation.Done {
			logrus.Infof("Enabled %s API", serviceName)
			break
		} else {
			logrus.Infof("Enabling still running. Polling again in 2 Seconds")
			time.Sleep(2 * time.Second)
		}
	}

	return nil
}
