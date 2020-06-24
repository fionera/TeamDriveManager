package list

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewListGroupCommand() cli.Command {
	return cli.Command{
		Name:   "group",
		Usage:  "List all groups",
		Action: CmdListGroup,
		Flags:  []cli.Flag{},
	}
}

func CmdListGroup(c *cli.Context) {
	filter := strings.Join(c.Args(), " ")

	if filter != "" {
		logrus.Infof("Using filter `%s`", filter)
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
		if !strings.HasPrefix(group.Name, filter) {
			continue
		}

		logrus.Infof("%s", group.Name)
		i++
	}

	logrus.Infof("Found %d Groups", i)
}
