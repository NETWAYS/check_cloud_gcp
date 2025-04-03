package cmd

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
	"os"
)

var (
	Timeout  = 30
	JsonFile string
)

var rootCmd = &cobra.Command{
	Use:   "check_cloud_gcp",
	Short: "Check plugin to check Google Cloud virtual machines.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		go check.HandleTimeout(Timeout)

		if JsonFile == "" && os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
			err := fmt.Errorf("please specify the service GCP account key file")
			check.ExitError(err)
		} else if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
			JsonFile = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		}
	},
	Run: Help,
}

func Execute(version string) {
	defer check.CatchPanic()

	rootCmd.Version = version

	if err := rootCmd.Execute(); err != nil {
		check.ExitError(err)
	}
}

func Help(cmd *cobra.Command, strings []string) {
	fmt.Println(cmd.Short)
	fmt.Println()

	_ = cmd.Usage()

	os.Exit(3)
}

func init() {
	rootCmd.AddCommand(computeCmd)
	rootCmd.SetHelpFunc(Help)

	p := rootCmd.PersistentFlags()
	p.StringVarP(&JsonFile, "json-file", "j", JsonFile, "GCP service account key file")
}
