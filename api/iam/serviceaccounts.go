package iam

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/api/iam/v1"
)

func (a *Api) CreateServiceAccount(projectId, accountId, displayName string) (*iam.ServiceAccount, error) {

	if displayName == "" {
		displayName = strings.Replace(accountId, "-", " ", -1)
	}

	serviceAccount, err := a.iam.Projects.ServiceAccounts.Create("projects/"+projectId, &iam.CreateServiceAccountRequest{
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

func (a *Api) CreateServiceAccountKey(serviceAccount *iam.ServiceAccount) (*iam.ServiceAccountKey, error) {
	serviceAccountKey, err := a.iam.Projects.ServiceAccounts.Keys.Create(serviceAccount.Name, &iam.CreateServiceAccountKeyRequest{}).Do()
	if err != nil {
		return nil, err
	}

	return serviceAccountKey, nil
}

func (a *Api) ListServiceAccounts(projectId string) ([]*iam.ServiceAccount, error) {
	var serviceAccounts []*iam.ServiceAccount

	err := a.iam.Projects.ServiceAccounts.List("projects/"+projectId).Pages(context.Background(), func(list *iam.ListServiceAccountsResponse) error {
		serviceAccounts = append(serviceAccounts, list.Accounts...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return serviceAccounts, nil
}

func (a *Api) DeleteServiceAccount(projectId, accountId string) error {
	_, err := a.iam.Projects.ServiceAccounts.Delete(fmt.Sprintf("projects/%s/serviceAccounts/%s", projectId, accountId)).Do()

	if err != nil {
		return err
	}

	return nil
}
