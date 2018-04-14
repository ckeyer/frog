package cmd

import (
	"syscall"

	"github.com/ckeyer/commons/pid"
	"github.com/ckeyer/frog/daemon"
	"github.com/ckeyer/logrus"
	"github.com/spf13/cobra"
)

func ReloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reload",
		Short: "reload config.",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Debugf("reload config.")
			proc, err := pid.Open(daemon.PidFile)
			if err != nil {
				logrus.Errorf("found pid failed, %s", err)
				return
			}

			err = proc.Signal(syscall.SIGUSR1)
			if err != nil {
				logrus.Errorf("send reload signal failed, %s", err)
				return
			}
		},
	}
	return cmd
}
