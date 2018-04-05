package daemon

import (
	"strings"
	"time"
)

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
