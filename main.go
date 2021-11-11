package main

import (
	"github.com/apache/apisix-go-plugin-runner/pkg/runner"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg := runner.RunnerConfig{}
	cfg.LogLevel = zapcore.DebugLevel
	runner.Run(cfg)
}
