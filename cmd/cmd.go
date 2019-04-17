package cmd

import (
	"github.com/codegangsta/cli"
	. "github.com/fionera/TeamDriveManager/config"
	"github.com/sirupsen/logrus"
	"os"
)

var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "config",
		Value:       "config.json",
		Usage:       "The Configfile to use",
		Destination: &App.AppConfigFile,
	},
	cli.StringFlag{
		Name:        "service-account-file",
		Destination: &App.Flags.ServiceAccountFile,
	},
	cli.StringFlag{
		Name:        "impersonate",
		Destination: &App.Flags.Impersonate,
	},
	cli.StringFlag{
		Name:        "service-account-folder",
		Destination: &App.Flags.ServiceAccountFolder,
	},
	cli.IntFlag{
		Name:        "concurrency, c",
		Destination: &App.Flags.Concurrency,
	},
	cli.BoolFlag{
		Name:        "debug",
		Destination: &App.Flags.Debug,
	},
}

var Commands []cli.Command

func RegisterCommand(command cli.Command) {
	Commands = append(Commands, command)
}

func CommandNotFound(c *cli.Context, command string) {
	logrus.Errorf("%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
