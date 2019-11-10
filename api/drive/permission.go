package drive

import (
	"context"

	"google.golang.org/api/drive/v3"
)

func (a *Api) ListPermissions(driveId string) ([]*drive.Permission, error) {
	return a.listPermissions(driveId, false)
}

func (a *Api) ListPermissionsAdmin(driveId string) ([]*drive.Permission, error) {
	return a.listPermissions(driveId, true)
}

func (a *Api) listPermissions(driveId string, admin bool) ([]*drive.Permission, error) {
	var permissions []*drive.Permission
	err := a.drive.Permissions.List(driveId).
		SupportsTeamDrives(true).
		Fields("permissions(id,emailAddress,domain,role,displayName)").
		UseDomainAdminAccess(admin).Pages(context.Background(), func(list *drive.PermissionList) error {
		permissions = append(permissions, list.Permissions...)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (a *Api) CreatePermission(driveId, role, address, addressType string) (*drive.Permission, error) {
	return a.createPermission(driveId, role, address, addressType, false)
}

func (a *Api) CreatePermissionAdmin(driveId, role, address, addressType string) (*drive.Permission, error) {
	return a.createPermission(driveId, role, address, addressType, true)
}

func (a *Api) createPermission(driveId, role, address, addressType string, admin bool) (*drive.Permission, error) {
	permission, err := a.drive.Permissions.Create(driveId, &drive.Permission{
		Role:         role,
		Type:         addressType,
		EmailAddress: address,
	}).SendNotificationEmail(false).
		UseDomainAdminAccess(admin).
		SupportsTeamDrives(true).Do()

	if err != nil {
		return nil, err
	}

	return permission, nil
}

func (a *Api) DeletePermission(driveId, permissionId string) error {
	return a.deletePermission(driveId, permissionId, false)
}

func (a *Api) DeletePermissionAdmin(driveId, permissionId string) error {
	return a.deletePermission(driveId, permissionId, true)
}

func (a *Api) deletePermission(driveId, permissionId string, admin bool) error {
	err := a.drive.Permissions.Delete(driveId, permissionId).
		UseDomainAdminAccess(admin).
		SupportsTeamDrives(true).Do()

	if err != nil {
		return err
	}

	return nil
}
