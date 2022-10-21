package common

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewHTTPClient(filename string, scopes ...string) (client *http.Client, credential *google.Credentials, err error) {
	credData, err := os.ReadFile(filename)
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
