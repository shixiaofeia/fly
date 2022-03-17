package logging

import (
	"fmt"

	"go.uber.org/zap"
)

type FieldLog struct {
	fields []zap.Field
}

// NewField  自定义添加log field.
func NewField() *FieldLog {
	return &FieldLog{fields: make([]zap.Field, 0)}
}

func (slf *FieldLog) WithString(key, value string) *FieldLog {
	slf.fields = append(slf.fields, zap.String(key, value))
	return slf
}

func (slf *FieldLog) WithInt(key string, value int) *FieldLog {
	slf.fields = append(slf.fields, zap.Int(key, value))
	return slf
}

func (slf *FieldLog) Debug(msg string) {
	logger.Debug(msg, slf.fields...)
}

func (slf *FieldLog) Debugf(format string, v ...interface{}) {
	logger.Debug(fmt.Sprintf(format, v...), slf.fields...)
}

func (slf *FieldLog) Info(msg string) {
	logger.Info(msg, slf.fields...)
}

func (slf *FieldLog) Infof(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...), slf.fields...)
}

func (slf *FieldLog) Warn(msg string) {
	logger.Warn(msg, slf.fields...)
}

func (slf *FieldLog) Warnf(format string, v ...interface{}) {
	logger.Warn(fmt.Sprintf(format, v...), slf.fields...)
}

func (slf *FieldLog) Error(msg string) {
	logger.Error(msg, slf.fields...)
}

func (slf *FieldLog) Errorf(format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...), slf.fields...)
}

func (slf *FieldLog) Fatal(msg string) {
	logger.Fatal(msg, slf.fields...)
}

func (slf *FieldLog) Fatalf(format string, v ...interface{}) {
	logger.Fatal(fmt.Sprintf(format, v...), slf.fields...)
}
