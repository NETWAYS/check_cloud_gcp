package compute

import (
	"context"
	"fmt"

	"github.com/NETWAYS/check_cloud_gcp/internal/compute/common"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type Client struct {
	Client     *compute.Service
	Credential *google.Credentials
}

var Scopes = []string{
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/compute",
	"https://www.googleapis.com/auth/compute.readonly",
	"https://www.googleapis.com/auth/devstorage.full_control",
	"https://www.googleapis.com/auth/devstorage.read_only",
	"https://www.googleapis.com/auth/devstorage.read_write",
}

func NewClient(jsonCred string) (c *Client, err error) {
	ctx := context.Background()

	c = &Client{}
	httpClient, cred, err := common.NewHTTPClient(jsonCred, Scopes...)

	if err != nil {
		err = fmt.Errorf("could not create http client: %w", err)
		return
	}

	c.Credential = cred
	c.Client, err = compute.NewService(ctx, option.WithHTTPClient(httpClient))

	if err != nil {
		err = fmt.Errorf("could not create client: %w", err)
		return
	}

	return
}

func (c *Client) LoadInstanceByName(zone, name string) (instance *Instance, err error) {
	local, err := c.Client.Instances.Get(c.Credential.ProjectID, zone, name).Do()
	if err != nil {
		err = fmt.Errorf("cloud not load instance '%s' in zone '%s': %w", name, zone, err)
		return
	}

	instance = &Instance{Instance: local}

	return
}

func (c *Client) LoadInstancesByZone(zone, filter string) (instances *Instances, err error) {
	instances = &Instances{}

	instanceList, err := c.Client.Instances.List(c.Credential.ProjectID, zone).Filter(filter).Do()
	if err != nil {
		err = fmt.Errorf("cloud not load instances in zone '%s': %w", zone, err)
		return
	}

	for _, instance := range instanceList.Items {
		instances.Instances = append(instances.Instances, &Instance{instance})
	}

	return
}

func (c *Client) LoadAllInstances(zone, filter string) (instances *Instances, err error) {
	instances = &Instances{}

	call := c.Client.Zones.List(c.Credential.ProjectID)

	if zone != "" {
		call.Filter(fmt.Sprintf(`name = "%s"`, zone))
	}

	zoneList, err := call.Do()
	if err != nil {
		err = fmt.Errorf("could not load zones: %w", err)
		return
	}

	for _, zone := range zoneList.Items {
		var zoneInstances *Instances

		zoneInstances, err = c.LoadInstancesByZone(zone.Name, filter)

		if err != nil {
			return
		}

		for _, i := range zoneInstances.Instances {
			instances.Instances = append(instances.Instances, i)
		}
	}

	return
}
