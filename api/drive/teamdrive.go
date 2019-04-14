package drive

import (
	"google.golang.org/api/drive/v3"
	"math/rand"
	"strconv"
)

func (a *Api) CreateTeamDrive(name string) (*drive.TeamDrive, error) {
	teamDrive, err := a.drive.Teamdrives.Create(strconv.Itoa(rand.Int()), &drive.TeamDrive{
		Name: name,
	}).Do()

	if err != nil {
		return nil, err
	}

	return teamDrive, nil
}
