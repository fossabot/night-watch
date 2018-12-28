package watch_append

import (
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/assert"
	waTesting "night-watch/watch-append/testing"
)


func TestStates_Scan(t *testing.T) {
	fs := waTesting.CreateTestFs(t)
	metric := NewWatchMetric()

	asf := NewStates()
	asf.Scan(
		filepath.Join(fs.Root, "*"),
		[]string{},
		&metric)
	assert.Equal(t, 5, len(asf.States))


	asf = NewStates()
	asf.Scan(
		filepath.Join(fs.Root, "*.log"),
		[]string{},
		&metric)
	assert.Equal(t, 4, len(asf.States))

	asf = NewStates()
	asf.Scan(
		filepath.Join(fs.Root, "*.log"),
		[]string{"_"},
		&metric)
	assert.Equal(t, 3, len(asf.States))
}

func TestStates_Save_And_Load(t *testing.T) {
	fs := waTesting.CreateTestFs(t)

	metric := NewWatchMetric()
	osf := NewStates()
	path := filepath.Join(fs.Root, "state.json")

	osf.Scan(
		filepath.Join(fs.Root, "*"),
		[]string{},
		&metric)
	osf.Save(path)
	asf, err := LoadStates(path)

	if err != nil{
		assert.Error(t, err)
	}
	isEqual := assert.ObjectsAreEqual(asf, osf)
	assert.Equal(t, true, isEqual)
}