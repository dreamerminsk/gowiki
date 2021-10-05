package log

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"



type Entry struct {
Time time.Time
Message string
}

func Logf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	Log(msg)
}

func Log(msg string) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := getShortName(pc)
	if funcName == "log.Logf" {
		pc, _, _, _ = runtime.Caller(2)
		funcName = getShortName(pc)
	}
	fmt.Printf("[%s] [%s]\r\n\t%s\r\n", time.Now().Format(timeFormat), funcName, msg)
}

func getShortName(f uintptr) string {
	fullName := runtime.FuncForPC(f).Name()
	return fullName[strings.LastIndex(fullName, "/")+1:]
}
