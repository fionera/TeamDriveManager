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
	ServiceAccountFile string
	Impersonate        string
	Organization       string // The Organization where the API Projects are created
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
	SaveConfig(AppConfig{})
}

func LoadConfig() {
	logrus.Debugf("Loading Configfile: %s", App.ConfigFile)
	content, err := ioutil.ReadFile(App.ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Infof("Configfile doesnt exist. Creating default one. Please edit it to your needs")
			CreateDefaultConfig()
			os.Exit(0)
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

	if config.ServiceAccountFile == "" {
		logrus.Errorf("Please edit the config '%s' and add the Service Account File", App.ConfigFile)
		os.Exit(1)
	}

	SaveConfig(config)

	if App.Flags.ServiceAccountFile != "" {
		config.ServiceAccountFile = App.Flags.ServiceAccountFile
	}

	if App.Flags.Impersonate != "" {
		config.Impersonate = App.Flags.Impersonate
	}

	App.AppConfig = config
}
