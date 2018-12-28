package cmd

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	waTesting "night-watch/watch-append/testing"
	"path/filepath"
	"strings"
	"testing"

	"fmt"
)

func runWatchAppend(t *testing.T, fs waTesting.FSMock) map[string]gjson.Result{
	output, err:= executeCommand(rootCmd,

		"watch-append","--once",
		"-p", filepath.Join(fs.Root, "*.log"),
		"-m", filepath.Join(fs.Root, "old_status.json"),
	)
	if err != nil {
		t.Error(err)
	}


	result, err := parseInfluxFormat(output)
	if err != nil {
		t.Error(err)
	}
	return result
}

func TestExecute_NoChange(t *testing.T) {
	fs := waTesting.NewFSMock(t)
	fs.CreateFile("a.log", 1e4)
	fs.CreateFile("b.log", 1e4)
	fs.CreateFile("c.log", 1e4)
	fs.CreateFile("d.log", 1e4)

	// first Result
	firstR := runWatchAppend(t, fs)
	assert.Equal(t,int64(0), firstR["file-append-total-size"].Int())

	// append Result
	fs.AppendFile("a.log", 1e4)
	appendR := runWatchAppend(t, fs)
	assert.Equal(t,int64(1e4), appendR["file-append-total-size"].Int() )


	// noChange Result
	noChangeR := runWatchAppend(t, fs)
	assert.Equal(t,int64(1e4), noChangeR["file-append-total-size"].Int())

	// rotate Result
	fs.RotateFile("a.log",  1e4, false)
	rotateR := runWatchAppend(t, fs)
	assert.Equal(t,int64(2e4), rotateR["file-append-total-size"].Int())

	// deep rotate Result
	fs.RotateFile("a.log",  1e4, true)
	deepRotateR := runWatchAppend(t, fs)
	assert.Equal(t,int64(3e4), deepRotateR["file-append-total-size"].Int())


}

func parseInfluxFormat(in string) (map[string]gjson.Result, error){
	s := strings.Split(in, " ")
	if len(s) != 2 {
		return nil, fmt.Errorf("parse error, have space char %d, expect 1", len(s) -1)
	}
	m := map[string]gjson.Result{}

	for _, kv := range strings.Split(s[1], ","){
		if !strings.Contains(kv, "="){
			return nil, fmt.Errorf("parse error, not have '=',actual:'%s', expect:{Key}={Value}", kv)
		}
		t := strings.Split(kv, "=")
		m[t[0]] = gjson.Parse(t[1])
	}

	return  m , nil
}
