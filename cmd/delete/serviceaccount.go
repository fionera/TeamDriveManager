package delete

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	iam2 "google.golang.org/api/iam/v1"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewDeleteServiceAccountCommand() cli.Command {
	return cli.Command{
		Name:   "serviceaccount",
		Usage:  "Delete all ServiceAccounts from a Project",
		Action: CmdDeleteServiceAccount,
		Flags:  []cli.Flag{},
	}
}

func CmdDeleteServiceAccount(c *cli.Context) {
	projectId := c.Args().First()

	if projectId == "" {
		logrus.Error("Please supply the ProjectID to use")
		return
	}

	tokenSource, err := api.NewTokenSource(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate)
	if err != nil {
		logrus.Panic(err)
		return
	}

	iamApi, err := api.NewIAMService(tokenSource)
	if err != nil {
		logrus.Panic(err)
		return
	}

	logrus.Info("Listing Service Accounts")
	accounts, err := api.ListServiceAccounts(iamApi, projectId)

	var serviceAccountRequests sync.WaitGroup
	var running int
	for _, account := range accounts {
		serviceAccountRequests.Add(1)
		running++

		go func(account *iam2.ServiceAccount) {
			defer serviceAccountRequests.Done()
		deleteAccount:
			logrus.Infof("Deleting account `%s`", account.Email)
			_, err := api.DeleteServiceAccount(iamApi, projectId, account.Email)
			if err != nil {
				logrus.Error(err)
				goto deleteAccount
			}
		}(account)

		if running > App.Flags.Concurrency {
			serviceAccountRequests.Wait()
			running = 0
		}
	}

	serviceAccountRequests.Wait()

	logrus.Info("Done :3")
}
