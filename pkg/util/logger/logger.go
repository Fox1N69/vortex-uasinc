package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var e *logrus.Entry

type writeHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

type Logger struct {
	*logrus.Entry
}

func (hook *writeHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}

	return err
}

func (hook *writeHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func Init(mode string) {
	l := logrus.New()
	l.SetReportCaller(true)

	switch mode {
	case "dev":
		l.Formatter = &logrus.JSONFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s %d", filename, f.Line)
			},
		}
	case "prod":
		l.Formatter = &logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
			CallerPrettyfier: func(f *runtime.Frame) (function, file string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		}
	default:
		l.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s %d", filename, f.Line)
			},
		}
	}

	err := os.MkdirAll("logs", 0755)
	if err != nil {
		panic(fmt.Errorf("failed to create logs directory: %w", err))
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(fmt.Errorf("failed to open log file: %w", err))
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writeHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	switch mode {
	case "dev":
		l.SetLevel(logrus.TraceLevel)
	case "release":
		l.SetLevel(logrus.InfoLevel)
	default:
		l.SetLevel(logrus.InfoLevel)
	}

	e = logrus.NewEntry(l)
}

func GetLogger() Logger {
	return Logger{e}
}

func (l *Logger) GetLoggerWithFild(k string, v interface{}) Logger {
	return Logger{l.WithField(k, v)}
}
