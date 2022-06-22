package logging

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"time"

	"go.uber.org/zap"
)

type JsonLog struct {
	fields []zap.Field
	val    string
}

// NewJsonLog  自定义添加log field.
func NewJsonLog() Encoder {
	return &JsonLog{fields: make([]zap.Field, 0)}
}

// Config 自定义配置.
func (slf *JsonLog) Config() zapcore.Encoder {
	var (
		cfg = zap.NewProductionEncoderConfig()
	)

	// 时间格式自定义
	cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	// 打印路径自定义
	cfg.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(getFilePath(caller))
	}
	// 级别显示自定义
	cfg.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(level.String())
	}
	return zapcore.NewJSONEncoder(cfg)
}

// WithKey 添加单个键.
func (slf *JsonLog) WithKey(key string) Encoder {
	slf.val = slf.val + key + " "
	return slf
}

// WithField 添加字段.
func (slf *JsonLog) WithField(key, val string) Encoder {
	slf.fields = append(slf.fields, zap.String(key, val))
	return slf
}

func (slf *JsonLog) Debug(msg string) {
	_logger.Debug(slf.val+msg, slf.fields...)
}

func (slf *JsonLog) Debugf(format string, v ...interface{}) {
	_logger.Debug(fmt.Sprintf(slf.val+format, v...), slf.fields...)
}

func (slf *JsonLog) Info(msg string) {
	_logger.Info(slf.val+msg, slf.fields...)
}

func (slf *JsonLog) Infof(format string, v ...interface{}) {
	_logger.Info(fmt.Sprintf(slf.val+format, v...), slf.fields...)
}

func (slf *JsonLog) Warn(msg string) {
	_logger.Warn(slf.val+msg, slf.fields...)
}

func (slf *JsonLog) Warnf(format string, v ...interface{}) {
	_logger.Warn(fmt.Sprintf(slf.val+format, v...), slf.fields...)
}

func (slf *JsonLog) Error(msg string) {
	_logger.Error(slf.val+msg, slf.fields...)
}

func (slf *JsonLog) Errorf(format string, v ...interface{}) {
	_logger.Error(fmt.Sprintf(slf.val+format, v...), slf.fields...)
}

func (slf *JsonLog) Fatal(msg string) {
	_logger.Fatal(slf.val+msg, slf.fields...)
}

func (slf *JsonLog) Fatalf(format string, v ...interface{}) {
	_logger.Fatal(fmt.Sprintf(slf.val+format, v...), slf.fields...)
}
