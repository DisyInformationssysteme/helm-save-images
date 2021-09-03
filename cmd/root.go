package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Verbose      bool
	ValueFiles   []string
	ChartVersion string
	DryRun       bool
)

var rootCmd = &cobra.Command{
	Use:   "helm save-images [chart]",
	Short: "save-images is a plugin to save all docker images",
	Long:  `save-images is a plugin to save all docker images necessary to create a given helm release to your local file system.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run(args[0], ChartVersion, ValueFiles, DryRun)
	},
}

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&ValueFiles, "values", "f", []string{}, "specify values in a YAML file or a URL (can specify multiple)")
	rootCmd.PersistentFlags().StringVarP(&ChartVersion, "version", "", "", "chart version")
	rootCmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "", false, "print list of images only")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
