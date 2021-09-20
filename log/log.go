package log

import (
"fmt"
)

const timeFormat = "2006-01-02T15:04:05"

func log(msg string) {
        fmt.Printf("[%s] [%s]\r\n\t%s\r\n", time.Now().Format(timeFormat), "webClient->Get", msg)
}
