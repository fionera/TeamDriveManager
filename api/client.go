package api

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"os"
)

func CreateClient(serviceAccountFile, impersonate string, scopes []string) (*http.Client, error) {
	loadedCreds, err := ioutil.ReadFile(os.ExpandEnv(serviceAccountFile))
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(loadedCreds, scopes...)
	if err != nil {
		return nil, err
	}

	if impersonate != "" {
		conf.Subject = impersonate
	}

	ctx := context.Background()

	return oauth2.NewClient(ctx, conf.TokenSource(ctx)), nil
}
