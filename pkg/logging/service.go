package logging

import "go.uber.org/zap"

// NewWithField  自定义log.
func NewWithField(key, value string) *zap.SugaredLogger {
	return Log.With(key, value)
}
