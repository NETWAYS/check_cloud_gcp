package compute

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstance(t *testing.T) {
	client, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseUrl("/zones/europe-west3-c/instances/instance-1?alt=json&prettyPrint=false"),
		newJsonFileResponder("./testdata/instance_running.json")) //nolint bodyclose

	instance, err := client.LoadInstanceByName("europe-west3-c", "instance-1")
	assert.NoError(t, err)

	assert.Equal(t, 0, instance.GetStatus())

	output := instance.GetOutput()
	assert.Contains(t, output, "powerstate=RUNNING")
	assert.Contains(t, output, "size=e2-micro")
}

func TestInstance_terminated(t *testing.T)  {
	client, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseUrl("/zones/europe-west3-c/instances/instance-1?alt=json&prettyPrint=false"),
		newJsonFileResponder("./testdata/instance_terminated.json")) //nolint bodyclose

	instance, err := client.LoadInstanceByName("europe-west3-c", "instance-1")
	assert.NoError(t, err)

	assert.Equal(t, 2, instance.GetStatus())

	output := instance.GetOutput()
	assert.Contains(t, output, "powerstate=TERMINATED")
	assert.Contains(t, output, "size=e2-micro")
}
