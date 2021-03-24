package cmd

import "github.com/spf13/cobra"

var computeCmd = &cobra.Command{
	Use:   "compute",
	Short: "Checks in the Compute Engine context",
	Run:   Help,
}
