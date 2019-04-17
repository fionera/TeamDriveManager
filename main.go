package main

import (
	"github.com/fionera/TeamDriveManager/config"
	"github.com/fionera/TeamDriveManager/setup"
	"github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"

	"github.com/codegangsta/cli"
	. "github.com/fionera/TeamDriveManager/cmd"
	_ "github.com/fionera/TeamDriveManager/cmd/all"
)

func main() {

	app := cli.NewApp()
	app.Name = "TeamDriveManager"
	app.Version = "0.1.0"
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
		config.SaveConfig(config.App.AppConfig) //TODO: When panic on load destroys config

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Panic(err)
	}
}
