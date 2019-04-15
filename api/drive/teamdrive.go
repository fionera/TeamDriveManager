package drive

import (
	"context"
	"github.com/satori/go.uuid"
	"google.golang.org/api/drive/v3"
)

func (a *Api) CreateTeamDrive(name string) (*drive.TeamDrive, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	teamDrive, err := a.drive.Teamdrives.Create(id.String(), &drive.TeamDrive{
		Name: name,
	}).Do()

	if err != nil {
		return nil, err
	}

	return teamDrive, nil
}

func (a *Api) ListTeamDrives() ([]*drive.TeamDrive, error) {
	var teamDrives []*drive.TeamDrive

	err := a.drive.Teamdrives.List().PageSize(100).Pages(context.Background(), func(list *drive.TeamDriveList) error {
		teamDrives = append(teamDrives, list.TeamDrives...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return teamDrives, nil
}
