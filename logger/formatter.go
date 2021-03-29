package logger

import (
	"bytes"
	"strings"

	"github.com/sirupsen/logrus"
)

type defaultFormatter struct {
	logrus.TextFormatter
}

func (f *defaultFormatter) Format(entry *Entry) ([]byte, error) {

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimeFormat
	}

	b.WriteString(entry.Time.Format(timestampFormat))

	b.WriteString(" [")
	b.WriteString(strings.ToUpper(entry.Level.String())[0:1])
	b.WriteString("] ")

	b.WriteString(entry.Message)
	b.WriteString(" ")

	for k, v := range entry.Data {
		if s, ok := v.(string); ok {
			b.WriteString(strings.Join([]string{k, "=", s, " "}, ""))
		}
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}
