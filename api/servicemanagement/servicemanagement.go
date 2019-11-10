package servicemanagement

import (
	"net/http"

	"google.golang.org/api/servicemanagement/v1"
)

const (
	DriveApi               = "drive.googleapis.com"
	ServiceManagementScope = servicemanagement.ServiceManagementScope
	CloudPlatformScope     = servicemanagement.CloudPlatformScope
)

type Api struct {
	sm *servicemanagement.APIService
}

func NewApi(client *http.Client) (*Api, error) {
	var api Api

	sm, err := servicemanagement.New(client)
	if err != nil {
		return nil, err
	}

	api.sm = sm

	return &api, nil
}
