package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "frog",
		Short: "docker images sync",
	}
)

func Execute() error {
	initRootCmd()

	return rootCmd.Execute()
}
