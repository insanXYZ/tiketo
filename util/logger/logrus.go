package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func InitLogger() {
	log = logrus.New()
	var level logrus.Level

	envLevel := strings.ToLower(os.Getenv("APP_LOG_LEVEL"))
	switch envLevel {
	case "trace":
		level = logrus.TraceLevel
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warn", "warning":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	case "fatal":
		level = logrus.FatalLevel
	case "panic":
		level = logrus.PanicLevel
	default:
		level = logrus.InfoLevel
	}

	log.SetLevel(level)
}

func Fatal(fields logrus.Fields, args ...any) {
	log.WithFields(fields).Fatal(args...)
}

func Info(fields logrus.Fields, args ...any) {
	log.WithFields(fields).Info(args...)
}

func Warn(fields logrus.Fields, args ...any) {
	log.WithFields(fields).Warn(args...)
}

func Debug(fields logrus.Fields, args ...any) {
	log.WithFields(fields).Debug(args...)
}
