package list

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewListServiceAccountCommand() cli.Command {
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

	tokenSource, err := api.NewTokenSource(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate)
	if err != nil {
		logrus.Error(err)
		return
	}

	iamApi, err := api.NewIAMService(tokenSource)
	if err != nil {
		logrus.Error(err)
		return
	}

	serviceAccounts, err := api.ListServiceAccounts(iamApi, projectId)
	if err != nil {
		logrus.Panic(err)
		return
	}

	for _, serviceAccount := range serviceAccounts {
		logrus.Infof("%s", serviceAccount.Name)
	}

	logrus.Infof("Found %d Service Accounts", len(serviceAccounts))
}
