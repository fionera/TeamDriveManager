package operations

import (
	"google.golang.org/api/iam/v1"
)

var iamService *iam.Service
var serviceAccountsService *iam.ProjectsServiceAccountsService

const (
	scopePrefix = "https://www.googleapis.com/auth/"
)

func init() {
	iamService = &iam.Service{
	}

	serviceAccountsService = iam.NewProjectsServiceAccountsService(iamService)
}

func DeleteServiceAccount(key string) error {
	_, err := serviceAccountsService.Delete(key).Do()
	return err
}

func CreateServiceAccount(projectId string, accountId string) (*iam.ServiceAccount, error) {
	return serviceAccountsService.Create("projects/"+projectId, &iam.CreateServiceAccountRequest{
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: kebabToCamelCase(accountId),
		},
		AccountId: accountId,
	}).Do()
}

func CreateServiceAccountKey(serviceAccount *iam.ServiceAccount) (*iam.ServiceAccountKey, error) {
	accountsKeysService := iam.NewProjectsServiceAccountsKeysService(iamService)

	return accountsKeysService.Create(serviceAccount.Name, &iam.CreateServiceAccountKeyRequest{}).Do()
}