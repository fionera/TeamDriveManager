package list

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewListMemberCommand() cli.Command {
	return cli.Command{
		Name:   "member",
		Usage:  "List all members of a group",
		Action: CmdListMember,
		Flags:  []cli.Flag{},
	}
}

func CmdListMember(c *cli.Context) {
	address := c.Args().First()

	if address == "" {
		logrus.Error("Please provide the group address")
		return
	}

	tokenSource, err := api.NewTokenSource(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate)
	if err != nil {
		logrus.Error(err)
		return
	}

	adminApi, err := api.NewAdminService(tokenSource)
	if err != nil {
		logrus.Error(err)
		return
	}

	groups, err := api.ListGroups(adminApi, App.AppConfig.Domain)
	if err != nil {
		logrus.Panic(err)
		return
	}

	var i int
	for _, group := range groups {
		logrus.Infof("%s", group.Name)
		i++
	}

	logrus.Infof("Found %d Groups", i)
}
