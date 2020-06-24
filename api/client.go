package api

import (
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

func NewTokenSource(serviceAccountFile, impersonate string) (*jwt.Config, error) {
	path, err := homedir.Expand(os.ExpandEnv(serviceAccountFile))
	if err != nil {
		return nil, err
	}

	loadedCreds, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(loadedCreds)
	if err != nil {
		return nil, err
	}

	if impersonate != "" {
		conf.Subject = impersonate
	}

	return conf, nil
}
