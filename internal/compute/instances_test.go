package compute

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstances(t *testing.T) {
	client, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseUrl("/zones/europe-west3-c/instances?alt=json&filter=&prettyPrint=false"),
		newJsonFileResponder("./testdata/instances.json")) //nolint bodyclose

	running, err := client.LoadInstancesByZone("europe-west3-c", "")
	assert.NoError(t, err)

	assert.Equal(t, 0, running.GetStatus())
	assert.Contains(t, running.GetOutput(), "[OK]")

	httpmock.RegisterResponder("GET",
		withBaseUrl("/zones/us-central1-a/instances?alt=json&filter=&prettyPrint=false"),
		newJsonFileResponder("./testdata/instances-terminated.json")) //nolint bodyclose

	terminated, err := client.LoadInstancesByZone("us-central1-a", "")
	assert.NoError(t, err)

	assert.Equal(t, 2, terminated.GetStatus())
	assert.Contains(t, terminated.GetOutput(), "[CRITICAL]")
}
