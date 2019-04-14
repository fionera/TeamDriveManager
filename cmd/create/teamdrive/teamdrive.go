package teamdrive

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fionera/TeamdriveManager/api"
	"github.com/fionera/TeamdriveManager/api/drive"
	"github.com/sirupsen/logrus"
)

const (
	scopePrefix = "https://www.googleapis.com/auth/"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "teamdrive",
		Usage:  "Create a Teamdrive",
		Action: CmdCreateTeamDrive,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "service-account-file",
			},
			cli.StringFlag{
				Name: "impersonate",
			},
		},
	}
}

func CmdCreateTeamDrive(c *cli.Context) {
	serviceAccountFile := c.String("service-account-file")
	impersonate := c.String("impersonate")

	client, err := api.CreateClient(serviceAccountFile, impersonate, []string{drive.DriveScope})
	if err != nil {
		logrus.Panic(err)
		return
	}

	driveApi, err := drive.NewApi(client)
	if err != nil {
		logrus.Panic(err)
		return
	}

	if !c.Args().Present() {
		fmt.Println("name?")
		return
	}

	_, err = driveApi.CreateTeamDrive(c.Args().Get(0))
	if err != nil {
		logrus.Panic(err)
		return
	}
}
