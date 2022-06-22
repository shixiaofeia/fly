package logging

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type ConsoleLog struct {
	val string
}

func NewConsoleLog() Encoder {
	return new(ConsoleLog)
}

// Config 自定义配置.
func (slf *ConsoleLog) Config() zapcore.Encoder {
	var (
		cfg = zap.NewProductionEncoderConfig()
	)

	// 时间格式自定义
	cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
	}
	// 打印路径自定义
	cfg.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString("[" + getFilePath(caller) + "]")
	}
	// 级别显示自定义
	cfg.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString("[" + level.String() + "]")
	}
	return zapcore.NewConsoleEncoder(cfg)
}

// WithKey 添加单个键.
func (slf *ConsoleLog) WithKey(key string) Encoder {
	slf.val = slf.val + "[" + key + "]    "
	return slf
}

// WithField 添加字段.
func (slf *ConsoleLog) WithField(key, val string) Encoder {
	slf.val = slf.val + fmt.Sprintf("[%s:%s]    ", key, val)
	return slf
}

func (slf *ConsoleLog) Debug(msg string) {
	_logger.Debug(slf.val + msg)
}

func (slf *ConsoleLog) Debugf(format string, v ...interface{}) {
	_logger.Debug(fmt.Sprintf(slf.val+format, v...))
}

func (slf *ConsoleLog) Info(msg string) {
	_logger.Info(slf.val + msg)
}

func (slf *ConsoleLog) Infof(format string, v ...interface{}) {
	_logger.Info(fmt.Sprintf(slf.val+format, v...))
}

func (slf *ConsoleLog) Warn(msg string) {
	_logger.Warn(slf.val + msg)
}

func (slf *ConsoleLog) Warnf(format string, v ...interface{}) {
	_logger.Warn(fmt.Sprintf(slf.val+format, v...))
}

func (slf *ConsoleLog) Error(msg string) {
	_logger.Error(slf.val + msg)
}

func (slf *ConsoleLog) Errorf(format string, v ...interface{}) {
	_logger.Error(fmt.Sprintf(slf.val+format, v...))
}

func (slf *ConsoleLog) Fatal(msg string) {
	_logger.Fatal(slf.val + msg)
}

func (slf *ConsoleLog) Fatalf(format string, v ...interface{}) {
	_logger.Fatal(fmt.Sprintf(slf.val+format, v...))
}
