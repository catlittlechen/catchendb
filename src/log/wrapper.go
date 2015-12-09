package log

import (
	"catchendb/src/vendor/code.google.com/p/log4go"
	"fmt"
)

func Init(configPath string) {
	log4go.LoadConfiguration(configPath)
}

func Critical(format string, args ...interface{}) {
	err := log4go.Critical(format, args...)
	if err != nil {
		fmt.Println(err)
	}
}

func Errorf(format string, args ...interface{}) {
	err := log4go.Error(format, args...)
	if err != nil {
		fmt.Println(err)
	}
}

func Warn(format string, args ...interface{}) {
	err := log4go.Warn(format, args...)
	if err != nil {
		fmt.Println(err)
	}
}

func Info(format string, args ...interface{}) {
	log4go.Info(format, args...)
}

func Trace(format string, args ...interface{}) {
	log4go.Trace(format, args...)
}

func Debug(format string, args ...interface{}) {
	log4go.Debug(format, args...)
}

func Fine(format string, args ...interface{}) {
	log4go.Fine(format, args...)
}

func Finest(format string, args ...interface{}) {
	log4go.Finest(format, args...)
}
