package combinations

import (
	"encoding/base64"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/cloudresourcemanager"
	"github.com/fionera/TeamDriveManager/api/iam"
	"github.com/fionera/TeamDriveManager/api/servicemanagement"
	"github.com/fionera/TeamDriveManager/cmd/assign/serviceaccount"
	. "github.com/fionera/TeamDriveManager/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

func NewProjectAccountsKeysCommand() cli.Command {
	return cli.Command{
		Name:   "project_accounts_keys",
		Usage:  "Create a Project, fill it with 100 Accounts and create the Keys for it",
		Action: CmdCreateProjectAccountsKeys,
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

func CmdCreateProjectAccountsKeys(c *cli.Context) {
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

	client, err := api.CreateClient(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate, []string{cloudresourcemanager.CloudPlatformScope, servicemanagement.ServiceManagementScope})
	if err != nil {
		logrus.Panic(err)
		return
	}

	crmApi, err := cloudresourcemanager.NewApi(client)
	if err != nil {
		logrus.Panic(err)
		return
	}

	smApi, err := servicemanagement.NewApi(client)
	if err != nil {
		logrus.Panic(err)
		return
	}

	iamApi, err := iam.NewApi(client)
	if err != nil {
		logrus.Panic(err)
		return
	}

	logrus.Info("Listing Projects")
	projects, err := crmApi.ListProjects(organization)
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
		err = crmApi.CreateProject(projectId, organization)
		if err != nil {
			logrus.Panic(err)
			return
		}
	}

	err = smApi.EnableApi("project:"+projectId, servicemanagement.DriveApi)
	if err != nil {
		logrus.Panic(err)
		return
	}

	var serviceAccountRequests sync.WaitGroup
	var running int
	for i := 1; i <= 100; i++ {
		serviceAccountRequests.Add(1)
		running++

		go func(i int) {
			defer serviceAccountRequests.Done()

			accountId := fmt.Sprintf("service-account-%d", i)

		createServiceAccount:
			logrus.Infof("Creating Service Account: %s", accountId)
			serviceAccount, err := iamApi.CreateServiceAccount(projectId, accountId, "")
			if err != nil {
				logrus.Error(err)
				goto createServiceAccount
			}

		createApiKey:
			logrus.Infof("Creating Key for Account: %s", accountId)
			serviceAccountKey, err := iamApi.CreateServiceAccountKey(serviceAccount)
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
		}(i)

		if running > App.Flags.Concurrency {
			serviceAccountRequests.Wait()
			running = 0
		}
	}

	serviceAccountRequests.Wait()

	App.AppConfig.Projects = append(App.AppConfig.Projects, projectId)

	logrus.Infof("Done :3")

	boolResponse := false

	confirm := &survey.Confirm{
		Message: "Do you want to assign them to the Service Account Group?",
		Default: true,
	}

	err = survey.AskOne(confirm, &boolResponse, nil)
	if err != nil {
		logrus.Panic(err)
		return
	}

	if !boolResponse {
		return
	}

	logrus.Infof("Waiting 5 seconds before querying google for the accounts")
	time.Sleep(5 * time.Second)

	serviceaccount.CmdAssignServiceAccount(c)
}
