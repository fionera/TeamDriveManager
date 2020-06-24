package combinations

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/Jeffail/gabs"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	iam2 "google.golang.org/api/iam/v1"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewRegenerateKeysCommand() cli.Command {
	return cli.Command{
		Name:   "regenerate_keys",
		Usage:  "Recreates JSON files used for Service Accounts",
		Action: CmdRegenerateKeys,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "project-id",
			},
			cli.StringFlag{
				Name: "organization",
			},
		},
	}
}

func CmdRegenerateKeys(c *cli.Context) {
	projectId := c.Args().First()
	organization := c.String("organization")

	if projectId == "" {
		logrus.Error("Please supply the ProjectID to use")
		return
	}

	if organization == "" {
		if App.AppConfig.Organization != "" {
			organization = App.AppConfig.Organization
		} else {
			logrus.Error("Please supply the Organization to use")
			return
		}
	}

	tokenSource, err := api.NewTokenSource(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate)
	if err != nil {
		logrus.Panic(err)
		return
	}

	crmApi, err := api.NewCloudResourceManagerService(tokenSource)
	if err != nil {
		logrus.Panic(err)
		return
	}

	iamApi, err := api.NewIAMService(tokenSource)
	if err != nil {
		logrus.Panic(err)
		return
	}

	logrus.Info("Listing Projects")
	projects, err := api.ListProjects(crmApi, organization)
	if err != nil {
		logrus.Panic(err)
		return
	}

	var found bool
	for _, project := range projects {
		if project.ProjectId == projectId {
			found = true
		}
	}

	if !found {
		logrus.Panicf("Project %s was not found", projectId)
		return
	}

	var serviceAccountRequests sync.WaitGroup
	var running int
	serviceAccounts, err := api.ListServiceAccounts(iamApi, projectId)
	for _, serviceAccount := range serviceAccounts {
		serviceAccountRequests.Add(1)
		running++

		go func(account *iam2.ServiceAccount) {
			defer serviceAccountRequests.Done()
		getServiceAccount:
			serviceAccountObject, err := api.GetServiceAccount(iamApi, projectId, account.Email)
			if err != nil {
				logrus.Error(err)
				goto getServiceAccount
			}
		deleteApiKey:
			logrus.Infof("Deleting key for `%s`", account.Email)
			_, err = api.DeleteServiceAccountKey(iamApi, projectId, account.Email)
			if err != nil {
				logrus.Error(err)
				goto deleteApiKey
			}
		createApiKey:
			logrus.Infof("Creating new key for `%s`", account.Email)
			serviceAccountKey, err := api.CreateServiceAccountKey(iamApi, serviceAccountObject)
			if err != nil {
				logrus.Error(err)
				goto createApiKey
			}

			json, err := serviceAccountKey.MarshalJSON()
			if err != nil {
				logrus.Panic(err)
				return
			}

			container, err := gabs.ParseJSON(json)
			if err != nil {
				logrus.Panicf("Error parsing JSON: %s", err)
				return
			}

			privateKeyData := container.Path("privateKeyData").String()
			jsonData, err := base64.StdEncoding.DecodeString(privateKeyData[1 : len(privateKeyData)-1])
			if err != nil {
				logrus.Panicf("Error reading key: %s", err)
				return
			}

			err = os.Mkdir(App.AppConfig.ServiceAccountFolder, 0755)
			if err != nil && !os.IsExist(err) {
				logrus.Panicf("Error changing type: %s", err)
				return
			}

			err = ioutil.WriteFile(App.AppConfig.ServiceAccountFolder+"/"+serviceAccount.ProjectId+"_"+strings.ReplaceAll(serviceAccount.DisplayName, " ", "_")+".json", jsonData, 0755)
			if err != nil {
				logrus.Panic(err)
				return
			}
		}(serviceAccount)

		if running > App.Flags.Concurrency {
			serviceAccountRequests.Wait()
			running = 0
		}
	}

	serviceAccountRequests.Wait()

	logrus.Infof("Done :3")

	return
}
