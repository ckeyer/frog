package cmd

import (
	"github.com/ckeyer/commons/version"
	"github.com/spf13/cobra"
)

func Execute() error {
	rootCmd := RootCmd()

	rootCmd.AddCommand(ReloadCmd())
	rootCmd.AddCommand(printVersionCmd())

	return rootCmd.Execute()
}

func printVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			version.Print(nil)
		},
	}
}
