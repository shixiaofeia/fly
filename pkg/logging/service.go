package logging

import (
	"fmt"
)

func Sync() {
	_ = logger.Sync()
}

func Debug(msg string) {
	logger.Debug(msg)
}

func Debugf(format string, v ...interface{}) {
	logger.Debug(fmt.Sprintf(format, v...))
}

func Info(msg string) {
	logger.Info(msg)
}

func Infof(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...))
}

func Warn(msg string) {
	logger.Warn(msg)
}

func Warnf(format string, v ...interface{}) {
	logger.Warn(fmt.Sprintf(format, v...))
}

func Error(msg string) {
	logger.Error(msg)
}

func Errorf(format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...))
}

func Fatal(msg string) {
	logger.Fatal(msg)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatal(fmt.Sprintf(format, v...))
}
