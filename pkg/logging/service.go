package logging

import (
	"fmt"
)

func Sync() {
	_ = _logger.Sync()
}

func Debug(msg string) {
	_logger.Debug(msg)
}

func Debugf(format string, v ...interface{}) {
	_logger.Debug(fmt.Sprintf(format, v...))
}

func Info(msg string) {
	_logger.Info(msg)
}

func Infof(format string, v ...interface{}) {
	_logger.Info(fmt.Sprintf(format, v...))
}

func Warn(msg string) {
	_logger.Warn(msg)
}

func Warnf(format string, v ...interface{}) {
	_logger.Warn(fmt.Sprintf(format, v...))
}

func Error(msg string) {
	_logger.Error(msg)
}

func Errorf(format string, v ...interface{}) {
	_logger.Error(fmt.Sprintf(format, v...))
}

func Fatal(msg string) {
	_logger.Fatal(msg)
}

func Fatalf(format string, v ...interface{}) {
	_logger.Fatal(fmt.Sprintf(format, v...))
}
