package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"strings"
)

type (
	Conf struct {
		Path    string // 日志路径
		Encoder string // 编码器选择
	}
	logItem struct {
		FileName string
		Level    zap.LevelEnablerFunc
	}
	Encoder interface {
		Config() zapcore.Encoder
		WithKey(key string) Encoder
		WithField(key, val string) Encoder
		Debug(msg string)
		Debugf(format string, v ...interface{})
		Info(msg string)
		Infof(format string, v ...interface{})
		Warn(msg string)
		Warnf(format string, v ...interface{})
		Error(msg string)
		Errorf(format string, v ...interface{})
		Fatal(msg string)
		Fatalf(format string, v ...interface{})
	}
)

var (
	maxSize    = 200 // 每个日志文件最大尺寸200M
	maxBackups = 20  // 日志文件最多保存20个备份
	maxAge     = 30  // 保留最大天数
	_logger    *zap.Logger
	_pool      = buffer.NewPool()
	c          Conf

	ConsoleEncoder = "console" // 控制台输出
	JsonEncoder    = "json"    // json输出
)

// Init 初始化日志.
func Init(conf Conf) {
	c = conf
	prefix, suffix := getFileSuffixPrefix(c.Path)

	infoPath := path.Join(prefix + ".info" + suffix)
	errPath := path.Join(prefix + ".err" + suffix)
	items := []logItem{
		{
			FileName: infoPath,
			Level: func(level zapcore.Level) bool {
				return level <= zap.InfoLevel
			},
		},
		{
			FileName: errPath,
			Level: func(level zapcore.Level) bool {
				return level > zap.InfoLevel
			},
		},
	}

	NewLogger(items)
}

// NewLogger 日志.
func NewLogger(items []logItem) {
	var (
		cfg   zapcore.Encoder
		cores []zapcore.Core
	)
	switch c.Encoder {
	case JsonEncoder:
		cfg = NewJsonLog().Config()
	case ConsoleEncoder:
		cfg = NewConsoleLog().Config()
	default:
		cfg = NewConsoleLog().Config()
	}

	for _, v := range items {
		hook := lumberjack.Logger{
			Filename:   v.FileName,
			MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: maxBackups, // 日志文件最多保存多少个备份
			MaxAge:     maxAge,     // 文件最多保存多少天
			Compress:   true,       // 是否压缩
			LocalTime:  true,       // 备份文件名本地/UTC时间
		}
		core := zapcore.NewCore(
			cfg, // 编码器配置;
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
			v.Level, // 日志级别
		)
		cores = append(cores, core)
	}

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开发模式
	development := zap.Development()
	// 二次封装
	skip := zap.AddCallerSkip(1)
	// 构造日志
	_logger = zap.New(zapcore.NewTee(cores...), caller, development, skip)
	return
}

// GetEncoder 获取自定义编码器.
func GetEncoder() Encoder {
	switch c.Encoder {
	case JsonEncoder:
		return NewJsonLog()
	case ConsoleEncoder:
		return NewConsoleLog()
	default:
		return NewConsoleLog()
	}
}

// GetLogger 获取日志记录器.
func GetLogger() *zap.Logger {
	return _logger
}

// getFileSuffixPrefix 文件路径切割
func getFileSuffixPrefix(fileName string) (prefix, suffix string) {
	paths, _ := path.Split(fileName)
	base := path.Base(fileName)
	suffix = path.Ext(fileName)
	prefix = strings.TrimSuffix(base, suffix)
	prefix = path.Join(paths, prefix)
	return
}

// getFilePath 自定义获取文件路径.
func getFilePath(ec zapcore.EntryCaller) string {
	if !ec.Defined {
		return "undefined"
	}
	buf := _pool.Get()
	buf.AppendString(ec.Function)
	buf.AppendByte(':')
	buf.AppendInt(int64(ec.Line))
	caller := buf.String()
	buf.Free()
	return caller
}
