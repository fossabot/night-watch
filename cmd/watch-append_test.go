package cmd

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"night-watch/watch-append"
	"os"
	"strings"
	"testing"
	"log"
)

var testpath string
func init(){
	tmp, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
		return
	}
	testpath = tmp
}

func setupTestCase(t *testing.T) {
	t.Log("setup test case")
	if _, err := os.Stat(testpath); err == nil {
		os.RemoveAll(testpath)
	}
	if _, err := os.Stat(testpath); os.IsNotExist(err) {
		os.Mkdir(testpath, 0744)
	}
}

func TestStatusInfo(t *testing.T) {
	setupTestCase(t)
	f1 := []byte("test\ndashbase\n")
	//length := len(f1)
	err := ioutil.WriteFile(testpath+string(os.PathSeparator)+"dat1.log", f1, 0644)
	checkError(t, err)
	f2 := []byte("test\ndashbase\n")
	err = ioutil.WriteFile(testpath+string(os.PathSeparator)+"dat2.log", f2, 0644)
	checkError(t, err)
	f3 := []byte("test\ndashbase\n")
	err = ioutil.WriteFile(testpath+string(os.PathSeparator)+"_dat3.log", f3, 0644)
	checkError(t, err)

	state := watch_append.NewStates()
	metric := watch_append.NewWatchMetric()
	metric.Start()

	state.Scan(testpath+string(os.PathSeparator)+"*.log", []string{"_"}, &metric)
	statusfile := testpath + string(os.PathSeparator) + "status.txt"
	state.Save(statusfile)
	dat, err := ioutil.ReadFile(statusfile)
	checkError(t, err)

	result := string(dat)
	if len(result) == 0 {
		assert.Fail(t, "empty state file")
	}
	osf, err := watch_append.LoadStates(statusfile)

	if len(osf.States) <= 0 {
		assert.Fail(t, "wrong state file")
	}

	checkData(t, osf)
}

func TestAppend(t *testing.T) {
	osf, err := watch_append.LoadStates(testpath + string(os.PathSeparator) + "status.txt")
	if err != nil {
		checkError(t, err)
	}
	metric := watch_append.NewWatchMetric()
	metric.Start()

	f, err := os.OpenFile(testpath+string(os.PathSeparator)+"dat1.log", os.O_APPEND|os.O_WRONLY, 0600)
	checkError(t, err)
	defer f.Close()
	if _, err = f.WriteString("add a line\n"); err != nil {
		checkError(t, err)
	}

	asf := watch_append.NewStates()
	asf.Scan(testpath+string(os.PathSeparator)+"*.log", []string{"_"}, &metric)

	diff := watch_append.NewDiff(asf, osf, pattern, &metric)
	diff.Diff()
	asf.TotalSize = diff.Result.TotalSize
	asf.Save(testpath + string(os.PathSeparator) + "status.txt")
	assert.Equal(t, int64(11), diff.Result.TotalSize)
	assert.Equal(t, int64(2), diff.Result.Count)
}

func checkData(t *testing.T, osf watch_append.States) {
	for _, v := range osf.States {
		if (strings.Contains(v.Source, "dat1") || strings.Contains(v.Source, "dat2")) && !strings.Contains(v.Source, "dat3") {
			t.Log(v.Source)
			assert.True(t, true, "string contains dat1 or dat2 and not dat3")
		} else {
			assert.Fail(t, "missing log")
		}
	}
}

func checkError(t *testing.T, e error) {
	if e != nil {
		assert.Fail(t, e.Error())
	}
}
