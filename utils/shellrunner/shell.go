// +build !darwin,!windows

package shellrunner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// GetUIDGID ...
func GetUIDGID(username string) (uid int, gid int, err error) {
	u, err := user.Lookup(username)
	if err != nil {
		err = fmt.Errorf("username `%s` lookup failed, message is %s", username, err.Error())
		return
	}

	uid, err = strconv.Atoi(u.Uid)
	if err != nil {
		err = fmt.Errorf("username `%s` uid(%s) to int failed, message is %s",
			username, u.Uid, err.Error())
		return
	}

	gid, err = strconv.Atoi(u.Gid)
	if err != nil {
		err = fmt.Errorf("username `%s` gid(%s) to int failed, message is %s",
			username, u.Gid, err.Error())
	}

	return
}

// RunCommandWithEnvUser ...
// 		output = stdout + stderr
// 		exit = -1(timeout)
func RunCommandWithEnvUser(CMDs, env []string, username string, timeout ...int) (output string, exit int) {

	var (
		err    error
		args   []string
		buf    bytes.Buffer
		ctx    context.Context
		cancel context.CancelFunc
	)

	exit = 1 // default exit code, failed

	if len(CMDs) == 0 {
		output = "command cannot be empty"
		return
	}

	// default current user
	if len(username) == 0 {
		var u *user.User
		u, err = user.Current()
		if err != nil {
			output = fmt.Sprintf("get current user failed, message is %s", err.Error())
			return
		}

		username = u.Name
	}

	if len(timeout) > 0 && timeout[0] > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout[0])*time.Second)
		defer cancel()
	} else {
		ctx = context.TODO()
	}

	bin := CMDs[0]
	if len(CMDs) > 1 {
		args = CMDs[1:]
	}

	var uid, gid int
	uid, gid, err = GetUIDGID(username)
	if err != nil {
		output = fmt.Sprintf("user(%s) does not exist, message is %s", username, err.Error())
		return
	}

	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:    true,
		Credential: &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)},
	}

	rr := make(chan runnerResult)
	go func() {
		defer close(rr)

		err = cmd.Run()
		rr <- runnerResult{buf.String(), err}
	}()

	select {
	case <-ctx.Done():
		output = fmt.Sprintf("execute `%s` timeout", strings.Join(CMDs, " "))
		err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if err != nil {
			output = output + fmt.Sprintf(", kill(%d) failed, message is %s", cmd.Process.Pid, err.Error())
		}

		exit = -1
		return

	case r := <-rr:
		output = r.output

		if r.err != nil {
			if exitError, ok := r.err.(*exec.ExitError); ok {
				exit = exitError.Sys().(syscall.WaitStatus).ExitStatus()
			}
		} else {
			exit = 0
		}

		return
	}
}
