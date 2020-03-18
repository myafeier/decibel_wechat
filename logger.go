package wechat

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

func NewDefaultLogger() ILogger {
	return &Logger{
		Prefix: "[Decibel_wechat]",
	}
}

type Logger struct {
	Prefix string
}

func (self *Logger) Debug(data ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Println(self.Prefix, "[Debug]", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("%s:%d >>>", filepath.Base(file), line), data)

}
func (self *Logger) Info(data ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Println(self.Prefix, "[Info]", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("%s:%d >>>", filepath.Base(file), line), data)
}
func (self *Logger) Error(data ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Println(self.Prefix, "[Error]", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("%s:%d >>>", filepath.Base(file), line), data)
}
