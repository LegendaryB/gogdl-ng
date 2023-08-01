package logging

import (
	"io"
	"os"
	"path"
	"runtime"
	"strconv"

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
	logger.SetReportCaller(true)

	logger.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)

			return "", fileName
		},
	})

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
