package common

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

func NewHTTPClient(filename string, scopes ...string) (client *http.Client, credential *google.Credentials, err error) {
	credData, err := ioutil.ReadFile(filename)
	if err != nil {
		err = fmt.Errorf("could not load credenial file '%s': %w", filename, err)
		return
	}

	ctx := context.Background()

	credential, err = google.CredentialsFromJSON(ctx, credData, scopes...)
	if err != nil {
		err = fmt.Errorf("could not load credentials: %w", err)
		return
	}

	client = oauth2.NewClient(ctx, credential.TokenSource)

	return
}
