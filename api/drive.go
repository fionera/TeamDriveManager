package api

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var DriveRoles = []string{"owner", "organizer", "fileOrganizer", "writer", "commenter", "reader"}

func NewDriveService(jwtConfig *jwt.Config) (*drive.Service, error) {
	jwtConfig.Scopes = append(jwtConfig.Scopes, drive.DriveScope)
	return drive.NewService(context.TODO(), option.WithTokenSource(jwtConfig.TokenSource(context.TODO())))
}

func ListPermissions(driveApi *drive.Service, driveId string) ([]*drive.Permission, error) {
	return listPermissions(driveApi, driveId, false)
}

func ListPermissionsAdmin(driveApi *drive.Service, driveId string) ([]*drive.Permission, error) {
	return listPermissions(driveApi, driveId, true)
}

func listPermissions(driveApi *drive.Service, driveId string, admin bool) ([]*drive.Permission, error) {
	var permissions []*drive.Permission
	err := driveApi.Permissions.List(driveId).
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

func CreatePermission(driveApi *drive.Service, driveId, role, address, addressType string) (*drive.Permission, error) {
	return createPermission(driveApi, driveId, role, address, addressType, false)
}

func CreatePermissionAdmin(driveApi *drive.Service, driveId, role, address, addressType string) (*drive.Permission, error) {
	return createPermission(driveApi, driveId, role, address, addressType, true)
}

func createPermission(driveApi *drive.Service, driveId, role, address, addressType string, admin bool) (*drive.Permission, error) {
	permission, err := driveApi.Permissions.Create(driveId, &drive.Permission{
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

func DeletePermission(driveApi *drive.Service, driveId, permissionId string) error {
	return deletePermission(driveApi, driveId, permissionId, false)
}

func DeletePermissionAdmin(driveApi *drive.Service, driveId, permissionId string) error {
	return deletePermission(driveApi, driveId, permissionId, true)
}

func deletePermission(driveApi *drive.Service, driveId, permissionId string, admin bool) error {
	err := driveApi.Permissions.Delete(driveId, permissionId).
		UseDomainAdminAccess(admin).
		SupportsTeamDrives(true).Do()

	if err != nil {
		return err
	}

	return nil
}

func CreateTeamDrive(driveApi *drive.Service, name string) (*drive.Drive, error) {
	id := uuid.NewV4()

	teamDrive, err := driveApi.Drives.Create(id.String(), &drive.Drive{
		Name: name,
	}).Do()

	if err != nil {
		return nil, err
	}

	return teamDrive, nil
}

func ListTeamDrives(driveApi *drive.Service) ([]*drive.Drive, error) {
	return listTeamDrives(driveApi, false)
}

func ListAllTeamDrives(driveApi *drive.Service) ([]*drive.Drive, error) {
	return listTeamDrives(driveApi, true)
}

func listTeamDrives(driveApi *drive.Service, admin bool) ([]*drive.Drive, error) {
	var teamDrives []*drive.Drive

	logrus.Debugf("Getting Drive List with admin: %t", admin)
	// * field because Google Api is fucking broken and doesn't return hidden
	err := driveApi.Drives.List().UseDomainAdminAccess(admin).Fields("*").PageSize(100).Pages(context.Background(), func(list *drive.DriveList) error {
		teamDrives = append(teamDrives, list.Drives...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return teamDrives, nil
}

func ListAllObjects(driveApi *drive.Service, driveId string) ([]*drive.File, error) {
	var fileList []*drive.File

	logrus.Debug("Retrieving all files.")
	err := driveApi.Files.List().IncludeItemsFromAllDrives(true).SupportsAllDrives(true).Corpora("drive").DriveId(driveId).PageSize(1000).Pages(context.Background(), func(list *drive.FileList) error {
		fileList = append(fileList, list.Files...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}

func DeleteObject(driveApi *drive.Service, fileId string) error {
	return driveApi.Files.Delete(fileId).SupportsAllDrives(true).Do()
}

func EmptyTrash(driveApi *drive.Service) *drive.FilesEmptyTrashCall {
	return driveApi.Files.EmptyTrash()
}

func DeleteTeamDrive(driveApi *drive.Service, id string) error {
	return driveApi.Drives.Delete(id).Do()
}

func HideTeamDrive(driveApi *drive.Service, driveId string) (*drive.Drive, error) {
	return driveApi.Drives.Hide(driveId).Do()
}

func UnHideTeamDrive(driveApi *drive.Service, driveId string) (*drive.Drive, error) {
	return driveApi.Drives.Unhide(driveId).Do()
}
