package assign

import (
	"os"
	"sync"
	"time"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/api/googleapi"
	gIAM "google.golang.org/api/iam/v1"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
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

	adminApi, err := api.NewAdminService(tokenSource)
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
			serviceAccounts, err := api.ListServiceAccounts(iamApi, projectId)
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
	if exists, err := api.GroupExists(adminApi, serviceAccountGroupAddress); !exists {
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

		_, err := api.CreateGroup(adminApi, "Service Account Group", serviceAccountGroupAddress)
		if err != nil {
			logrus.Error(err)
			goto checkServiceAccountGroup
		}
		logrus.Info("Successfully created Service Account Group")
	}

listServiceAccountGroupMembers:
	members, err := api.ListMembers(adminApi, serviceAccountGroupAddress)
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
				_, err := api.AddMember(adminApi, serviceAccountGroupAddress, address)
				if err != nil {
					if gerr, ok := err.(*googleapi.Error); ok {
						switch gerr.Code {
						case 409:
							logrus.Info("Account already exists. Skipping.")
							time.Sleep(100 * time.Millisecond)
						default:
							logrus.Error("An error occurred when adding an account. Retrying...", err)
							time.Sleep(100 * time.Millisecond)
							goto addMemberToGroup
						}
					} else {
						logrus.Fatal("An unknown error occurred: ", err)
					}
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
				err := api.RemoveMember(adminApi, serviceAccountGroupAddress, address)
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
