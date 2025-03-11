package glog

import (
	"log/slog"
	"os"
	"path"
	"sync"
	"time"
)

var m sync.Mutex
var logDir = "log"
var dirPerm os.FileMode = 0755
var filePerm os.FileMode = 0644

func newName() string {
	format := "20060102-150405"
	baseName := path.Base(os.Args[0]) + ".log." + time.Now().Format(format)
	return path.Join(logDir, baseName)
}

func newLogger() *slog.Logger {
	m.Lock()
	defer m.Unlock()
	err := os.Mkdir(logDir, dirPerm)
	if err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}
	fileName := newName()
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, filePerm)
	if err != nil {
		panic(err)
	}
	h := slog.NewTextHandler(f, nil)
	return slog.New(h)
}

func init() {
	slog.SetDefault(newLogger())
}
