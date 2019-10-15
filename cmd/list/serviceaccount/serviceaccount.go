package teamdrive

import (
	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/iam"
	. "github.com/fionera/TeamDriveManager/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "serviceaccount",
		Usage:  "List all Service Accounts",
		Action: CmdListServiceAccount,
		Flags:  []cli.Flag{},
	}
}

func CmdListServiceAccount(c *cli.Context) {
	projectId := c.Args().First()

	if projectId == "" {
		logrus.Error("Please provide a Project ID")
		return
	}

	logrus.Infof("Using projectId `%s`", projectId)

	client, err := api.CreateClient(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate, []string{iam.CloudPlatformScope})
	if err != nil {
		logrus.Error(err)
		return
	}

	iamApi, err := iam.NewApi(client)
	if err != nil {
		logrus.Error(err)
		return
	}

	serviceAccounts, err := iamApi.ListServiceAccounts(projectId)
	if err != nil {
		logrus.Panic(err)
		return
	}

	for _, serviceAccount := range serviceAccounts {
		logrus.Infof("%s", serviceAccount.Name)
	}

	logrus.Infof("Found %d Service Accounts", len(serviceAccounts))
}
