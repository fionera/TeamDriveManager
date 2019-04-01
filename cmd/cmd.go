package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

var GlobalFlags = []cli.Flag{}

var Commands []cli.Command

func RegisterCommand(command cli.Command) {
	Commands = append(Commands, command)
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
