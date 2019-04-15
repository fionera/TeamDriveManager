package config

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type Flags struct {
	AppConfig
	Concurrency int
}

type AppConfig struct {
	ServiceAccountFile   string
	ServiceAccountFolder string // The Folder where to save the ServiceAccount Files
	Impersonate          string
	Organization         string // The Organization where the API Projects are created
	Domain               string // The domain under which groups are created

	Projects            []string // The ProjectIds of all Projects that should be considered to be used for Service Accounts
	ServiceAccountGroup string   // The address of the service account group
}

type Application struct {
	ConfigFile string
	AppConfig  AppConfig
	Flags      Flags
}

var App *Application

func init() {
	App = &Application{}
}

func SaveConfig(config AppConfig) {
	bytes, err := json.Marshal(config)
	if err != nil {
		logrus.Panic(err)
		return
	}

	err = ioutil.WriteFile(App.ConfigFile, bytes, 0755)
	if err != nil {
		logrus.Panic(err)
		return
	}
}

func CreateDefaultConfig() {
	SaveConfig(AppConfig{
		Projects: []string{},
	})
}

func LoadConfig() {
	logrus.Debugf("Loading Configfile: %s", App.ConfigFile)
	content, err := ioutil.ReadFile(App.ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Infof("Configfile doesnt exist. Creating empty one")
			CreateDefaultConfig()
			return
		}

		logrus.Panic(err)
		return
	}

	var config AppConfig
	err = json.Unmarshal(content, &config)
	if err != nil {
		logrus.Panic(err)
		return
	}

	if App.Flags.ServiceAccountGroup == "" {
		config.ServiceAccountGroup = "serviceaccounts"
	}

	if App.Flags.ServiceAccountFile != "" {
		config.ServiceAccountFile = App.Flags.ServiceAccountFile
	}

	if App.Flags.Impersonate != "" {
		config.Impersonate = App.Flags.Impersonate
	}

	if App.Flags.ServiceAccountFolder != "" {
		config.ServiceAccountFolder = App.Flags.ServiceAccountFolder
	}

	App.AppConfig = config
}
