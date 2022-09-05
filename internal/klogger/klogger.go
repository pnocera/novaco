package klogger

import (
	"github.com/pnocera/novaco/internal/utils"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type KLogger struct {
	Name string
	Log  *logrus.Logger
}

func NewKLogger(name string) *KLogger {
	logger := &KLogger{
		Name: name,
	}

	assets, _ := utils.Assets()
	dir := utils.Join(assets, "logs", name)

	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  utils.Join(dir, "info.log"),
		logrus.ErrorLevel: utils.Join(dir, "error.log"),
	}

	logger.Log = logrus.New()
	logger.Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))

	return logger
}

func (l *KLogger) Info(msg string, v ...interface{}) {
	l.Log.Info(msg, v)
}
func (l *KLogger) Error(msg string, v ...interface{}) {
	l.Log.Error(msg, v)
}
