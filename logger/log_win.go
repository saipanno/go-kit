// +build windows

package logger

import (
	"log/syslog"

	"github.com/sirupsen/logrus"
)

func setSyslogHook(log *logrus.Logger, host string, slLevel syslog.Priority, app string) (err error) {

	err = errors.New("don't support syslog on windows")
	return
}
