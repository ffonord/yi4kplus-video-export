package logger

import "github.com/sirupsen/logrus"

type Logger struct {
	logger *logrus.Logger
}

func New() *Logger {
	return &Logger{
		logger: logrus.New(),
	}
}

func (l *Logger) SetLevel(logLevel string) error {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	l.logger.SetLevel(level)

	return nil
}

func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
