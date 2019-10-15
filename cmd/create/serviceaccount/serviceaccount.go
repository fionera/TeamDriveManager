package serviceaccount

import (
	"encoding/base64"
	"github.com/Jeffail/gabs"
	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/cloudresourcemanager"
	"github.com/fionera/TeamDriveManager/api/iam"
	"github.com/fionera/TeamDriveManager/api/servicemanagement"
	. "github.com/fionera/TeamDriveManager/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
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

	err = ioutil.WriteFile(App.AppConfig.ServiceAccountFolder+"/"+serviceAccount.ProjectId+"_"+serviceAccount.DisplayName+".json", jsonData, 0755)
	if err != nil {
		logrus.Panic(err)
		return
	}
}
