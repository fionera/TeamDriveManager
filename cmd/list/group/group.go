package teamdrive

import (
	"github.com/codegangsta/cli"
	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/admin"
	. "github.com/fionera/TeamDriveManager/config"
	"github.com/sirupsen/logrus"
	"strings"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "group",
		Usage:  "List all groups",
		Action: CmdListTeamDrive,
		Flags:  []cli.Flag{},
	}
}

func CmdListTeamDrive(c *cli.Context) {
	filter := strings.Join(c.Args(), " ")

	if filter != "" {
		logrus.Infof("Using filter `%s`", filter)
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
		if !strings.HasPrefix(group.Name, filter) {
			continue
		}

		logrus.Infof("%s", group.Name)
		i++
	}

	logrus.Infof("Found %d Groups", i)
}
