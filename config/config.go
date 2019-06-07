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
	Debug       bool
}

type GroupAssignment struct {
	TeamDriveName  string            // The TeamDrive name without the Prefix
	GroupAddresses map[string]string // A map which has the role name as key and the address as value
}

type RawUsers map[string]map[string]string

type TeamDriveConfig struct {
	NamePrefix string // The string with which every TeamDrive has to begin with

	GlobalUsers RawUsers // A map which has the role name as Key and a map where the users  that are added to all TeamDrives if not overridden anywhere as key and a comment the value

	BlackList map[string][]string // A Map which has the TeamDrive names without the in NamePrefix defined prefix as key and an array of email addresses as value

	GroupAssignments []GroupAssignment // An array which contains one or multiple GroupAssignments
}

type AppConfig struct {
	ServiceAccountFile   string
	ServiceAccountFolder string // The Folder where to save the ServiceAccount Files
	Impersonate          string
	Organization         string // The Organization where the API Projects are created
	Domain               string // The domain under which groups are created

	Projects            []string // The ProjectIds of all Projects that should be considered to be used for Service Accounts
	ServiceAccountGroup string   // The address of the service account group

	TeamDriveConfig TeamDriveConfig
}

type Application struct {
	AppConfigFile string
	AppConfig     AppConfig
	Flags         Flags
}

var App *Application

func init() {
	App = &Application{}
}

func SaveConfig(config interface{}) {
	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		logrus.Panic(err)
		return
	}

	var savePath string
	switch config.(type) {
	case AppConfig:
		savePath = App.AppConfigFile
		break
	default:
		logrus.Panic("unknown config type")
		break
	}

	err = ioutil.WriteFile(savePath, bytes, 0755)
	if err != nil {
		logrus.Panic(err)
		return
	}
}

func CreateDefaultAppConfig() {
	SaveConfig(AppConfig{
		Domain:   "domain.com",
		Projects: []string{},
		TeamDriveConfig: TeamDriveConfig{
			GlobalUsers: map[string]map[string]string{
				"reader": {
					"you@domain.com":      "myself",
					"example@example.org": "random dude",
				},
			},
			BlackList: map[string][]string{
				"Example TeamDrive Name": {
					"example@example.org",
				},
			},
			GroupAssignments: []GroupAssignment{
				{
					TeamDriveName: "Example TeamDrive Name",
					GroupAddresses: map[string]string{
						"writer": "example_teamdrive_name_writer@domain.com",
					},
				},
			},
		},
	})
}

func LoadConfig() {
	logrus.Debugf("Loading AppConfig from %s", App.AppConfigFile)
	appConfigContent, err := ioutil.ReadFile(App.AppConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Infof("Configfile doesnt exist. Creating empty one")
			CreateDefaultAppConfig()
			return
		}

		logrus.Panic(err)
		return
	}

	var appConfig AppConfig
	err = json.Unmarshal(appConfigContent, &appConfig)
	if err != nil {
		logrus.Panic(err)
		return
	}

	if App.Flags.ServiceAccountGroup != "" {
		appConfig.ServiceAccountGroup = "serviceaccounts"
	}

	if App.Flags.ServiceAccountFile != "" {
		appConfig.ServiceAccountFile = App.Flags.ServiceAccountFile
	}

	if App.Flags.Impersonate != "" {
		appConfig.Impersonate = App.Flags.Impersonate
	}

	if App.Flags.ServiceAccountFolder != "" {
		appConfig.ServiceAccountFolder = App.Flags.ServiceAccountFolder
	}

	App.AppConfig = appConfig
}
