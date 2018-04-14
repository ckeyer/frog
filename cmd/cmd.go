package cmd

func Execute() error {
	rootCmd := RootCmd()

	rootCmd.AddCommand(ReloadCmd())

	return rootCmd.Execute()
}
