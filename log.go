package glog

import (
	"log/slog"
	"os"
	"path"
	"path/filepath"
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

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == "source" {
		source := a.Value.Any().(*slog.Source)
		fileName := filepath.Base(source.File)
		return slog.Attr{
			Key:   a.Key,
			Value: slog.StringValue(fileName),
		}
	}
	return a
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
	options := slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo, ReplaceAttr: replaceAttr}
	h := slog.NewTextHandler(f, &options)
	return slog.New(h)
}

func init() {
	slog.SetDefault(newLogger())
}
