package tasks

import (
	"context"
	"time"
)

type Task struct {
	Title  string
	Start  time.Time
	Finish time.Time
	Work   func(ctx context.Context)
}

type TaskRunner interface {
	Run()
}
