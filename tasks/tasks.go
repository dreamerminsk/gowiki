package tasks

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

type Task struct {
	Title  string
	Start  time.Time
	Finish time.Time
	Message chan<- string
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
		defer func() { t.Finish = time.Now() }()
		t.Work(ctx)
	}()
}

func (t *Task) String() string {
	d := time.Since(t.Start)
	if t.Finish.After(t.Start) {
		d = t.Finish.Sub(t.Start)
	}
	s := fmt.Sprintf("&{%s,\t%s,\t%s}",
		t.Title,
		t.Start.Format(timeFormat),
		d)
	return s
}
