package assign

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/fionera/TeamDriveManager/api"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewAssignGroupCmd() cli.Command {
	return cli.Command{
		Name:   "group",
		Usage:  "Assign all users from the config to the corresponding TeamDrives by using Groups",
		Action: CmdAssignGroup,
		Flags:  []cli.Flag{},
	}
}

var useDomainWide = false

func CmdAssignGroup(c *cli.Context) {
	if App.AppConfig.TeamDriveConfig.NamePrefix == "" && !AskBool("Using no filter is dangerous! Do you really want to continue?", false) {
		return
	}

	logrus.Infof("Using filter `%s`", App.AppConfig.TeamDriveConfig.NamePrefix)
	tokenSource, err := api.NewTokenSource(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate)
	if err != nil {
		logrus.Error(err)
		return
	}

	driveApi, err := api.NewDriveService(tokenSource)
	if err != nil {
		logrus.Error(err)
		return
	}

	adminApi, err := api.NewAdminService(tokenSource)
	if err != nil {
		logrus.Error(err)
		return
	}

	listFunc := api.ListTeamDrives
	if useDomainWide = AskBool("Do you want to list ALL TeamDrives using DomainAdminAccess?", false); useDomainWide {
		listFunc = api.ListAllTeamDrives
	}

	teamDrives, err := listFunc(driveApi)
	if err != nil {
		logrus.Error(err)
		return
	}

	var worker sync.WaitGroup
	var running int
	for _, teamDrive := range teamDrives {
		if !strings.HasPrefix(teamDrive.Name, App.AppConfig.TeamDriveConfig.NamePrefix) {
			continue
		}

		worker.Add(1)
		running++
		go CheckTeamDrive(teamDrive, driveApi, adminApi, &worker)

		if running > App.Flags.Concurrency {
			running = 0
			worker.Wait()
		}
	}
	worker.Wait()
}

func CheckTeamDrive(teamDrive *drive.Drive, driveApi *drive.Service, adminApi *admin.Service, wait *sync.WaitGroup) {
	defer wait.Done()

	logrus.Infof("Checking TeamDrive `%s`", teamDrive.Name)
	users := GatherUsersForTeamDrive(teamDrive.Name)
	groupAssignment := GetGroupAssignmentForTeamDrive(teamDrive.Name)

	for role := range users {
		if groupAssignment.GroupAddresses[role] == "" {
			// Generate the missing addresses for the groups
			groupAddress := GenerateGroupAddressForTeamDrive(teamDrive.Id, role)

		checkGroupExists:
			logrus.Debugf("Looking for Group `%s`", groupAddress)
			exists, err := api.GroupExists(adminApi, groupAddress+"@"+App.AppConfig.Domain)
			if err != nil {
				logrus.Error(err)
				goto checkGroupExists
			}

			if !exists { //&& AskBool(fmt.Sprintf("Cant find group for Teamdrive `%s` with Role `%s`. Should it be created?", teamDrive.Name, role), true) {
				_, err := api.CreateGroup(adminApi, fmt.Sprintf("%s | %s", strings.ToUpper(string(role[0]))+role[1:], teamDrive.Name), groupAddress+"@"+App.AppConfig.Domain)
				if err != nil {
					logrus.Error(err)
					goto checkGroupExists
				}
			}

			groupAssignment.GroupAddresses[role] = groupAddress
		}
	}

	//Check if all users are in the correct groups and if the groups are even there
	for groupRole, groupAddress := range groupAssignment.GroupAddresses {
	listMembers:
		logrus.Debugf("Listing Members for Group `%s` with role `%s`", groupAddress, groupRole)
		groupMembers, err := api.ListMembers(adminApi, groupAddress+"@"+App.AppConfig.Domain)
		if err != nil {
			logrus.Error(err)
			goto listMembers
		}

		// Make a usable list from all members
		var members []string
		for _, member := range groupMembers {
			members = append(members, strings.ToLower(member.Email)) //ToLower since google doesnt respect RFCs
		}

		var userList []string
		for userRole, userAddresses := range users {
			// Skip if not correct role to use
			if userRole == groupRole {
				for _, userAddress := range userAddresses {
					userList = append(userList, strings.ToLower(userAddress)) //ToLower since google doesnt respect RFCs
				}
				break
			}
		}

		// Check if any member is going to be added
		// Go over every user that should have access
		for _, userAddress := range userList {
			// When not in members list -> add him
			if !stringInArray(userAddress, members) {
			addUser:
				logrus.Debugf("Adding `%s` to Group `%s`", userAddress, groupAddress+"@"+App.AppConfig.Domain)
				_, err := api.AddMember(adminApi, groupAddress+"@"+App.AppConfig.Domain, userAddress)
				if err != nil { //Skip existing err
					logrus.Error(err)
					if gerr, ok := err.(*googleapi.Error); ok {
						if gerr.Code != 409 && gerr.Code != 404 {
							goto addUser
						}
					} else if !ok {
						logrus.Panic(err)
					}
				}
			}
		}

		// Check if any member is going to be removed
		// Go over every member that currently has access
		for _, member := range members {
			// When not in user list -> add him
			if !stringInArray(member, userList) {
			removeUser:
				logrus.Debugf("Removing `%s` from Group `%s`", member, groupAddress+"@"+App.AppConfig.Domain)
				err := api.RemoveMember(adminApi, groupAddress+"@"+App.AppConfig.Domain, member)
				if err != nil {
					logrus.Error(err)
					goto removeUser
				}
			}
		}
	}

	listFunc := api.ListPermissions
	if useDomainWide {
		listFunc = api.ListPermissionsAdmin
	}

	deleteFunc := api.DeletePermission
	if useDomainWide {
		deleteFunc = api.DeletePermissionAdmin
	}

	createFunc := api.CreatePermission
	if useDomainWide {
		createFunc = api.CreatePermissionAdmin
	}

listPermissions:
	logrus.Infof("Listing Permissions for TeamDrive `%s`", teamDrive.Name)
	permissions, err := listFunc(driveApi, teamDrive.Id)
	if err != nil {
		logrus.Error(err)
		time.Sleep(500 * time.Millisecond)
		goto listPermissions
	}

	var toAdd = map[string]string{}
	var toRemove []string

	// Check if group is added
	for role, address := range groupAssignment.GroupAddresses {

		// Try to find the current group with correct role
		found := false
		for _, permission := range permissions {
			if permission.Role == role && permission.EmailAddress == address+"@"+App.AppConfig.Domain {
				found = true
				break
			}
		}

		// When not found -> add it
		if !found {
			toAdd[role] = address + "@" + App.AppConfig.Domain
		}
	}

	// Check permissions
	for _, permission := range permissions {
		// Get group address for this role
		groupForRole := groupAssignment.GroupAddresses[permission.Role]

		// If mail is correct -> continue
		if permission.EmailAddress == groupForRole+"@"+App.AppConfig.Domain {
			continue
		}

		logrus.Debugf("Scheduling deletion of Permission `%s` for `%s` with role `%s` from TeamDrive `%s`", permission.Id, permission.EmailAddress, permission.Role, teamDrive.Name)
		toRemove = append(toRemove, permission.Id)
	}

	for role, address := range toAdd {
	createPermission:
		logrus.Debugf("Creating Permission for `%s` with role `%s` on TeamDrive `%s`", address, role, teamDrive.Name)
		_, err := createFunc(driveApi, teamDrive.Id, role, address, "group")
		if err != nil {
			logrus.Error(err)
			goto createPermission
		}
	}

	for _, permissionId := range toRemove {
	deletePermission:
		logrus.Debugf("Deleting Permission `%s` from TeamDrive `%s`", permissionId, teamDrive.Name)
		err := deleteFunc(driveApi, teamDrive.Id, permissionId)
		if err != nil {
			logrus.Error(err)
			goto deletePermission
		}
	}
}

// GenerateGroupAddressForTeamDrive returns a sha256 hash
// printed in hex containing the TeamDriveId and the role combined with an underscore
func GenerateGroupAddressForTeamDrive(id, role string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(id+"_"+role)))
}

func RemoveTeamDrivePrefix(name string) string {
	return strings.Replace(name, App.AppConfig.TeamDriveConfig.NamePrefix, "", 1)
}

func GetGroupAssignmentForTeamDrive(name string) *GroupAssignment {
findGroupAssignment:
	name = RemoveTeamDrivePrefix(name)
	for _, groupAssignment := range App.AppConfig.TeamDriveConfig.GroupAssignments {
		if groupAssignment.TeamDriveName == name {
			return &groupAssignment
		}
	}

	App.AppConfig.TeamDriveConfig.GroupAssignments = append(App.AppConfig.TeamDriveConfig.GroupAssignments, GroupAssignment{
		TeamDriveName:  name,
		GroupAddresses: map[string]string{},
	})

	goto findGroupAssignment
}

// GatherUsersForTeamDrive finds the correct BlackList entry,
// iterates over the GlobalUsers and adds every non blacklisted person into the array corresponding to his role
func GatherUsersForTeamDrive(name string) map[string][]string {
	sortedUsers := make(map[string][]string)

	var blackList []string
	for blackListName, addresses := range App.AppConfig.TeamDriveConfig.BlackList {
		if blackListName == App.AppConfig.TeamDriveConfig.NamePrefix+name {
			blackList = addresses
		}
	}

	for role, addresses := range App.AppConfig.TeamDriveConfig.GlobalUsers {
		if blackList != nil {
			for _, address := range addresses {
				if !stringInArray(address, blackList) {
					sortedUsers[role] = append(sortedUsers[role], address)
				}
			}
		} else {
			sortedUsers[role] = GetStringKeysFromMap(addresses)
		}
	}

	return sortedUsers
}

func GetStringKeysFromMap(data map[string]string) []string {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

func stringInArray(key string, array []string) bool {
	for _, v := range array {
		if key == v {
			return true
		}
	}

	return false
}

func AskBool(message string, def bool) (response bool) {
	confirm := &survey.Confirm{
		Message: message,
		Default: def,
	}

	err := survey.AskOne(confirm, &response, nil)
	if err != nil {
		logrus.Panic(err)
	}

	return
}
