package cmd

import (
	"github.com/spf13/cobra"
	"night-watch/watch-append"
	"fmt"
	"strings"
	"time"
	"os"
	"go.uber.org/zap"
)

var pattern string
var excludeFiles []string
var metaPath string
var isOnce bool
var interval int

func init() {
	rootCmd.AddCommand(watchAppendCmd)
	watchAppendCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "an absolute path to watch change")
	watchAppendCmd.Flags().StringArrayVarP(&excludeFiles, "exclude-file", "e", []string{}, "A list of regular expressions to match the files that you want nightwatch to ignore")
	watchAppendCmd.Flags().StringVarP(&metaPath, "mate-path", "m", "/tmp/night-watch.old_status.json", "the path to read/save metadata")
	watchAppendCmd.Flags().BoolVarP(&isOnce, "once", "o", true, "run it once")
	watchAppendCmd.Flags().IntVar(&interval, "interval", 1, "each run interval")
	watchAppendCmd.MarkFlagRequired("pattern")
}

var watchAppendCmd = &cobra.Command{
	Use:   "watch-append",
	Short: "watch-append",
	Long:  `watch-append`,
	Run: func(cmd *cobra.Command, args []string) {
		for {
			run()
			if isOnce {
				return
			}
			time.Sleep(time.Second * time.Duration(interval))
		}
	},
}

func run() {
	metric := watch_append.NewWatchMetric()
	metric.Start()
	// first run watch append, create first state and save
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		asf := watch_append.NewStates()
		asf.Scan(pattern, excludeFiles, &metric)
		asf.Save(metaPath)
		fmt.Fprintln(os.Stdout, create_influx_format_outout(watch_append.DiffResult{}, metric))
		return
	}

	// osf == old_status_of_file
	osf, err := watch_append.LoadStates(metaPath)
	if err != nil {
		zap.S().Error("Load old status failed.")
	}

	// asf == actual_status_of_file
	asf := watch_append.NewStates()
	asf.Scan(pattern, excludeFiles, &metric)

	diff := watch_append.NewDiff(asf, osf, pattern, &metric)
	diff.Diff()
	diff.Result.TotalSize += osf.TotalSize
	asf.SetTotalSize(diff.Result.TotalSize)
	asf.Save(metaPath)
	metric.End()

	fmt.Fprintln(os.Stdout, create_influx_format_outout(diff.Result, metric))
}

func create_influx_format_outout( result watch_append.DiffResult, metric watch_append.WatchMetric) string {
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
	str.Write([]byte(fmt.Sprintf("%s=%d,", "file-append-total-size", result.TotalSize)))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "scan-file-count", result.Count)))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "run-spent-time-nanosecond", metric.Spent)))
	str.Write([]byte(fmt.Sprintf("%s=%f,", "run-spent-time-second", float64(metric.Spent)/1e9)))


	str.Write([]byte(fmt.Sprintf("%s=%d,", "glob_count", metric.GlobMatchCount)))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "exclude_count", metric.ExcludeCount)))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "as_count", metric.AsCount)))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "os_count", metric.OsCount)))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "as_union_os_count", metric.AsUnionOs)))

	str.Write([]byte(fmt.Sprintf("%s=%d,", "diff-new", metric.DiffNewFile.Get())))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "diff-rotate", metric.DiffRotate.Get())))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "diff-nochange", metric.DiffNoChange.Get())))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "diff-append", metric.DiffAppend.Get())))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "diff-null", metric.DiffNullError.Get())))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "diff-r-ignore", metric.DiffRotateIgnore.Get())))
	str.Write([]byte(fmt.Sprintf("%s=%d,", "diff-r-noexists", metric.DiffRotateNotExists.Get())))
	// remove last char ','
	return strings.TrimSuffix(str.String(), ",")
}
