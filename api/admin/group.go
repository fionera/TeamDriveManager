package admin

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

func (a *Api) CreateGroup(name, address string) (*admin.Group, error) {
	group, err := a.admin.Groups.Insert(&admin.Group{
		Name:        name,
		Email:       address,
		Description: "Created by TeamDriveManager",
	}).Do()
	if err != nil {
		return nil, errors.Errorf("Error creating group: %s", err)
	}

	return group, nil
}

func (a *Api) ListGroups(domain string) ([]*admin.Group, error) {
	var groups []*admin.Group

	err := a.admin.Groups.List().Domain(domain).Pages(context.Background(), func(list *admin.Groups) error {
		groups = append(groups, list.Groups...)

		return nil
	})

	if err != nil {
		return nil, errors.Errorf("Error listing groups: %s", err)
	}

	return groups, nil
}

// GroupExists needs the full email address as parameter and returns true if it can find a group for it
func (a *Api) GroupExists(address string) (bool, error) {
	_, err := a.admin.Groups.Get(address).Do()
	if err != nil {
		if err.(*googleapi.Error).Code == 404 {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
