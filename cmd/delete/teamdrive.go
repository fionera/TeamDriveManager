package delete

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/api/googleapi"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewDeleteTeamDriveCommand() cli.Command {
	return cli.Command{
		Name:   "teamdrive",
		Usage:  "Delete a Teamdrive",
		Action: CmdDeleteTeamDrive,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "teamdrive-id",
			},
			cli.BoolFlag{
				Name: "force-delete",
			},
		},
	}
}

func CmdDeleteTeamDrive(c *cli.Context) {
	if !c.Args().Present() {
		logrus.Error("Please supply a teamdrive id")
		return
	}
	forceDelete := c.Bool("force-delete")
	teamdriveId := c.Args().First()

	tokenSource, err := api.NewTokenSource(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate)
	if err != nil {
		logrus.Error(err)
		return
	}

	driveApi, err := api.NewDriveService(tokenSource)
	if err != nil {
		logrus.Error(err)
		return
	}

	if forceDelete {
		driveFiles, err := api.ListAllObjects(driveApi, teamdriveId, "")
		if err != nil {
			logrus.Panic(err)
		}

		var fileDeleteRequests sync.WaitGroup
		var running int
		for i := 0; i < len(driveFiles); i++ {
			fileDeleteRequests.Add(1)
			running++

			go func(i int) {
				defer fileDeleteRequests.Done()

				err = api.DeleteObject(driveApi, driveFiles[i].Id)
				if err != nil {
					logrus.Debugf("Failed to delete object: %s", driveFiles[i].Id)
				}
				logrus.Infof("%05d: Deleted object: %s", i, driveFiles[i].Id)
			}(i)

			if running > App.Flags.Concurrency {
				fileDeleteRequests.Wait()
				running = 0
			}
		}

		api.EmptyTrash(driveApi)
	}
deleteTeamDrive:
	err = api.DeleteTeamDrive(driveApi, teamdriveId)
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok {
			switch gerr.Code {
			case 403:
				if forceDelete {
					logrus.Info("Waiting for all objects to finish deletion.")
					time.Sleep(10 * time.Second)
					goto deleteTeamDrive
				} else {
					logrus.Error("Teamdrive contains objects and therefore cannot be deleted.")
				}
				return
			default:
				logrus.Error("An error occurred when deleting account.", err)
				return
			}
		} else {
			logrus.Fatal("An unknown error occurred: ", err)
		}
	}

	logrus.Infof("Successfully deleted TeamDrive %s", teamdriveId)
}
