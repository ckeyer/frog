package cmd

import (
	"github.com/ckeyer/frog/config"
	"github.com/ckeyer/frog/daemon"
	"github.com/ckeyer/logrus"
	"github.com/spf13/cobra"
)

const (
	defaultDockerBin = "docker"
	defaultLogFile   = "/tmp/image_sync"
)

func initRootCmd() {
	var (
		cfgFilename string
		debug       bool
		dockerBin   string
		logFile     string
	)
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
	}
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		logrus.Debugf("starting...")
		cfg, err := config.OpenConfigFile(cfgFilename)
		if err != nil {
			logrus.Errorf("open config file failed, %s", err)
			return
		}
		if dockerBin != "" {
			cfg.DockerBin = dockerBin
		} else if cfg.DockerBin == "" {
			cfg.DockerBin = defaultDockerBin
		}
		if logFile != "" {
			cfg.LogFile = logFile
		} else if cfg.LogFile == "" {
			cfg.LogFile = defaultLogFile
		}

		logrus.Debugf("config.Global: %+v", cfg.Global)
		for _, r := range cfg.Registries {
			logrus.Debugf("config.Registry: %+v", r)
		}
		i := 1
		for _, task := range cfg.Tasks {
			for _, tag := range task.Tags {
				logrus.Debugf("task %3d: %s:%s -> %s:%s", i, task.Origin, tag, task.Target, tag)
				i++
			}
		}

		d := daemon.New(cfg)
		if err := d.Run(); err != nil {
			logrus.Error(err)
		}
	}

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "debug level")
	rootCmd.PersistentFlags().StringVarP(&cfgFilename, "config-file", "f", "", "config filepath.")
	rootCmd.PersistentFlags().StringVar(&dockerBin, "docker-bin-path", "", "docker binary filepath.")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "log filepath.")
}
