package setup

import (
	. "github.com/fionera/TeamDriveManager/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
)

func Setup() {
	boolResponse := false
	stringResponse := ""

	confirm := &survey.Confirm{
		Message: "Did you do the Setup from the Readme?",
	}

	err := survey.AskOne(confirm, &boolResponse, nil)
	if err != nil {
		logrus.Panic(err)
	}

	if !boolResponse {
		logrus.Info("Do it.")
		os.Exit(1)
	}

	input := &survey.Input{
		Message: "Whats your domain name?",
	}

	err = survey.AskOne(input, &stringResponse, nil)
	if err != nil {
		logrus.Panic(err)
	}

	App.AppConfig.Domain = stringResponse

	input = &survey.Input{
		Message: "Where should the Service Account files be saved?",
		Default: "Files",
	}

	err = survey.AskOne(input, &stringResponse, nil)
	if err != nil {
		logrus.Panic(err)
	}

	App.AppConfig.ServiceAccountFolder = stringResponse

	input = &survey.Input{
		Message: "How do you want the Service Account Group named?",
		Default: "serviceaccounts",
	}

	err = survey.AskOne(input, &stringResponse, nil)
	if err != nil {
		logrus.Panic(err)
	}

	App.AppConfig.ServiceAccountGroup = stringResponse

	input = &survey.Input{
		Message: "Whats the ID of your Organization? (You can see it in the API Console)",
	}

	err = survey.AskOne(input, &stringResponse, nil)
	if err != nil {
		logrus.Panic(err)
	}

	App.AppConfig.Organization = stringResponse

	input = &survey.Input{
		Message: "Where is the Service Account file stored? (Full Path)",
	}

	err = survey.AskOne(input, &stringResponse, nil)
	if err != nil {
		logrus.Panic(err)
	}

	App.AppConfig.ServiceAccountFile = stringResponse

	SaveConfig(App.AppConfig)

	input = &survey.Input{
		Message: "Whats your Username for this domain (Full Address)",
	}

	err = survey.AskOne(input, &stringResponse, nil)
	if err != nil {
		logrus.Panic(err)
	}

	App.AppConfig.Impersonate = stringResponse

	SaveConfig(App.AppConfig)
}
