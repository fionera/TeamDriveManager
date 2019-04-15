package serviceaccount

import (
	"github.com/codegangsta/cli"
	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/admin"
	"github.com/fionera/TeamDriveManager/api/iam"
	. "github.com/fionera/TeamDriveManager/config"
	"github.com/sirupsen/logrus"
	gIAM "google.golang.org/api/iam/v1"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
	"sync"
	"time"
)

func NewAssignServiceAccountsCmd() cli.Command {
	return cli.Command{
		Name:   "serviceaccount",
		Usage:  "Assign all ServiceAccounts from the projects defined in the config to the configured Group",
		Action: CmdAssignServiceAccount,
		Flags:  []cli.Flag{},
	}
}

func CmdAssignServiceAccount(c *cli.Context) {
	projectIds := App.AppConfig.Projects

	client, err := api.CreateClient(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate, []string{iam.CloudPlatformScope, admin.AdminDirectoryGroupScope})
	if err != nil {
		logrus.Error(err)
		return
	}

	iamApi, err := iam.NewApi(client)
	if err != nil {
		logrus.Error(err)
		return
	}

	adminApi, err := admin.NewApi(client)
	if err != nil {
		logrus.Error(err)
		return
	}

	var requests sync.WaitGroup
	var running int
	var mutex sync.Mutex
	var accounts []*gIAM.ServiceAccount

	for _, projectId := range projectIds {
		running++
		requests.Add(1)

		go func(projectId string) {
			defer requests.Done()

			logrus.Infof("Requesting Service Accounts for projectId `%s`", projectId)

		requestServiceAccounts:
			serviceAccounts, err := iamApi.ListServiceAccounts(projectId)
			if err != nil {
				logrus.Error(err)
				goto requestServiceAccounts
			}

			mutex.Lock()
			accounts = append(accounts, serviceAccounts...)
			mutex.Unlock()
		}(projectId)

		if running > App.Flags.Concurrency {
			requests.Wait()
			running = 0
		}
	}
	requests.Wait()

	logrus.Infof("Found %d Accounts", len(accounts))

	serviceAccountGroupAddress := App.AppConfig.ServiceAccountGroup + "@" + App.AppConfig.Domain

	logrus.Info("Checking for Service Account Group")
checkServiceAccountGroup:
	if exists, err := adminApi.GroupExists(serviceAccountGroupAddress); !exists {
		if err != nil {
			logrus.Panic(err)
			return
		}

		logrus.Warnf("Couldnt find Service Account Group with Address `%s`.", serviceAccountGroupAddress)
		cont := false
		prompt := &survey.Confirm{
			Message: "Should it be created?",
			Default: true,
		}
		err = survey.AskOne(prompt, &cont, nil)
		if err != nil {
			logrus.Error(err)
			return
		}

		if !cont {
			logrus.Info("Cancelling.")
			os.Exit(0)
		}

		_, err := adminApi.CreateGroup("Service Account Group", serviceAccountGroupAddress)
		if err != nil {
			logrus.Error(err)
			goto checkServiceAccountGroup
		}
		logrus.Info("Successfully created Service Account Group")
	}

listServiceAccountGroupMembers:
	members, err := adminApi.ListMembers(serviceAccountGroupAddress)
	if err != nil {
		logrus.Error(err)
		logrus.Info("This can happen when the group is new. Retrying in 2 Seconds.")
		time.Sleep(2 * time.Second) // Google is slow
		goto listServiceAccountGroupMembers
	}

	logrus.Info("Checking which accounts to add and which member to remove")

	var toAdd []string
	var toRemove []string
	for _, member := range members {
		found := false
		for _, account := range accounts {
			if member.Email == account.Email {
				found = true
				break
			}
		}

		if !found {
			toRemove = append(toRemove, member.Email)
		}
	}

	for _, account := range accounts {
		found := false
		for _, member := range members {
			if member.Email == account.Email {
				found = true
				break
			}
		}

		if !found {
			toAdd = append(toAdd, account.Email)
		}
	}

	logrus.Infof("Need to add %d Accounts and remove %d Members", len(toAdd), len(toRemove))
	cont := false
	prompt := &survey.Confirm{
		Message: "Do you really want to continue?",
		Default: true,
	}
	err = survey.AskOne(prompt, &cont, nil)
	if err != nil {
		logrus.Error(err)
		return
	}

	if !cont {
		logrus.Info("Cancelling.")
		os.Exit(0)
	}

	if len(toAdd) > 0 {
		logrus.Info("Start adding Accounts")
		running = 0
		for _, address := range toAdd {
			running++
			requests.Add(1)

			go func(address string) {
				defer requests.Done()

			addMemberToGroup:
				logrus.Debugf("Adding %s to Group", address)
				_, err := adminApi.AddMember(serviceAccountGroupAddress, address)
				if err != nil {
					logrus.Error("An error occurred when adding an account. Retrying...", err)
					time.Sleep(100 * time.Millisecond)
					goto addMemberToGroup
				}
			}(address)

			if running > App.Flags.Concurrency {
				requests.Wait()
				running = 0
			}
		}
		requests.Wait()
		logrus.Info("Done adding Accounts")
	}

	if len(toRemove) > 0 {
		logrus.Info("Start removing Members")
		running = 0
		for _, address := range toAdd {
			running++
			requests.Add(1)

			go func(address string) {
				defer requests.Done()

			removeMemeberFromGroup:
				logrus.Debugf("Removing %s from Group", address)
				err := adminApi.RemoveMember(serviceAccountGroupAddress, address)
				if err != nil {
					logrus.Error("An error occurred when removing a member. Retrying...", err)
					time.Sleep(100 * time.Millisecond)
					goto removeMemeberFromGroup
				}
			}(address)

			if running > App.Flags.Concurrency {
				requests.Wait()
				running = 0
			}
		}
		requests.Wait()
		logrus.Info("Done removing Members")
	}

	logrus.Info("Done :3")
}
