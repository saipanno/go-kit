package shellrunner

import (
	"errors"
)

type runnerResult struct {
	output string
	err    error
}

var ErrExecTimeout = errors.New("command execute timeout")
