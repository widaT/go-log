package golog

import (
	"testing"
)

func TestLog(t *testing.T) {
	err := Init("test", "INFO", "./log", true, "M", 2)
	if err != nil {
		t.Error("log.Init() fail")
	}

	for i := 0; i < 50; i = i + 1 {
		Logger.Warn("warning msg: %d", i)
		Logger.Info("info msg: %d", i)
	}
}
