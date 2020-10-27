package delete

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewSpecialMoveCommand() cli.Command {
	return cli.Command{
		Name:   "special_move",
		Usage:  "move to MyDrive",
		Action: CmdSpecialMove,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "teamdrive-id",
			},
		},
	}
}

func CmdSpecialMove(c *cli.Context) {
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

	query := fmt.Sprintf(`parents = "%s"`, teamdriveId)
	driveFiles, err := api.ListAllObjects(driveApi, teamdriveId, query)
	if err != nil {
		logrus.Panic(err)
	}
}
