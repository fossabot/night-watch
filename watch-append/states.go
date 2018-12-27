package watch_append

import (
	"syscall"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"
	"go.uber.org/zap"
	"night-watch/utils/match"
	"os"
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
	States    map[uint64]State
	RecordAt  syscall.Timespec `json:"record_at"`
	TotalSize int64            `json:"total_size"`
}

func (s *States) Save(path string) error {
	blob, err := json.Marshal(s)
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
	err = json.Unmarshal(blob, &states)
	return states, err
}

func NewStates() States {
	return States{
		States: map[uint64]State{},
		TotalSize: 0,
	}
}


// MatchAny checks if the text matches any of the regular expressions
func MatchAny(matchers []match.Matcher, text string) bool {
	for _, m := range matchers {
		if m.MatchString(text) {
			return true
		}
	}
	return false
}

func (s *States) scan(files []string, excludeMatch []match.Matcher, metric *WatchMetric) error {
	for _, path := range files {
		filename := filepath.Base(path)
		if MatchAny(excludeMatch, filename){
			metric.ExcludeCount += 1
			continue
		}

		fi, err := os.Stat(path)
		if err != nil {
			zap.S().Debugw("scan get file os.stat failed",
				"file_path:", path,
				"err", err,
			)
			continue
		}

		if fi.IsDir() {
			continue
		}


		var stat syscall.Stat_t
		if err := syscall.Stat(path, &stat); err != nil {
			zap.S().Debugw("scan get file stat failed",
				"file_path:", path,
				"err", err,
			)
			continue
		}
		state := SysStatToState(path, stat)
		s.States[state.INode] = state
	}
	now := time.Now().UnixNano()
	s.RecordAt = syscall.Timespec{
		Sec:  now / 1e9,
		Nsec: now % 1e9,
	}
	return nil
}


func (s *States) Scan(pattern string, excludeFiles []string, metric *WatchMetric) error {
	files, _ := filepath.Glob(pattern)
	metric.GlobMatchCount = len(files)
	var excludeMatches []match.Matcher
	for _, s := range excludeFiles{
		excludeMatches = append(excludeMatches, match.MustCompile(s))
	}
	zap.S().Debugw("begin scan", "pattern", pattern)
	s.scan(files, excludeMatches,metric)
	zap.S().Debugw("end scan", "files_count", len(files))
	return nil
}
