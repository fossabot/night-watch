package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"night-watch/watch-append"
	"fmt"
	"os"
	"time"
	"strings"
)

var pattern string
var metaPath string
func init() {
	rootCmd.AddCommand(watchAppendCmd)
	watchAppendCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "an absolute path to watch change")
	watchAppendCmd.Flags().StringVarP(&metaPath, "mate-path", "m", "/tmp/night-watch.old_status.json", "the path to read/save metadata")
	watchAppendCmd.MarkFlagRequired("pattern")
}

var watchAppendCmd = &cobra.Command{
	Use:   "watch-append",
	Short: "watch-append",
	Long:  `watch-append`,
	Run: func(cmd *cobra.Command, args []string) {
		start_at := time.Now().UnixNano()
		if _, err := os.Stat(metaPath); os.IsNotExist(err) {
			asf := watch_append.NewStates()
			asf.Scan(pattern)
			asf.Save(metaPath)
			fmt.Printf("%f\n", float64(0))
			return
		}

		// osf == old_status_of_file
		zap.S().Debug("load osf")
		osf, err := watch_append.LoadStates(metaPath)
		if err != nil {
			zap.S().Error("Load old status failed.")
		}
		zap.S().Debug("finished osf")

		// asf == actual_status_of_file
		zap.S().Debug(">> load asf")
		asf := watch_append.NewStates()
		asf.Scan(pattern)

		zap.S().Debug(">> >>	start diff")
		diff := watch_append.NewDiff(asf, osf)
		diff.Diff()
		diff.Result.TotalSize += osf.TotalSize
		zap.S().Debug("<< << 	end diff")

		asf.SetTotalSize(diff.Result.TotalSize)
		asf.Save(metaPath)

		zap.S().Debug("<< finished asf")



		spend := time.Now().UnixNano() - start_at
		fmt.Fprintln(os.Stdout, create_influx_format_outout(diff.Result, spend))
	},
}


func create_influx_format_outout(result watch_append.DiffResult, spend int64) string{
	//	weather,location=us-midwest temperature=82 1465839830100400200
	//  	|    -------------------- --------------  |
	//  	|             |             |             |
	//  	|             |             |             |
	//	+-----------+--------+-+---------+-+---------+
	//	|measurement|,tag_set| |field_set| |timestamp|
	//	+-----------+--------+-+---------+-+---------+
	str := strings.Builder{}

	// measurement
	str.Write([]byte("nightwatch-metric"))

	// tag_set
	str.Write([]byte(" "))

	// field_set

	str.Write([]byte(fmt.Sprintf("%s=%f,", "file-append-speed", result.Speed)))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "file-append-total-size", result.TotalSize)))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "scan-file-count", result.Count)))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "run-spent-time-nanosecond", spend)))
	str.Write([]byte(fmt.Sprintf("%s=%f,", "run-spent-time-second", float64(spend) / 1e9)))

	// remove last char ','
	return strings.TrimSuffix(str.String(), ",")
}