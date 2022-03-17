package logging

import "testing"

func TestNewLogger(t *testing.T) {
	Init("../../logs/fly.log")
	defer Sync()
	Debug("debug msg")
	Debugf("debugf %s", "fly")
	Info("info msg")
	Infof("infof %d", 10)
	Warn("warn msg")
	Warnf("warnf %v", true)
	Error("err msg")
	Errorf("errorf %v", []int{1, 2, 3})
	Fatal("fatal msg")
}

func TestFatalf(t *testing.T) {
	Init("../../logs/fly.log")
	defer Sync()
	Fatalf("fatalf %v", map[string]interface{}{"name": "master"})
}
