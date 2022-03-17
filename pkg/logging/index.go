package logging

import (
	"os"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type logItem struct {
	FileName string
	Level    zap.LevelEnablerFunc
}

var (
	maxSize    = 200 // 每个日志文件最大尺寸200M
	maxBackups = 20  // 日志文件最多保存20个备份
	maxAge     = 30

	logger *zap.Logger
	_pool  = buffer.NewPool()
)

// Init 初始化日志.
func Init(logPath string) {
	prefix, suffix := getFileSuffixPrefix(logPath)

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

// NewLogger
func NewLogger(items []logItem) {
	var (
		cfg   = zap.NewProductionEncoderConfig()
		cores []zapcore.Core
	)

	// 时间格式自定义
	cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	// 打印路径自定义
	cfg.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(getFilePath(caller))
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
			zapcore.NewJSONEncoder(cfg), // 编码器配置
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
	logger = zap.New(zapcore.NewTee(cores...), caller, development, skip)
	return
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
