package daemon

import (
	"fmt"
	"os/exec"

	"github.com/ckeyer/logrus"
)

var (
	dockerBin = "docker"
)

func dockerLogin(reg, user, password string) error {
	cmd := exec.Command(dockerBin,
		"login", reg,
		"-u", user,
		"-p", password,
	)

	bs, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("%s", bs)
		return err
	}

	return nil
}

func dockerExec(args ...string) error {
	cmd := exec.Command(dockerBin, args...)

	logrus.Debugf("exec: (%v)", cmd.Args)
	bs, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("%s", bs)
		return err
	}
	logrus.Debugf("exec: (%v) successful. %s", cmd.Args, bs)

	return nil
}

func dockerPull(name string) error {
	return dockerExec("pull", name)
}

func dockerTag(src, dst string) error {
	return dockerExec("tag", src, dst)
}

func dockerPush(name string) error {
	return dockerExec("push", name)
}

func dockerDelete(name string) error {
	return dockerExec("rmi", name)
}

func PullTagPushDelete(origin, target, tag string, del bool) error {
	src := fmt.Sprintf("%s:%s", origin, tag)
	dst := fmt.Sprintf("%s:%s", target, tag)

	if err := dockerPull(src); err != nil {
		return err
	}
	defer func() {
		if del {
			dockerDelete(src)
		}
	}()

	if err := dockerTag(src, dst); err != nil {
		return err
	}
	defer func() {
		if del {
			dockerDelete(dst)
		}
	}()

	return dockerPush(dst)
}
