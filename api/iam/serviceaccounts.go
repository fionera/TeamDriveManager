package iam

import (
	"google.golang.org/api/iam/v1"
	"strings"
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
