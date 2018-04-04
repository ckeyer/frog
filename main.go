package main

import (
	"github.com/ckeyer/logrus"
	"github.com/spf13/cobra"
)

var (
	dockerBin string
	logFile   string
)

func rootCmd() *cobra.Command {
	var (
		cfgFilename string
		debug       bool
	)

	cmd := &cobra.Command{
		Use:   "frog",
		Short: "docker images sync",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Debugf("starting...")
			cfg, err := OpenConfigFile(cfgFilename)
			if err != nil {
				logrus.Errorf("open config file failed, %s", err)
				return
			}

			logrus.Debugf("config: %+v", cfg)

			if err := Run(cfg); err != nil {
				logrus.Error(err)
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "debug level")
	cmd.PersistentFlags().StringVarP(&cfgFilename, "config-file", "f", "", "config filepath.")
	cmd.PersistentFlags().StringVar(&dockerBin, "docker-bin-path", "docker", "docker binary filepath.")
	cmd.PersistentFlags().StringVar(&logFile, "log-file", "/tmp/image_sync", "log filepath.")

	return cmd
}

func main() {
	rootCmd().Execute()
}
