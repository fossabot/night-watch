package main

import (
	"night-watch/cmd"
	"runtime/pprof"
)

func main() {
	cmd.Execute()
	defer pprof.StopCPUProfile()

}
