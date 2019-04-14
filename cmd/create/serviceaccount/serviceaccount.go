package serviceaccount

import (
	"github.com/codegangsta/cli"
	"github.com/fionera/TeamdriveManager/api"
	"github.com/fionera/TeamdriveManager/api/cloudresourcemanager"
	"github.com/fionera/TeamdriveManager/api/iam"
	"github.com/fionera/TeamdriveManager/api/servicemanagement"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "serviceaccount",
		Usage:  "Create a ServiceAccount",
		Action: CmdCreateProject,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "service-account-file",
			},
			cli.StringFlag{
				Name: "impersonate",
			},
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
	serviceAccountFile := c.String("service-account-file")
	impersonate := c.String("impersonate")
	projectId := c.String("project-id")
	accountId := c.String("account-id")

	client, err := api.CreateClient(serviceAccountFile, impersonate, []string{cloudresourcemanager.CloudPlatformScope, servicemanagement.ServiceManagementScope})
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

	err = ioutil.WriteFile(serviceAccount.ProjectId+"_"+serviceAccount.DisplayName+".json", json, 0755)
	if err != nil {
		logrus.Panic(err)
		return
	}
}
