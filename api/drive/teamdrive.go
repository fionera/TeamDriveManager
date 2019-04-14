package drive

import (
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
