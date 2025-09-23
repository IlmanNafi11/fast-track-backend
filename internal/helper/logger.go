package helper

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger(env string) {
	Logger = logrus.New()

	Logger.SetOutput(os.Stdout)

	if env == "production" {
		Logger.SetFormatter(&logrus.JSONFormatter{})
		Logger.SetLevel(logrus.WarnLevel)
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
		Logger.SetLevel(logrus.InfoLevel)
	}
}

func GetLogger() *logrus.Logger {
	if Logger == nil {
		InitLogger("development")
	}
	return Logger
}

func Info(message string, fields ...logrus.Fields) {
	logger := GetLogger()
	if len(fields) > 0 {
		logger.WithFields(fields[0]).Info(message)
	} else {
		logger.Info(message)
	}
}

func Warn(message string, fields ...logrus.Fields) {
	logger := GetLogger()
	if len(fields) > 0 {
		logger.WithFields(fields[0]).Warn(message)
	} else {
		logger.Warn(message)
	}
}

func Error(message string, err error, fields ...logrus.Fields) {
	logger := GetLogger()
	entry := logger.WithError(err)
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Error(message)
}

func Fatal(message string, err error, fields ...logrus.Fields) {
	logger := GetLogger()
	entry := logger.WithError(err)
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Fatal(message)
}
