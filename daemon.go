package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ckeyer/logrus"
)

func Run(cfg *Config) error {
	for _, reg := range cfg.Registries {
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
			cfg.Global.lastStart = time.Now()
			doTasks(cfg.Tasks, cfg.Global)
		}

		last := cfg.Global.lastStart
		wait = time.Duration(cfg.Global.Period)
		if last.Add(wait).Before(time.Now()) {
			wait = time.Second * 5
		} else {
			wait = last.Add(wait).Sub(time.Now())
		}
	}

	return nil
}

func doTasks(tasks []Task, opt Global) {
	for _, task := range tasks {
		if err := doOneTask(task, opt); err != nil {
			logrus.Errorf("do %+v failed, %s", task, err)
		}
	}
}

func doOneTask(task Task, opt Global) error {
	file := fmt.Sprintf("%s.%s", logFile, opt.lastStart.Format("2006-01-02T15-04-05"))
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logrus.Error("save log failed, %s", err)
	}
	defer f.Close()

	for _, tag := range task.Tags {
		log := &taskLog{
			origin:  task.Origin,
			target:  task.Target,
			tag:     tag,
			startAt: time.Now(),
		}
		log.err = PullTagPushDelete(task.Origin, task.Target, tag, opt.DeleteEveryTime)
		f.WriteString(log.String())
	}
	return nil
}

type taskLog struct {
	origin  string
	target  string
	tag     string
	startAt time.Time
	err     error
}

func (t taskLog) String() string {
	var (
		status = "[SUC]"
		errstr = ""
	)
	if t.err != nil {
		status = "[ERR]"
		errstr = t.err.Error()
	}
	strs := []string{
		status,
		t.origin,
		t.tag,
		t.target,
		time.Now().Sub(t.startAt).String(),
		errstr,
		"\n",
	}

	return strings.Join(strs, " ")
}
