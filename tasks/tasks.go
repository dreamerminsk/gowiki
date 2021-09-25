package tasks

import (
	"context"
	"reflect"
	"runtime"
	"strings"
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

func New(f func(ctx context.Context)) *Task {
	return &Task{
		Title: getShortName(f),
		Work:  f,
	}
}

func getShortName(f interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return fullName[strings.LastIndex(fullName, "/")+1:]
}

func (t *Task) Run(ctx context.Context) {
	t.Start = time.Now()
	go t.Work(ctx)
}
