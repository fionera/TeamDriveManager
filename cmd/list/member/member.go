package teamdrive

import (
	"github.com/codegangsta/cli"
	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/admin"
	. "github.com/fionera/TeamDriveManager/config"
	"github.com/sirupsen/logrus"
)

func NewCommand() cli.Command {
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

	client, err := api.CreateClient(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate, []string{admin.AdminDirectoryGroupScope})
	if err != nil {
		logrus.Error(err)
		return
	}

	adminApi, err := admin.NewApi(client)
	if err != nil {
		logrus.Error(err)
		return
	}

	groups, err := adminApi.ListGroups(App.AppConfig.Domain)
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
