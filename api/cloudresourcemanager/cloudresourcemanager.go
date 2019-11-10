package cloudresourcemanager

import (
	"net/http"

	"google.golang.org/api/cloudresourcemanager/v1"
)

const (
	CloudPlatformScope = cloudresourcemanager.CloudPlatformScope
)

type Api struct {
	crm *cloudresourcemanager.Service
}

func NewApi(client *http.Client) (*Api, error) {
	var api Api

	crm, err := cloudresourcemanager.New(client)
	if err != nil {

		return nil, err
	}

	api.crm = crm

	return &api, nil
}
