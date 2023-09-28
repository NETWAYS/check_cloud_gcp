package compute

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/NETWAYS/check_cloud_gcp/internal/compute/common"

	"github.com/NETWAYS/go-check-network/http/mock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

var GcpProjectName = "foobar"

func withBaseUrl(url string) string {
	return "https://compute.googleapis.com/compute/v1/projects/" + GcpProjectName + url
}

func testClientWithMock() (client *Client, cleanup func()) {
	httpmock.Activate()
	cleanup = httpmock.DeactivateAndReset

	checkhttpmock.ActivateRecorder()

	var (
		err        error
		credential = &google.Credentials{ProjectID: GcpProjectName}
		ctx        = context.Background()
		httpClient = &http.Client{}
	)

	if cred := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); cred != "" {
		httpClient, credential, err = common.NewHTTPClient(cred, Scopes...)
		if err != nil {
			panic(err)
		}

		// Set project name, so that mock tests keep working
		GcpProjectName = credential.ProjectID
	}

	service, err := compute.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		panic(err)
	}

	client = &Client{
		Credential: credential,
		Client:     service,
	}

	return
}

func newJsonFileResponder(fileName string) func(request *http.Request) (*http.Response, error) {
	return func(request *http.Request) (*http.Response, error) {
		data, err := os.ReadFile(fileName)
		if err != nil {
			return nil, err
		}

		return httpmock.NewBytesResponse(200, data), nil
	}
}

func TestClient_LoadInstanceByName(t *testing.T) {
	client, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseUrl("/zones/europe-west3-c/instances/instance-1?alt=json&prettyPrint=false"),
		newJsonFileResponder("./testdata/instance_running.json")) //nolint bodyclose

	instance, err := client.LoadInstanceByName("europe-west3-c", "instance-1")
	assert.NoError(t, err)
	assert.Equal(t, "instance-1", instance.Instance.Name)
}

func TestClient_LoadInstancesByZone(t *testing.T) {
	client, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseUrl("/zones/europe-west3-c/instances?alt=json&filter=&prettyPrint=false"),
		newJsonFileResponder("./testdata/instances.json")) //nolint bodyclose

	instances, err := client.LoadInstancesByZone("europe-west3-c", "")
	assert.NoError(t, err)

	assert.Equal(t, 0, instances.GetStatus())
	assert.Contains(t, instances.GetOutput(), "[OK]")
}

func TestClient_LoadAllInstances(t *testing.T) {
	client, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseUrl("/zones?alt=json&filter=name+%3D+%22europe-west3-c%22&prettyPrint=false"),
		newJsonFileResponder("./testdata/zones.json")) //nolint bodyclose

	httpmock.RegisterResponder("GET",
		withBaseUrl("/zones/europe-west3-c/instances?alt=json&filter=&prettyPrint=false"),
		newJsonFileResponder("./testdata/instances.json")) //nolint bodyclose

	instances, err := client.LoadAllInstances("europe-west3-c", "")
	assert.NoError(t, err)

	assert.Equal(t, 0, instances.GetStatus())

	out := instances.GetOutput()
	assert.Contains(t, out, "[OK]")
	assert.NotContains(t, out, "[CRITICAL]")
}
