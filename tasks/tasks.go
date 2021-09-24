package tasks

import (
	"context"
	"reflect"
	"runtime"
	"time"
)

type Task struct {
	Title  string
	Start  time.Time
	Finish time.Time
	Work   func(ctx context.Context)
}

type TaskRunner interface {
	Run(ctx context.Context)
}

func (t *Task) GetName() string {
	return runtime.FuncForPC(reflect.ValueOf(t.Work).Pointer()).Name()
}

func (t *Task) Run(ctx context.Context) {
	t.Start = time.Now()
	go t.Work(ctx)
}
