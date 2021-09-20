package log

import (
"fmt"
"runtime"
)

const timeFormat = "2006-01-02T15:04:05"

func log(msg string) {
        pc, _, _, _ := runtime.Caller(1)
    funcName := runtime.FuncForPC(pc).Name()
        fmt.Printf("[%s] [%s]\r\n\t%s\r\n", time.Now().Format(timeFormat), funcName, msg)
}
