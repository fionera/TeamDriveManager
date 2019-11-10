package iam

import (
	"net/http"

	"google.golang.org/api/iam/v1"
)

const (
	CloudPlatformScope = iam.CloudPlatformScope
)

type Api struct {
	iam *iam.Service
}

func NewApi(client *http.Client) (*Api, error) {
	var api Api

	sm, err := iam.New(client)
	if err != nil {
		return nil, err
	}

	api.iam = sm

	return &api, nil
}
