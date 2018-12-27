package watch_append

import (
	"syscall"
	"time"
)

func SysStatToState(path string,stat syscall.Stat_t) State {
	now := time.Now().UnixNano()
	return State{
		INode:    stat.Ino,
		Source:   path,
		Size:     stat.Size,
		CreateAt: stat.Ctimespec,
		ModifyAt: stat.Mtimespec,
		RecordAt: syscall.Timespec{
			Sec:  now / 1e9,
			Nsec: now % 1e9,
		},
	}

}