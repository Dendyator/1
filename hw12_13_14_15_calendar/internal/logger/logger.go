package logger

import (
	"os"

	"github.com/sirupsen/logrus" //nolint:depguard
)

type Logger struct {
	*logrus.Logger
}

func New(level string) *Logger {
	logger := logrus.New()

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	logger.SetLevel(lvl)

	logger.SetOutput(os.Stdout)

	return &Logger{logger}
}
