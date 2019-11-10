package admin

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/admin/directory/v1"
)

func (a *Api) ListMembers(groupAddress string) ([]*admin.Member, error) {
	var members []*admin.Member

	err := a.admin.Members.List(groupAddress).Pages(context.Background(), func(list *admin.Members) error {
		members = append(members, list.Members...)

		return nil
	})
	if err != nil {
		return nil, errors.Errorf("Error listing group members: %s", err)
	}

	return members, nil
}

func (a *Api) AddMember(groupAddress, memberAddress string) (*admin.Member, error) {
	member, err := a.admin.Members.Insert(groupAddress, &admin.Member{
		Email:            memberAddress,
		DeliverySettings: "NONE",
	}).Do()

	if err != nil {
		return nil, err
	}

	return member, nil
}

func (a *Api) RemoveMember(groupAddress, memberAddress string) error {
	err := a.admin.Members.Delete(groupAddress, memberAddress).Do()

	if err != nil {
		return errors.Errorf("Error removing group member: %s", err)
	}

	return nil
}
