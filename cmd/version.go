package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"night-watch/pkg/build"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of NightWatch",
	Long:  `Print the version number of NightWatch`,
	Run: func(cmd *cobra.Command, args []string) {
		info := build.Get()
		fmt.Printf("NightWatch BuildDate: %s, commit: %s\n", info.BuildDate, info.GitCommit)
	},
}