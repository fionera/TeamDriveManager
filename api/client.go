package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func CreateClient(serviceAccountFile, impersonate string, scopes []string) (*http.Client, error) {
	path, err := homedir.Expand(os.ExpandEnv(serviceAccountFile))
	if err != nil {
		return nil, err
	}

	loadedCreds, err := ioutil.ReadFile(path)
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
