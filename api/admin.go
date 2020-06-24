package api

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/jwt"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

func NewAdminService(jwtConfig *jwt.Config) (*admin.Service, error) {
	jwtConfig.Scopes = append(jwtConfig.Scopes, admin.AdminDirectoryGroupScope)
	return admin.NewService(context.TODO(), option.WithTokenSource(jwtConfig.TokenSource(context.TODO())))
}

func ListMembers(a *admin.Service, groupAddress string) ([]*admin.Member, error) {
	var members []*admin.Member

	err := a.Members.List(groupAddress).Pages(context.Background(), func(list *admin.Members) error {
		members = append(members, list.Members...)

		return nil
	})
	if err != nil {
		return nil, errors.Errorf("Error listing group members: %s", err)
	}

	return members, nil
}

func AddMember(a *admin.Service, groupAddress, memberAddress string) (*admin.Member, error) {
	return a.Members.Insert(groupAddress, &admin.Member{
		Email:            memberAddress,
		DeliverySettings: "NONE",
	}).Do()
}

func RemoveMember(a *admin.Service, groupAddress, memberAddress string) error {
	err := a.Members.Delete(groupAddress, memberAddress).Do()

	if err != nil {
		return errors.Errorf("Error removing group member: %s", err)
	}

	return nil
}

func CreateGroup(a *admin.Service, name, address string) (*admin.Group, error) {
	group, err := a.Groups.Insert(&admin.Group{
		Name:        name,
		Email:       address,
		Description: "Created by TeamDriveManager",
	}).Do()
	if err != nil {
		return nil, errors.Errorf("Error creating group: %s", err)
	}

	return group, nil
}

func ListGroups(a *admin.Service, domain string) ([]*admin.Group, error) {
	var groups []*admin.Group

	err := a.Groups.List().Domain(domain).Pages(context.Background(), func(list *admin.Groups) error {
		groups = append(groups, list.Groups...)

		return nil
	})

	if err != nil {
		return nil, errors.Errorf("Error listing groups: %s", err)
	}

	return groups, nil
}

// GroupExists needs the full email address as parameter and returns true if it can find a group for it
func GroupExists(a *admin.Service, address string) (bool, error) {
	_, err := a.Groups.Get(address).Do()
	if err != nil {
		if err.(*googleapi.Error).Code == 404 {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
