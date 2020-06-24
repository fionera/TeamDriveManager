package create

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewCreateGroupCommand() cli.Command {
	return cli.Command{
		Name:   "config",
		Usage:  "Create a Group",
		Action: CmdCreateGroup,
		Flags:  []cli.Flag{},
	}
}

func CmdCreateGroup(c *cli.Context) {
	tokenSource, err := api.NewTokenSource(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate)
	if err != nil {
		logrus.Panic(err)
		return
	}

	adminApi, err := api.NewAdminService(tokenSource)
	if err != nil {
		logrus.Panic(err)
		return
	}

	name := c.Args().First()
	var address string
	if name == "" {
		logrus.Errorf("Please supply a name")
		return
	}

	if c.Args().Get(1) != "" {
		address = name
		name = c.Args().Get(1)
	} else {
		address = strings.ReplaceAll(name, " ", "_")
	}

	address = strings.ToLower(address)

	if !strings.Contains(address, "@") {
		address += "@" + App.AppConfig.Domain
	}

	logrus.Infof("Creating Group: %s<%s>", name, address)
	group, err := api.CreateGroup(adminApi, name, address)
	if err != nil {
		logrus.Panic(err)
		return
	}

	logrus.Infof("Successfully created Group: %s<%s>", group.Name, group.Email)
}
