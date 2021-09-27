package tasks

import (
	"context"
	"fmt"
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
        go func() {
           defer func() {t.Finish = time.Now()}()
           t.Work(ctx)
        }()
}

func (t *Task) String() string {
	s := fmt.Sprintf("&{%s, %s, %s}",
		t.Title,
		t.Start.Format("2006-01-02 15:04:05"),
		t.Finish.Format("2006-01-02 15:04:05"))
	return s
}
