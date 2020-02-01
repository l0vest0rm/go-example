package main

import (
	"fmt"
)

var (
	isDebug = true //调试开关是否打开
)

// Printf 打印日志
func Printf(format string, a ...interface{}) (n int, err error) {
	if isDebug {
		return fmt.Printf(format, a...)
	}

	return 0, nil
}

// Println 打印日志
func Println(a ...interface{}) (n int, err error) {
	if isDebug {
		return fmt.Println(a...)
	}

	return 0, nil
}
