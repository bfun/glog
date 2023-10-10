package glog_test

import (
	"log/slog"
	"testing"
)

func TestInfo(t *testing.T) {
	for i := 0; i < 1024*1024; i++ {
		slog.Info("test msg", "i", i)
	}
}

func BenchmarkInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slog.Info("benchmark", "i", i)
	}
}
