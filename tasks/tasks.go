package tasks

import "time"

type Task struct {
	Title  string
	Start  time.Time
	Finish time.Time
}
