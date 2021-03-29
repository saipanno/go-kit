// +build darwin

package shellrunner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

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
	if len(username) != 0 {

		output = "does not support specified user on darwin"
		return
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

	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = &buf
	cmd.Stderr = &buf

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

			output = fmt.Sprintf("execute `%s` failed, output is %s, error message is %s",
				strings.Join(CMDs, " "), output, r.err.Error())
		} else {
			exit = 0
		}

		return
	}
}
