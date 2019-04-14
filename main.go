package main

import (
	"github.com/fionera/TeamdriveManager/config"
	"os"

	"github.com/codegangsta/cli"
	. "github.com/fionera/TeamdriveManager/cmd"
	_ "github.com/fionera/TeamdriveManager/cmd/all"
)

func main() {

	app := cli.NewApp()
	app.Name = "TeamdriveManager"
	app.Version = "0.1.0"
	app.Author = "fionera"
	app.Email = "teamdrive-manager@fionera.de"
	app.Usage = ""

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound
	app.Before = func(context *cli.Context) error {
		config.LoadConfig()

		return nil
	}

	app.Run(os.Args)
}
