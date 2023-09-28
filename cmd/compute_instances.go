package cmd

import (
	"fmt"

	"github.com/NETWAYS/check_cloud_gcp/internal/compute"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
	"github.com/spf13/cobra"
)

var (
	Filter           string
	IgnoreApiWarning bool
)

var vmInstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Checks multiple GCP instances",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := compute.NewClient(JsonFile)
		if err != nil {
			check.ExitError(err)
		}

		var (
			instances     *compute.Instances
			overallStatus int
			output        string
		)

		instances, err = client.LoadAllInstances(Zone, Filter)
		if err != nil {
			check.ExitError(err)
		}

		states := map[string]int{}
		for _, instance := range instances.Instances {
			states[instance.Instance.Status]++
		}

		overallStatus = result.WorstState(overallStatus, instances.GetStatus())

		output += instances.GetOutput()

		summary := fmt.Sprintf("%d Instances found", len(instances.Instances))

		if len(instances.Instances) <= 0 {
			overallStatus = check.Unknown
		}

		for state, count := range states {
			summary += fmt.Sprintf(" - %d %s", count, state)
		}

		if Zone == "" && !IgnoreApiWarning {
			output = "\nWarning: Please filter for zones, e.g. at least by using a wildcard like: --zone \"europe-*\"\n" +
				"Otherwise, the plugin will query *every* zone worldwide for instances!\n" + output
		}

		check.ExitRaw(overallStatus, summary+"\n"+output)
	},
}

func init() {
	f := vmInstancesCmd.Flags()
	f.StringVarP(&Zone, "zone", "z", "",
		`GCP Zone name, can include wildcards (e.g. "europe-*")`)
	f.StringVarP(&Filter, "filter", "f", "",
		`Filter expression that filters resources e.g. '(cpuPlatform = "Intel Broadwell") AND (name != "instance1")'`)
	f.BoolVar(&IgnoreApiWarning, "ignore-api-warning", false,
		"Disables warning when querying without a zone filter (please do not do that)")

	computeCmd.AddCommand(vmInstancesCmd)
}
