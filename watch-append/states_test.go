package watch_append

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/assert"
)

type FSMock struct {
	root string
	t    *testing.T
}


func NewFSMock(t *testing.T) FSMock {
	tempPath, _ := ioutil.TempDir("","fs_mock")

	return FSMock{
		root: tempPath,
		t: t,
	}
}



func (m *FSMock) createTempFile(filename string){
	path := filepath.Join(m.root, filename)
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	emptyFile, err := os.Create(path)
	if err != nil {
		m.t.Error("create file failed", err)
	}
	emptyFile.Close()
}

func CreateTestFs(t *testing.T) *FSMock{
	fs := NewFSMock(t)
	fs.createTempFile("a.log")
	fs.createTempFile("b.log")
	fs.createTempFile("c.log")
	fs.createTempFile("_a.log")
	fs.createTempFile("a/a.log")
	fs.createTempFile("a.log.1")
	return &fs
}


func TestStates_Scan(t *testing.T) {
	fs := CreateTestFs(t)
	metric := NewWatchMetric()

	asf := NewStates()
	asf.Scan(
		filepath.Join(fs.root, "*"),
		[]string{},
		&metric)
	assert.Equal(t, 5, len(asf.States))


	asf = NewStates()
	asf.Scan(
		filepath.Join(fs.root, "*.log"),
		[]string{},
		&metric)
	assert.Equal(t, 4, len(asf.States))

	asf = NewStates()
	asf.Scan(
		filepath.Join(fs.root, "*.log"),
		[]string{"_"},
		&metric)
	assert.Equal(t, 3, len(asf.States))
}

func TestStates_Save_And_Load(t *testing.T) {
	fs := CreateTestFs(t)

	metric := NewWatchMetric()
	osf := NewStates()
	path := filepath.Join(fs.root, "state.json")

	osf.Scan(
		filepath.Join(fs.root, "*"),
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