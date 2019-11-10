package admin

import (
	"net/http"

	"google.golang.org/api/admin/directory/v1"
)

const (
	AdminDirectoryGroupScope = admin.AdminDirectoryGroupScope
)

type Api struct {
	admin *admin.Service
}

func NewApi(client *http.Client) (*Api, error) {
	var api Api

	adminService, err := admin.New(client)
	if err != nil {
		return nil, err
	}

	api.admin = adminService

	return &api, nil
}
