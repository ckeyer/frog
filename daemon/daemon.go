package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ckeyer/commons/pid"
	"github.com/ckeyer/frog/config"
	"github.com/ckeyer/logrus"
)

type Daemon struct {
	*config.Config
	lastStart time.Time
	pidFile   string
	chStop    chan struct{}
}

func New(cfg *config.Config) *Daemon {
	dockerBin = cfg.DockerBin
	return &Daemon{
		Config:    cfg,
		lastStart: time.Now(),
		pidFile:   "/var/run/frog.pid",
		chStop:    make(chan struct{}),
	}
}

func (d *Daemon) Run() error {
	go d.waitStop()

	if err := d.initPidFile(); err != nil {
		return err
	}

	for _, reg := range d.Registries {
		err := dockerLogin(reg.Server, reg.Username, reg.Password)
		if err != nil {
			return err
		}
		logrus.Infof("login %s use %s successful.", reg.Server, reg.Username)
	}

	wait := time.Second * 1
	for {
		logrus.Debugf("wait %s for next time.", wait)
		select {
		case <-time.Tick(wait):
			d.lastStart = time.Now()
			d.doTasks()
		case <-d.chStop:
			aft := time.Second
			logrus.Debugf("stop Run() after %s", aft)
			// wait for delete pid file.
			time.Sleep(aft)
			return nil
		}

		last := d.lastStart
		wait = time.Duration(d.Period)
		if last.Add(wait).Before(time.Now()) {
			wait = time.Second * 5
		} else {
			wait = last.Add(wait).Sub(time.Now())
		}
	}

	return nil
}

func (d *Daemon) doTasks() {
	for _, task := range d.Config.Tasks {
		select {
		case <-d.chStop:
			return
		default:
			if err := d.doOneTask(task); err != nil {
				logrus.Errorf("do %+v failed, %s", task, err)
			}
		}
	}
}

func (d *Daemon) doOneTask(task config.Task) error {
	file := fmt.Sprintf("%s.%s", d.LogFile, d.lastStart.Format("2006-01-02T15-04-05"))
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logrus.Error("save log failed, %s", err)
	}
	defer f.Close()

	for _, tag := range task.Tags {
		select {
		case <-d.chStop:
			return fmt.Errorf("user stoped.")
		default:
			log := &taskLog{
				origin:  task.Origin,
				target:  task.Target,
				tag:     tag,
				startAt: time.Now(),
			}
			logrus.Debugf("start %s:%s -> %s:%s", task.Origin, tag, task.Target, tag)
			log.err = PullTagPushDelete(task.Origin, task.Target, tag, d.DeleteEveryTime)
			f.WriteString(log.String())
		}
	}
	return nil
}

func (d *Daemon) Stop() {
	select {
	case <-d.chStop:
	default:
		close(d.chStop)
		logrus.Info("frog stopping...")
	}
}

// initPidFile
func (d *Daemon) initPidFile() error {
	return pid.Generate(d.pidFile, d.chStop)
}

// waitStop
func (d *Daemon) waitStop() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	switch s := <-sigChan; s {
	case syscall.SIGINT, syscall.SIGTERM:
		d.Stop()
	}
}
