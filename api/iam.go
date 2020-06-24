package api

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

func NewIAMService(jwtConfig *jwt.Config) (*iam.Service, error) {
	jwtConfig.Scopes = append(jwtConfig.Scopes, iam.CloudPlatformScope)
	return iam.NewService(context.TODO(), option.WithTokenSource(jwtConfig.TokenSource(context.TODO())))
}

func CreateServiceAccount(iamApi *iam.Service, projectId, accountId, displayName string) (*iam.ServiceAccount, error) {

	if displayName == "" {
		displayName = strings.Replace(accountId, "-", " ", -1)
	}

	serviceAccount, err := iamApi.Projects.ServiceAccounts.Create("projects/"+projectId, &iam.CreateServiceAccountRequest{
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: displayName,
		},
		AccountId: accountId,
	}).Do()
	if err != nil {
		return nil, err
	}

	return serviceAccount, nil
}

func CreateServiceAccountKey(iamApi *iam.Service, serviceAccount *iam.ServiceAccount) (*iam.ServiceAccountKey, error) {
	return iamApi.Projects.ServiceAccounts.Keys.Create(serviceAccount.Name, &iam.CreateServiceAccountKeyRequest{}).Do()
}

func DeleteServiceAccountKey(iamApi *iam.Service, projectId, accountId string) (*iam.Empty, error) {
	return iamApi.Projects.ServiceAccounts.Keys.Delete(fmt.Sprintf("projects/%s/serviceAccounts/%s", projectId, accountId)).Do()
}

func ListServiceAccounts(iamApi *iam.Service, projectId string) ([]*iam.ServiceAccount, error) {
	var serviceAccounts []*iam.ServiceAccount

	err := iamApi.Projects.ServiceAccounts.List("projects/"+projectId).Pages(context.Background(), func(list *iam.ListServiceAccountsResponse) error {
		serviceAccounts = append(serviceAccounts, list.Accounts...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return serviceAccounts, nil
}

func GetServiceAccount(iamApi *iam.Service, projectId, accountId string) (*iam.ServiceAccount, error) {
	return iamApi.Projects.ServiceAccounts.Get(fmt.Sprintf("projects/%s/serviceAccounts/%s", projectId, accountId)).Context(context.Background()).Do()
}

func DeleteServiceAccount(iamApi *iam.Service, projectId, accountId string) (*iam.Empty, error) {
	return iamApi.Projects.ServiceAccounts.Delete(fmt.Sprintf("projects/%s/serviceAccounts/%s", projectId, accountId)).Do()
}
