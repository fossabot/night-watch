package testing

import (
	"testing"
	"io/ioutil"
	"path/filepath"
	"os"
	"log"
)

type FSMock struct {
	Root string
	t    *testing.T
}

func NewFSMock(t *testing.T) FSMock {
	tempPath, _ := ioutil.TempDir("", "fs_mock")

	return FSMock{
		Root: tempPath,
		t:    t,
	}
}

func (m *FSMock) CreateFile(filename string, size int64) string {
	path := filepath.Join(m.Root, filename)
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	f, err := os.Create(path)
	if err != nil {
		m.t.Error("create file failed", err)
	}
	if err := f.Truncate(size); err != nil {
		log.Fatal(err)
	}
	_ = f.Close()
	return path
}

func (m *FSMock) AppendFile(filename string, size int64) string {
	path := filepath.Join(m.Root, filename)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		m.t.Error("create file failed", err)
	}
	buf := make([]byte, size)
	if _, err := f.Write(buf); err != nil {
		log.Fatal(err)
	}
	_ = f.Close()
	return path
}


func (m *FSMock) RotateFile(filename string, size int64, deep bool)  {
	path := filepath.Join(m.Root, filename)
	to := filepath.Join(m.Root, filename + ".1")
	if deep {
		to = filepath.Join(m.Root, "deep-rotate" + filename + ".1")
	}
	err := os.Rename(path, to)
	if err != nil {
		m.t.Error("create file failed", err)
	}
	m.CreateFile(filename, size)
}


func (m *FSMock) CreateEmptyFile(filename string) string {
	path := filepath.Join(m.Root, filename)
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	emptyFile, err := os.Create(path)
	if err != nil {
		m.t.Error("create file failed", err)
	}
	_ = emptyFile.Close()
	return path
}

func CreateTestFs(t *testing.T) *FSMock {
	fs := NewFSMock(t)
	fs.CreateEmptyFile("a.log")
	fs.CreateEmptyFile("b.log")
	fs.CreateEmptyFile("c.log")
	fs.CreateEmptyFile("_a.log")
	fs.CreateEmptyFile("a/a.log")
	fs.CreateEmptyFile("a.log.1")
	return &fs
}
