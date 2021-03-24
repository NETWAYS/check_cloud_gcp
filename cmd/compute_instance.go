package cmd

import (
	"check_cloud_gcp/internal/compute"
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
)

var (
	Zone         string
	InstanceName string
)

var vmInstanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Checks a single GCP instance",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := compute.NewClient(JsonFile)
		if err != nil {
			check.ExitError(err)
		}

		var instance *compute.Instance

		if Zone != "" && InstanceName != "" {
			instance, err = client.LoadInstanceByName(Zone, InstanceName)
			if err != nil {
				check.ExitError(err)
			}
		} else {
			check.ExitError(fmt.Errorf("please specify zone name and instance name"))
		}

		output := instance.GetOutput()

		check.Exit(instance.GetStatus(), output)
	},
}

func init() {
	vmInstanceCmd.Flags().StringVarP(&Zone, "zone", "z", "", "GCP zone name")
	vmInstanceCmd.Flags().StringVarP(&InstanceName, "name", "n", "", "Look for instance by name")

	computeCmd.AddCommand(vmInstanceCmd)
}
