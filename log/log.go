package log

import (
	"fmt"
	"runtime"
"reflect"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

func Log(msg string) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := getShortName(pc)
	fmt.Printf("[%s] [%s]\r\n\t%s\r\n", time.Now().Format(timeFormat), funcName, msg)
}




func getShortName(f interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return fullName[strings.LastIndex(fullName, "/")+1:]
}
