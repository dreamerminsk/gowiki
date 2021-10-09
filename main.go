package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dreamerminsk/gowiki/tasks"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	ticker := time.NewTicker(time.Duration(60) * time.Second)
	defer ticker.Stop()

	keyChan := make(chan os.Signal, 1)
	signal.Notify(keyChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := [...]*tasks.Task{
		tasks.New(tasks.InitCategories),
		tasks.New(tasks.InitForums),
		tasks.New(tasks.InitUsers),
		tasks.New(tasks.UpdateForums),
		tasks.New(tasks.UpdateTopics),
	}
	for _, t := range queue {
		go t.Run(ctx)
	}

	for {
		select {
		case <-keyChan:
			for _, t := range queue {
				fmt.Println(t)
			}
			return
		case <-ticker.C:
			for _, t := range queue {
				fmt.Println(t)
			}
		}
	}

}
