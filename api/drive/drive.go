package drive

import (
	"net/http"

	"google.golang.org/api/drive/v3"
)

const (
	DriveScope = drive.DriveScope
)

var Roles = []string{"owner", "organizer", "fileOrganizer", "writer", "commenter", "reader"}

type Api struct {
	drive *drive.Service
}

func NewApi(client *http.Client) (*Api, error) {
	var api Api

	driveService, err := drive.New(client)
	if err != nil {
		return nil, err
	}

	api.drive = driveService

	return &api, nil
}
