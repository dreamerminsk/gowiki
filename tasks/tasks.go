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
	Title   string
	Start   time.Time
	Finish  time.Time
	Message string
	MsgChan chan string
	Work    func(ctx context.Context, t *Task)
}

type TaskRunner interface {
	Run(ctx context.Context)
}

func New(f func(ctx context.Context, t *Task)) *Task {
	return &Task{
		Title:   getShortName(f),
		MsgChan: make(chan string),
		Work:    f,
	}
}

func getShortName(f interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return fullName[strings.LastIndex(fullName, "/")+1:]
}

func (t *Task) Run(ctx context.Context) {
	t.Start = time.Now()
	defer close(t.MsgChan)
	go func() {
		defer func() { t.Finish = time.Now() }()
		t.Work(ctx, t)
	}()
	for {
		select {
		case msg := <-t.MsgChan:
			t.Message = msg
		case <-ctx.Done():
			return
		}
	}
}

func (t *Task) String() string {
	d := time.Since(t.Start)
	if t.Finish.After(t.Start) {
		d = t.Finish.Sub(t.Start)
	}
	s := fmt.Sprintf("&{%s,\t%s,\t%s\n\r  %s}",
		t.Title,
		t.Start.Format(timeFormat),
		d, t.Message)
	return s
}
