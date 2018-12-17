package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"runtime/pprof"
)

var rootCmd = &cobra.Command{
	Use:   "nightwatch",
	Short: "Dashbase NightWatch",
	Long: `Dashbase NightWatch`,
}

var debug bool
var isLogSave bool
var logPath string
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
	rootCmd.PersistentFlags().BoolVar(&isLogSave, "log", false, "save nightwatch log")
	rootCmd.PersistentFlags().StringVar(&logPath, "log-path", "/tmp/nightwatch.log", "only when --log=true, save log to this path")

}

func initConfig() {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	if isLogSave {
		cfg.OutputPaths = []string{
			logPath,
		}
	}
	if debug {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

		f, err := os.Create("./cpu.pprof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}