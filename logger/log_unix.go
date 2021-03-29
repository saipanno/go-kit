// +build !windows,!nacl,!plan9

package logger

import (
	"log/syslog"

	"github.com/sirupsen/logrus"
	slhooks "github.com/sirupsen/logrus/hooks/syslog"
)

func setSyslogHook(log *logrus.Logger, host string, slLevel syslog.Priority, app string) (hook logrus.Hook, err error) {

	hook, err = slhooks.NewSyslogHook("udp", host, slLevel, app)
	if err != nil {
		log.Errorf("unable to connect to syslog server, message is %s", err.Error())
		return
	}

	hooks := make(logrus.LevelHooks)
	hooks.Add(hook)

	log.ReplaceHooks(hooks)

	return
}
