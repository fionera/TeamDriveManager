package serviceaccount

import (
	"github.com/codegangsta/cli"
	"github.com/fionera/TeamdriveManager/api"
	"github.com/fionera/TeamdriveManager/api/cloudresourcemanager"
	"github.com/fionera/TeamdriveManager/api/iam"
	"github.com/fionera/TeamdriveManager/api/servicemanagement"
	. "github.com/fionera/TeamdriveManager/config"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "serviceaccount",
		Usage:  "Create a ServiceAccount",
		Action: CmdCreateProject,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "project-id",
			},
			cli.StringFlag{
				Name: "account-id",
			},
		},
	}
}

func CmdCreateProject(c *cli.Context) {
	projectId := c.String("project-id")
	accountId := c.String("account-id")

	if projectId == "" {
		logrus.Error("Please supply the ProjectID to use")
		return
	}

	if accountId == "" {
		logrus.Error("Please supply the AccountID to use")
		return
	}

	client, err := api.CreateClient(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate, []string{cloudresourcemanager.CloudPlatformScope, servicemanagement.ServiceManagementScope})
	if err != nil {
		logrus.Panic(err)
		return
	}

	iamApi, err := iam.NewApi(client)
	if err != nil {
		logrus.Panic(err)
		return
	}

	logrus.Infof("Creating Service Account: %s", accountId)
	serviceAccount, err := iamApi.CreateServiceAccount(projectId, accountId, "")
	if err != nil {
		logrus.Panic(err)
		return
	}

	logrus.Infof("Creating Key for Account: %s", accountId)
	serviceAccountKey, err := iamApi.CreateServiceAccountKey(serviceAccount)
	if err != nil {
		logrus.Panic(err)
		return
	}

	json, err := serviceAccountKey.MarshalJSON()
	if err != nil {
		logrus.Panic(err)
		return
	}

	err = os.Mkdir(App.AppConfig.ServiceAccountFolder, 0755)
	if err != nil && !os.IsExist(err) {
		logrus.Panicf("Error changing type: %s", err)
		return
	}

	err = ioutil.WriteFile(App.AppConfig.ServiceAccountFolder+"/"+serviceAccount.ProjectId+"_"+serviceAccount.DisplayName+".json", json, 0755)
	if err != nil {
		logrus.Panic(err)
		return
	}
}
