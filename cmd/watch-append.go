package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"night-watch/watch-append"
	"fmt"
	"os"
)

var pattern string

func init() {
	rootCmd.AddCommand(watchAppendCmd)
	watchAppendCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "an absolute path to watch change")
	watchAppendCmd.MarkFlagRequired("pattern")
}

var watchAppendCmd = &cobra.Command{
	Use:   "watch-append",
	Short: "watch-append",
	Long:  `watch-append`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("/tmp/night-watch.old_status.json"); os.IsNotExist(err) {
			asf := watch_append.NewStates()
			asf.Scan(pattern)
			asf.Save("/tmp/night-watch.old_status.json")
			fmt.Printf("%f\n", float64(0))
			return
		}

		// osf == old_status_of_file
		zap.S().Debug("load osf")
		osf, err := watch_append.LoadStates("/tmp/night-watch.old_status.json")
		if err != nil {
			zap.S().Error("Load old status failed.")
		}
		zap.S().Debug("finished osf")

		// asf == actual_status_of_file
		zap.S().Debug("load asf")
		asf := watch_append.NewStates()
		asf.Scan(pattern)
		asf.Save("/tmp/night-watch.old_status.json")
		zap.S().Debug("finished asf")

		zap.S().Debug("start diff")
		diff := watch_append.NewDiff(asf, osf)
		fmt.Printf("%f\n", diff.GetTotalRate())
		zap.S().Debug("end diff")
	},
}
