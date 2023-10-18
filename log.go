package glog

import (
	"log/slog"
	"os"
	"path"
	"sync"
	"time"
)

var Size int64 = 1024 * 1024 * 10
var fileName string
var m sync.Mutex
var count int

func newName(bak bool) string {
	prefix := "log/" + path.Base(os.Args[0]) + ".log."
	format := "20060102"
	if bak {
		format = "20060102-150405.000"
	}
	return prefix + time.Now().Format(format)
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key != "time" {
		return a
	}
	a = slog.String("time", time.Now().Format("2006-01-02 15:04:05.000000"))
	s1 := newName(false)
	if s1 != fileName {
		slog.SetDefault(newLogger())
		return a
	}
	if count < 1000 {
		count++
		return a
	}
	count = 0
	stat, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}
	if stat.Size() < Size {
		return a
	}
	s2 := newName(true)
	err = os.Rename(fileName, s2)
	if err != nil {
		panic(err)
	}
	slog.SetDefault(newLogger())
	return a
}

func newLogger() *slog.Logger {
	m.Lock()
	defer m.Unlock()
	err := os.Mkdir("log", 0755)
	if err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}
	fileName = newName(false)
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	o := slog.HandlerOptions{ReplaceAttr: replaceAttr}
	h := slog.NewTextHandler(f, &o)
	return slog.New(h)
}

func init() {
	slog.SetDefault(newLogger())
}
