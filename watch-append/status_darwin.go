package watch_append

import (
	"syscall"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"
	"go.uber.org/zap"
)

type State struct {
	INode    uint64           `json:"inode"`
	Source   string           `json:"source"`
	Size     int64            `json:"size"`
	CreateAt syscall.Timespec `json:"created_at"`
	ModifyAt syscall.Timespec `json:"modify_at"`
	RecordAt syscall.Timespec `json:"record_at"`
}

type States struct {
	States   map[uint64]State
	RecordAt syscall.Timespec `json:"record_at"`
}

func (s *States) Save(path string) error {
	blob, err := json.Marshal(s.States)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, blob, 0644)
	return err
}

func LoadStates(path string) (States, error) {
	states := States{}
	blob, err := ioutil.ReadFile(path)
	if err != nil {
		return states, err
	}
	err = json.Unmarshal(blob, &states.States)
	return states, err
}

func NewStates() States {
	return States{
		States: map[uint64]State{},
	}
}

func (s *States) Scan(pattern string) error {
	files, _ := filepath.Glob(pattern)
	zap.S().Debugw("begin scan", "pattern", pattern)
	for _, path := range files {
		var stat syscall.Stat_t
		if err := syscall.Stat(path, &stat); err != nil {
			zap.S().Debugw("scan get file stat failed",
				"file_path:", path,
				"err", err,
			)
			continue
		}
		now := time.Now().UnixNano()
		state := State{
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
		s.States[state.INode] = state
	}
	now := time.Now().UnixNano()
	s.RecordAt = syscall.Timespec{
		Sec:  now / 1e9,
		Nsec: now % 1e9,
	}
	zap.S().Debugw("end scan", "files_count", len(files))
	return nil
}