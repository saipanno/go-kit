package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log/syslog"
	"os"
	"path/filepath"
	"strings"

	"github.com/saipanno/go-kit/utils"
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	log             = logrus.New()
	defaultFieldMap = logrus.FieldMap{
		logrus.FieldKeyTime:  "date",
		logrus.FieldKeyLevel: "level",
		logrus.FieldKeyMsg:   "message",
	}

	defaultTimeFormat = "2006-01-02 15:04:05"
)

type (
	// Entry ...
	Entry = logrus.Entry

	// Fields ...
	Fields = logrus.Fields
)

// SetConfig ...
func SetConfig(config *Config) (err error) {

	var wrs []io.Writer
	var slLevel syslog.Priority
	hooks := make(logrus.LevelHooks)

	if len(config.TimeFormat) > 0 {
		defaultTimeFormat = config.TimeFormat
	}

	if len(config.Level) != 0 {

		switch strings.ToUpper(config.Level) {
		case "DEBUG":
			log.SetLevel(logrus.DebugLevel)
			slLevel = syslog.LOG_DEBUG
		case "INFO":
			log.SetLevel(logrus.InfoLevel)
			slLevel = syslog.LOG_INFO
		case "WARN":
			log.SetLevel(logrus.WarnLevel)
			slLevel = syslog.LOG_WARNING
		case "ERROR":
			log.SetLevel(logrus.ErrorLevel)
			slLevel = syslog.LOG_ERR
		default:
			log.SetLevel(logrus.WarnLevel)
			slLevel = syslog.LOG_WARNING
		}
	} else {

		log.SetLevel(logrus.WarnLevel)
		slLevel = syslog.LOG_WARNING
	}

	if strings.ToUpper(config.Encoder) == "JSON" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: defaultTimeFormat,
			FieldMap:        defaultFieldMap,
		})
	} else {
		log.SetFormatter(&defaultFormatter{
			logrus.TextFormatter{
				TimestampFormat: defaultTimeFormat,
				FieldMap:        defaultFieldMap,
			}})
	}

	if len(config.Output) == 0 {
		config.Output = append(config.Output, "stdout")
	}

	if config.ReportCaller {
		hooks.Add(&CallerHook{})
	}

	for _, output := range config.Output {

		if strings.HasPrefix(output, "syslog://") {

			if config.App == "" {
				config.App = "go-kit"
			}

			var hook logrus.Hook
			hook, err = setSyslogHook(log, strings.Replace(output, "syslog://", "", 1), slLevel, config.App)
			if err != nil {
				log.Errorf("set syslog logger failed, message is %s", err.Error())
				return
			}

			hooks.Add(hook)
		} else if strings.ToUpper(output) == "STDOUT" {

			wrs = append(wrs, os.Stdout)
		} else if strings.ToUpper(output) == "STDERR" {

			wrs = append(wrs, os.Stderr)
		} else {

			ljack := &lumberjack.Logger{
				Filename: output,
			}

			err = os.MkdirAll(filepath.Dir(output), 0755)
			if err != nil {
				err = fmt.Errorf("create logfile dir(%s) failed, message is %s",
					filepath.Dir(output), err.Error())
				return
			}

			if config.Logrotate != nil {
				if config.Logrotate.MaxSize != 0 {
					ljack.MaxSize = config.Logrotate.MaxSize
				}
				if config.Logrotate.MaxBackups != 0 {
					ljack.MaxBackups = config.Logrotate.MaxBackups
				}
				if config.Logrotate.MaxAge != 0 {
					ljack.MaxAge = config.Logrotate.MaxAge
				}
				if !config.Logrotate.Compress {
					ljack.Compress = config.Logrotate.Compress
				}
			}

			wrs = append(wrs, ljack)
		}
	}

	logMultiWriter := io.MultiWriter(wrs...)
	log.SetOutput(logMultiWriter)

	log.ReplaceHooks(hooks)
	return
}

// IsDebug ...
func IsDebug() bool {

	return log.GetLevel() >= logrus.DebugLevel
}

// Debug ...
func Debug(message string) {

	log.Debug(message)
}

// Debugf ...
func Debugf(message string, args ...interface{}) {

	log.Debugf(message, args...)
}

// Info ...
func Info(message string) {

	log.Info(message)
}

// Infof ...
func Infof(message string, args ...interface{}) {

	log.Infof(message, args...)
}

// Warn ...
func Warn(message string) {

	log.Warn(message)
}

// Warnf ...
func Warnf(message string, args ...interface{}) {

	log.Warnf(message, args...)
}

// Error ...
func Error(message string) {

	log.Error(message)
}

// Errorf ...
func Errorf(message string, args ...interface{}) {

	log.Errorf(message, args...)
}

// Panic ...
func Panic(message string) {

	log.Panic(message)
}

// Panicf ...
func Panicf(message string, args ...interface{}) {

	log.Panicf(message, args...)
}

// LogfByLevel ...
// 根据LogLevel 自动选择 Infof 或 Debugf
func LogfByLevel(message string, args ...interface{}) {

	al := len(args)

	if al <= 1 {

		if al == 1 {
			message = message + ", raw data is " + utils.PrettyPrint(args[0])
		}

		if IsDebug() {
			Debug(message)
			return
		}

		Info(message)
		return
	}

	args = args[:al-1]
	t := fmt.Sprintf(message, args...)
	if IsDebug() {

		t = t + ", raw data is " + utils.PrettyPrint(args[al-1])
		Debug(t)
	} else {
		Infof(message, args...)
	}
}

// WithField ...
func WithField(key string, value interface{}) *Entry {

	return log.WithField(key, value)
}

// WithFields ...
func WithFields(fields Fields) *Entry {

	return log.WithFields(fields)
}

// Writer ...
func Writer() *io.PipeWriter {

	return log.Writer()
}

// SetStringConfig ...
func SetStringConfig(s string) (err error) {

	var config Config

	err = json.Unmarshal([]byte(s), &config)
	if err != nil {
		Errorf("marshal config failed, message is %s", err.Error())
		return
	}

	err = SetConfig(&config)
	return
}
