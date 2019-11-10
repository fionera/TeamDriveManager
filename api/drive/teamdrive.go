package drive

import (
	"context"

	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
)

func (a *Api) CreateTeamDrive(name string) (*drive.Drive, error) {
	id := uuid.NewV4()

	teamDrive, err := a.drive.Drives.Create(id.String(), &drive.Drive{
		Name: name,
	}).Do()

	if err != nil {
		return nil, err
	}

	return teamDrive, nil
}

func (a *Api) ListTeamDrives() ([]*drive.Drive, error) {
	return a.listTeamDrives(false)
}

func (a *Api) ListAllTeamDrives() ([]*drive.Drive, error) {
	return a.listTeamDrives(true)
}

func (a *Api) listTeamDrives(admin bool) ([]*drive.Drive, error) {
	var teamDrives []*drive.Drive

	logrus.Debugf("Getting Drive List with admin: %t", admin)
	err := a.drive.Drives.List().UseDomainAdminAccess(admin).PageSize(100).Pages(context.Background(), func(list *drive.DriveList) error {
		teamDrives = append(teamDrives, list.Drives...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return teamDrives, nil
}
