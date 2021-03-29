package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// CallerHook ...
type CallerHook struct {
}

const (
	maximumCallerDepth int = 25

	knownLogrusFrames int = 4
)

var (

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth = 1

	// qualified package name, cached at first use
	logrusPackage string

	localPackage string

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

// getCaller retrieves the name of the first non-logrus calling function
func getCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "logrus") {
				logrusPackage = getPackageName(funcName)
				break
			}
		}

		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "go-kit") {
				localPackage = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logrusPackage && pkg != localPackage {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// Fire ...
func (hook *CallerHook) Fire(entry *Entry) error {

	caller := getCaller()
	if caller != nil {
		entry.Data["caller"] = fmt.Sprintf("%s:%d", filepath.Base(caller.File), caller.Line)
	}

	return nil
}

// Levels ...
func (hook *CallerHook) Levels() []logrus.Level {

	return logrus.AllLevels
}
