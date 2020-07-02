package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"

	. "github.com/fionera/TeamDriveManager/cmd"
	_ "github.com/fionera/TeamDriveManager/cmd/all"
	"github.com/fionera/TeamDriveManager/config"
	"github.com/fionera/TeamDriveManager/setup"
)

// Version - defined default version if it's not passed through flags during build
var Version string = "master"

func main() {

	app := cli.NewApp()
	app.Name = "TeamDriveManager"
	app.Version = Version
	app.Author = "fionera"
	app.Email = "teamdrive-manager@fionera.de"
	app.Usage = ""

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound
	app.Before = func(context *cli.Context) error {
		config.LoadConfig()

		if config.App.AppConfig.ServiceAccountFile == "" || config.App.AppConfig.Domain == "" {
			cont := false
			prompt := &survey.Confirm{
				Message: "It seems TeamDriveManager isn't configured correctly. Start Setup?",
			}
			err := survey.AskOne(prompt, &cont, nil)
			if err != nil {
				return err
			}

			if !cont {
				logrus.Info("Exiting.")
				os.Exit(1)
			}

			setup.Setup()
		}

		if context.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}

		return nil
	}
	app.After = func(context *cli.Context) error {
		if config.App.AppConfig.ServiceAccountFile != "" {
			config.SaveConfig(config.App.AppConfig)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Panic(err)
	}
}
