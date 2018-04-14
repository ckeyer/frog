package daemon

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ckeyer/commons/pid"
	"github.com/ckeyer/frog/config"
	"github.com/ckeyer/logrus"
)

const (
	// PidFile = "/var/run/frog.pid"
	PidFile = "/tmp/frog.pid"
)

var (
	ConfigFilePath = ""

	ErrReload = errors.New("reload")
)

type Daemon struct {
	sync.Mutex
	*config.Config
	lastStart time.Time
	pidFile   string
	chStop    chan struct{}
	chReload  chan struct{}
}

func New(cfg *config.Config) *Daemon {
	dockerBin = cfg.DockerBin
	return &Daemon{
		Config:    cfg,
		lastStart: time.Now(),
		pidFile:   PidFile,
		chStop:    make(chan struct{}),
	}
}

func (d *Daemon) Run() error {
	go d.waitStop()

	if err := d.initPidFile(); err != nil {
		return err
	}

	if err := d.loginRegistries(d.Registries...); err != nil {
		return err
	}

	wait := time.Second * 1
	isReload := false
	for {
		if isReload {
			isReload = false
			wait = time.Second * 1
			logrus.Info("to reload.")
		}
		logrus.Debugf("wait %s for next time.", wait)
		d.chReload = make(chan struct{})
		select {
		case <-d.chStop:
			aft := time.Second
			logrus.Debugf("stop Run() after %s", aft)
			// wait for delete pid file.
			time.Sleep(aft)
			return nil
		case <-time.Tick(wait):
			d.lastStart = time.Now()
			if err := d.doTasks(d.Config.Tasks...); err != nil {
				if err == ErrReload {
					isReload = true
				}
			}
		case <-d.chReload:
			isReload = true
			continue
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

func (d *Daemon) loginRegistries(regs ...config.Registry) error {
	for _, reg := range regs {
		err := dockerLogin(reg.Server, reg.Username, reg.Password)
		if err != nil {
			logrus.Errorf("login %s use %s failed. %s", reg.Server, reg.Username, err)
			return err
		}
		logrus.Infof("login %s use %s successful.", reg.Server, reg.Username)
	}
	return nil
}

func (d *Daemon) doTasks(tasks ...config.Task) error {
	for _, task := range tasks {
		select {
		case <-d.chReload:
			return ErrReload
		case <-d.chStop:
			return nil
		default:
			if err := d.doOneTask(task); err != nil {
				logrus.Errorf("do %+v failed, %s", task, err)
				return err
			}
		}
	}
	return nil
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
		case <-d.chReload:
			return ErrReload
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
			time.Sleep(time.Second * 6)
			// log.err = PullTagPushDelete(task.Origin, task.Target, tag, d.DeleteEveryTime)
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

// Reload reload
func (d *Daemon) Reload() error {
	cfg, err := config.OpenConfigFile(ConfigFilePath)
	if err != nil {
		logrus.Error("reload config file failed, %s", err)
		return err
	}

	if err := d.loginRegistries(cfg.Registries...); err != nil {
		return nil
	}

	d.Lock()
	defer d.Unlock()
	d.Config = cfg
	close(d.chReload)

	return nil
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
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
	for {
		switch s := <-sigChan; s {
		case syscall.SIGUSR1, syscall.SIGUSR2:
			logrus.Debug("receive reload signal.")
			d.Reload()
		case syscall.SIGINT, syscall.SIGTERM:
			d.Stop()
		}
	}
}
