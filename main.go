package main


import (
	"go.uber.org/zap"
	"night-watch/cmd"
)

func main() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"/tmp/nightwatch.log",
	}
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
	cmd.Execute()
}
