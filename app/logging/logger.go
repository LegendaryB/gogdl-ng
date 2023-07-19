package logging

import (
	"io"
	"os"

	"github.com/LegendaryB/gogdl-ng/app/config"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

func NewLogger(conf config.LoggingConfiguration) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	lvl, err := logrus.ParseLevel(conf.LogLevel)

	if err != nil {
		return nil, err
	}

	logger.SetLevel(lvl)

	file, err := os.OpenFile(conf.LogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)

	if err != nil {
		return nil, err
	}

	if conf.LogToConsole {
		mw := io.MultiWriter(os.Stdout, file)
		logger.SetOutput(mw)
	} else {
		logger.SetOutput(file)
	}

	logrus.RegisterExitHandler(func() {
		file.Close()
	})

	return logger, nil
}
