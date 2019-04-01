package operations

import (
	admin "google.golang.org/api/admin/directory/v1"
)

var adminService *admin.Service

func AddUserToGroup(groupKey string, mailAddress string) (*admin.Member, error) {
	return adminService.Members.Insert(groupKey, &admin.Member{
		Email:            mailAddress,
		DeliverySettings: "NONE",
		Role:             "MEMBER",
	}).Do()
}