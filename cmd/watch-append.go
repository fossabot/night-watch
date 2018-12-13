package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"night-watch/watch-append"
	"fmt"
)

var pattern string
var sugar *zap.SugaredLogger

func init() {
	rootCmd.AddCommand(watchAppendCmd)
	watchAppendCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "an absolute path to watch change")
	watchAppendCmd.MarkFlagRequired("pattern")
	logger, _ := zap.NewProduction()
	sugar = logger.Sugar()
}

var watchAppendCmd = &cobra.Command{
	Use:   "watch-append",
	Short: "watch-append",
	Long:  `watch-append`,
	Run: func(cmd *cobra.Command, args []string) {
		// osf == old_status_of_file
		osf, err := watch_append.LoadStates("/tmp/night-watch.old_status.json")
		if err != nil {
			sugar.Error("Load old status failed.")
		}
		// asf == actual_status_of_file
		asf := watch_append.NewStates()
		asf.Scan(pattern)
		asf.Save("/tmp/night-watch.old_status.json")
		diff := watch_append.NewDiff(asf, osf)
		fmt.Printf("%f\n", diff.GetTotalRate())
	},
}
