package main

import (
	"runtime/pprof"
	"night-watch/cmd"
)

func main() {
	cmd.Execute()
	defer pprof.StopCPUProfile()

}
