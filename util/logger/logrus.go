package logger

import (
	"fmt"
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

	log.SetFormatter(&logrus.JSONFormatter{})
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

func EnteringMethod(method string) func() {
	Info(nil, fmt.Sprintf("Entering %s method", method))

	return func() {
		Info(nil, fmt.Sprintf("Exit %s method", method))
	}
}

func WarnMethod(method string, err error) {
	Warn(nil, fmt.Sprintf("Operation encountered on %s method :", method), err.Error())
}

func Error(fields logrus.Fields, args ...any) {
	log.WithFields(fields).Error(args...)
}
