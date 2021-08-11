package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var Log *zap.Logger

var (
	MaxSize    = 200 // 每个日志文件最大尺寸200M
	MaxBackups = 20  // 日志文件最多保存20个备份
)

// Init 初始化日志
func Init(logPath string, maxAge int, compress bool) {
	Log = NewLogger(logPath, zap.InfoLevel, MaxSize, MaxBackups, maxAge, compress)
}

// NewLogger
func NewLogger(logPath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logPath,    // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
		LocalTime:  true,       // 备份文件名本地/UTC时间
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "file",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		FunctionKey:    "func",                         // 方法名
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 路径
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 时间格式自定义
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel,                                                                     // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开发模式
	development := zap.Development()
	// 忽略封装器
	skip := zap.AddCallerSkip(1)
	// 构造日志
	logger := zap.New(core, caller, development, skip)
	return logger
}
