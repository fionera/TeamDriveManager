package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	. "github.com/fionera/TeamdriveManager/config"
	"os"
)

var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "config",
		Value:       "config.json",
		Usage:       "The Configfile to use",
		Destination: &App.ConfigFile,
	},
	cli.StringFlag{
		Name:        "service-account-file",
		Destination: &App.Flags.ServiceAccountFile,
	},
	cli.StringFlag{
		Name:        "impersonate",
		Destination: &App.Flags.Impersonate,
	},
	cli.IntFlag{
		Name:        "concurrency, c",
		Destination: &App.Flags.Concurrency,
	},
}

var Commands []cli.Command

func RegisterCommand(command cli.Command) {
	Commands = append(Commands, command)
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
