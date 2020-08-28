package combinations

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/cmd/assign"
	. "github.com/fionera/TeamDriveManager/config"
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

	smApi, err := api.NewServiceManagementService(tokenSource)
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
		err = api.CreateProject(crmApi, projectId, organization)
		if err != nil {
			logrus.Panic(err)
			return
		}
	}

	err = api.EnableApi(smApi, "project:"+projectId, api.DriveApi)
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
			serviceAccount, err := api.CreateServiceAccount(iamApi, projectId, accountId, "")
			if err != nil {
				if gerr, ok := err.(*googleapi.Error); ok {
					if gerr.Code == 429 {
						logrus.Fatalf("Project %s reached maximum service account amount limit.", projectId)
						return
					} else if gerr.Code == 409 {
						var serr error
						serviceAccount, serr = api.GetServiceAccount(iamApi, projectId, accountId+"@"+projectId+".iam.gserviceaccount.com")
						if serr != nil {
							logrus.Error(err)
							goto createServiceAccount
						}
						goto createApiKey
					} else {
						logrus.Error(err)
						goto createServiceAccount
					}
				} else if !ok {
					logrus.Error(err)
				}
				goto createServiceAccount
			}

		createApiKey:
			logrus.Infof("Creating Key for Account: %s", accountId)
			serviceAccountKey, err := api.CreateServiceAccountKey(iamApi, serviceAccount)
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

	assign.CmdAssignServiceAccount(c)
}
