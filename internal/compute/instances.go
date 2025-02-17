package compute

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
	"sort"
	"strings"
)

type Instances struct {
	Instances []*Instance
}

func (i *Instances) GetStatus() int {
	var states []int

	for _, instance := range i.Instances {
		states = append(states, instance.GetStatus())
	}

	return result.WorstState(states...)
}

func (i *Instances) GetOutput() (output string) {
	mapped := map[string][]*Instance{}

	// Index instances by zone
	for _, instance := range i.Instances {
		zoneInfo := strings.Split(instance.Instance.Zone, "/")
		zone := zoneInfo[len(zoneInfo)-1]

		mapped[zone] = append(mapped[zone], instance)
	}

	// Get sorted zones
	var zones []string
	for zone := range mapped {
		zones = append(zones, zone)
	}

	sort.Strings(zones)

	// Prepare output
	for _, zone := range zones {
		output += fmt.Sprintf("\n## %s\n\n", zone)

		for _, instance := range mapped[zone] {
			output += fmt.Sprintf("[%s] %s\n", check.StatusText(instance.GetStatus()), instance.GetOutput())
		}
	}

	return
}
